package milvus

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

//参考文档 https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/retriever/retriever_milvus/

// RetrieverService 封装 Milvus 检索器服务
type RetrieverService struct {
	retriever *milvus.Retriever
	config    *milvus.RetrieverConfig
}

// NewRetrieverService 创建新的检索器服务
func NewRetrieverService(ctx context.Context, config *milvus.RetrieverConfig) (*RetrieverService, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if config.Client == nil {
		return nil, fmt.Errorf("milvus client is nil")
	}
	if config.Collection == "" {
		return nil, fmt.Errorf("collection name is empty")
	}
	if config.Embedding == nil {
		return nil, fmt.Errorf("embedding is nil")
	}
	// 设置默认值
	if config.VectorField == "" {
		config.VectorField = "vector"
	}
	if len(config.OutputFields) == 0 {
		config.OutputFields = []string{"id", "content", "metadata"}
	}
	if config.TopK <= 0 {
		config.TopK = 10
	}
	// 构建检索器配置
	// 注意：需要显式设置 Sp（搜索参数）和 VectorConverter
	// 1. eino-ext 的 defaultSearchParam 函数错误地将维度值作为 radius 参数传递
	//    参考：https://github.com/cloudwego/eino-ext/blob/main/components/retriever/milvus/utils.go#L40
	// 2. eino-ext 的 defaultVectorConverter 返回 BinaryVector，但我们使用 FloatVector
	//    参考：https://github.com/cloudwego/eino-ext/blob/main/components/retriever/milvus/utils.go#L97
	retrieverConfig := &milvus.RetrieverConfig{
		Client:            config.Client,
		Collection:        config.Collection,
		Partition:         config.Partition,
		VectorField:       config.VectorField,
		OutputFields:      config.OutputFields,
		DocumentConverter: config.DocumentConverter,
		// 设置 VectorConverter 为 FloatVector 转换器（与 indexer 保持一致）
		VectorConverter: floatVectorConverter,
		MetricType:      config.MetricType,
		TopK:            config.TopK,
		Embedding:       config.Embedding,
	}
	if config.Sp == nil {
		// 创建一个简单的 AUTOINDEX 搜索参数，不使用 range search
		searchParam, err := entity.NewIndexAUTOINDEXSearchParam(1)
		if err != nil {
			return nil, fmt.Errorf("failed to create search param: %w", err)
		}
		retrieverConfig.Sp = searchParam
	} else {
		retrieverConfig.Sp = config.Sp
	}
	// 只有当 ScoreThreshold > 0 时才设置
	if config.ScoreThreshold > 0 {
		retrieverConfig.ScoreThreshold = config.ScoreThreshold
	}
	retriever, err := milvus.NewRetriever(ctx, retrieverConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create milvus retriever: %w", err)
	}

	return &RetrieverService{
		retriever: retriever,
		config:    config,
	}, nil
}

// Retrieve 检索相关文档
func (s *RetrieverService) Retrieve(ctx context.Context, query string) ([]*schema.Document, error) {
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}
	documents, err := s.retriever.Retrieve(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %w", err)
	}
	return documents, nil
}

// GetConfig 获取配置信息
func (s *RetrieverService) GetConfig() *milvus.RetrieverConfig {
	return s.config
}

// floatVectorConverter 将 float64 向量转换为 FloatVector
// 这个转换器与 indexer 中使用的向量类型保持一致
func floatVectorConverter(ctx context.Context, vectors [][]float64) ([]entity.Vector, error) {
	result := make([]entity.Vector, 0, len(vectors))
	for _, vector := range vectors {
		float32Vec := make([]float32, len(vector))
		for i, v := range vector {
			float32Vec[i] = float32(v)
		}
		result = append(result, entity.FloatVector(float32Vec))
	}
	return result, nil
}
