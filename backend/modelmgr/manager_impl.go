package modelmgr

import (
	"ai-eino-interview-agent/modelmgr/factory"
	"ai-eino-interview-agent/modelmgr/factory/builders"
	"ai-eino-interview-agent/modelmgr/manager"
	"context"
	"fmt"
	"time"
)

// defaultModelManager 是 ModelManager 接口的默认实现。
type defaultModelManager struct {
	mgr     manager.Manager // 内部模型管理器
	factory factory.Factory // 模型工厂
	logger  Logger          // 日志记录器
	timeout time.Duration   // 默认超时时间
}

// New 创建一个新的 ModelManager 实例。
func New(cfg *Config) (ModelManager, error) {
	if cfg == nil {
		return nil, NewError("INVALID_CONFIG", "config cannot be nil", ErrInvalidConfig)
	}

	// 1. 加载模型元数据
	var models []*manager.InternalModel
	var err error

	if cfg.LoadFromEnv {
		if cfg.TemplatePath == "" {
			return nil, NewError("INVALID_CONFIG", "template_path is required when load_from_env is true", ErrInvalidConfig)
		}
		models, err = manager.LoadFromEnv(cfg.TemplatePath)
	} else {
		// 解析配置路径（支持自动查找）
		configPath, err := ResolveConfigPath(cfg.ConfigPath)
		if err != nil {
			return nil, NewError("INVALID_CONFIG", fmt.Sprintf("failed to resolve config path: %v", err), ErrInvalidConfig)
		}

		// 记录使用的配置路径
		if cfg.Logger != nil {
			cfg.Logger.Info("Using config path", "path", configPath)
		}

		models, err = manager.LoadFromYAML(configPath)
	}

	if err != nil {
		return nil, NewError("LOAD_CONFIG_FAILED", "failed to load model configurations", err)
	}

	// 2. 创建内部管理器
	mgr := manager.NewStaticManager(models)

	// 3. 创建工厂并注册构建器
	factoryBuilders := make(map[factory.Protocol]factory.Builder)

	// 注册默认构建器
	factoryBuilders[factory.Protocol(ProtocolOpenAI)] = builders.NewOpenAIBuilder()
	factoryBuilders[factory.Protocol(ProtocolArk)] = builders.NewArkBuilder()

	// 注册自定义构建器
	if cfg.CustomBuilders != nil {
		for protocol, builder := range cfg.CustomBuilders {
			factoryBuilders[factory.Protocol(protocol)] = &builderAdapter{builder: builder}
		}
	}

	fact := factory.NewFactory(factoryBuilders)

	// 4. 设置默认超时时间
	timeout := cfg.DefaultTimeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &defaultModelManager{
		mgr:     mgr,
		factory: fact,
		logger:  cfg.Logger,
		timeout: timeout,
	}, nil
}

// ListModels 列出符合条件的模型列表。
func (m *defaultModelManager) ListModels(ctx context.Context, opts ...*ListOptions) ([]*Model, error) {
	var opt *ListOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = &ListOptions{}
	}

	// 转换为内部选项
	internalOpts := &manager.ListOptions{
		Limit:      opt.Limit,
		Cursor:     opt.Cursor,
		NameFilter: opt.NameFilter,
	}

	// 转换状态
	if len(opt.Status) > 0 {
		internalOpts.Status = make([]int, len(opt.Status))
		for i, s := range opt.Status {
			internalOpts.Status[i] = int(s)
		}
	}

	// 查询内部管理器
	internalModels, err := m.mgr.List(ctx, internalOpts)
	if err != nil {
		return nil, err
	}

	// 转换为公共模型
	models := make([]*Model, len(internalModels))
	for i, im := range internalModels {
		models[i] = convertToModel(im)
	}

	return models, nil
}

// GetModel 根据 ID 获取单个模型的详细信息。
func (m *defaultModelManager) GetModel(ctx context.Context, id int64) (*Model, error) {
	if id <= 0 {
		return nil, ErrInvalidModelID
	}

	internalModel, err := m.mgr.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if internalModel == nil {
		return nil, ErrModelNotFound
	}

	return convertToModel(internalModel), nil
}

// GetModels 批量获取多个模型的信息。
func (m *defaultModelManager) GetModels(ctx context.Context, ids []int64) ([]*Model, error) {
	if len(ids) == 0 {
		return []*Model{}, nil
	}

	internalModels, err := m.mgr.MGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	models := make([]*Model, len(internalModels))
	for i, im := range internalModels {
		models[i] = convertToModel(im)
	}

	return models, nil
}

// CreateChatModel 创建一个聊天模型实例，用于实际的对话交互。
func (m *defaultModelManager) CreateChatModel(ctx context.Context, modelID int64, opts ...*ModelOptions) (ChatModel, error) {
	// 1. 获取模型元数据
	internalModel, err := m.mgr.Get(ctx, modelID)
	if err != nil {
		return nil, err
	}
	if internalModel == nil {
		return nil, ErrModelNotFound
	}

	// 2. 检查协议是否支持
	if !m.factory.Support(internalModel.Meta.Protocol) {
		return nil, NewError("PROTOCOL_NOT_SUPPORTED",
			fmt.Sprintf("protocol %s is not supported", internalModel.Meta.Protocol),
			ErrProtocolNotSupported)
	}

	// 3. 构建配置
	config := internalModel.Meta.ConnConfig
	if config == nil {
		return nil, NewError("INVALID_MODEL_CONFIG", "model connection config is missing", ErrInvalidConfig)
	}

	// 4. 合并运行时选项
	buildConfig := &BuildConfig{
		Protocol:    Protocol(internalModel.Meta.Protocol),
		Model:       config.Model,
		APIKey:      config.APIKey,
		BaseURL:     config.BaseURL,
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		Timeout:     m.timeout,
		Extra:       make(map[string]interface{}),
	}

	// 应用运行时选项
	if len(opts) > 0 && opts[0] != nil {
		opt := opts[0]
		if opt.Temperature != nil {
			buildConfig.Temperature = opt.Temperature
		}
		if opt.MaxTokens != nil {
			buildConfig.MaxTokens = opt.MaxTokens
		}
		if opt.TopP != nil {
			buildConfig.Extra["top_p"] = *opt.TopP
		}
		if opt.CustomConfig != nil {
			for k, v := range opt.CustomConfig {
				buildConfig.Extra[k] = v
			}
		}
	}

	// 5. 使用工厂创建模型实例
	factoryModel, err := m.factory.Create(ctx, internalModel.Meta.Protocol, &factory.Config{
		Model:       buildConfig.Model,
		APIKey:      buildConfig.APIKey,
		BaseURL:     buildConfig.BaseURL,
		Temperature: buildConfig.Temperature,
		MaxTokens:   buildConfig.MaxTokens,
		Timeout:     buildConfig.Timeout.Milliseconds(),
		Extra:       buildConfig.Extra,
	})
	if err != nil {
		return nil, NewError("MODEL_CREATION_FAILED",
			fmt.Sprintf("failed to create model instance: %v", err),
			ErrModelCreationFailed)
	}

	// 6. 返回适配器
	return &chatModelAdapter{model: factoryModel}, nil
}
