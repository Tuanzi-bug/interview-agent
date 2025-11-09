package manager

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// LoadFromYAML 从目录中加载所有 YAML 配置文件。
func LoadFromYAML(dir string) ([]*InternalModel, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var models []*InternalModel
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml") {
			path := filepath.Join(dir, entry.Name())
			model, err := loadYAMLFile(path)
			if err != nil {
				return nil, err
			}
			models = append(models, model)
		}
	}

	return models, nil
}

// LoadFromEnv 从模板文件和环境变量加载模型配置。
func LoadFromEnv(templatePath string) ([]*InternalModel, error) {
	// 首先从模板文件加载基础配置
	models, err := LoadFromYAML(templatePath)
	if err != nil {
		return nil, err
	}
	for _, model := range models {
		if model.Meta.ConnConfig != nil {
			envKeyName := strings.ToUpper(string(model.Meta.Protocol)) + "_API_KEY"
			if envKey := os.Getenv(envKeyName); envKey != "" {
				model.Meta.ConnConfig.APIKey = envKey
			}
			envURLName := strings.ToUpper(string(model.Meta.Protocol)) + "_BASE_URL"
			if envURL := os.Getenv(envURLName); envURL != "" {
				model.Meta.ConnConfig.BaseURL = envURL
			}
		}
	}

	return models, nil
}

// loadYAMLFile 加载单个 YAML 文件并解析为 InternalModel。
func loadYAMLFile(path string) (*InternalModel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var model InternalModel
	if err := yaml.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, nil
}
