// Package config tool/config.go（配置结构体定义）
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// 1. 根配置结构体：对应 config.yaml 整个文件
// 类比 Spring Boot 中的配置类（用 @Configuration 注解的类）
type Config struct {
	Google GoogleConfig `yaml:"google"` // 对应 YAML 中的 google 节点
	OpenAI OpenAIConfig `yaml:"openai"` // 对应 YAML 中的 openai 节点（可选）
}

// 2. 谷歌配置结构体：对应 YAML 中的 google 节点
type GoogleConfig struct {
	APIKey         string `yaml:"api_key"`          // 对应 google.api_key
	SearchEngineID string `yaml:"search_engine_id"` // 对应 google.search_engine_id
}

// 3. OpenAI 配置结构体（可选，对应 YAML 中的 openai 节点）
type OpenAIConfig struct {
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model_name"`
	BaseURL string `yaml:"base_url"`
}

// LoadConfig：读取 config.yaml 文件，返回配置结构体
// 类比 Spring Boot 的配置自动加载（手动实现，但能复用）
func LoadConfig() (*Config, error) {
	// 1. 获取 config.yaml 文件的绝对路径（避免相对路径错误）
	// 项目根目录下的 config.yaml，这里用 filepath.Abs 转成绝对路径
	configPath, err := filepath.Abs("../../config.yaml")
	if err != nil {
		return nil, fmt.Errorf("获取配置文件路径失败：%v", err)
	}

	// 2. 打开 config.yaml 文件（类似 Spring Boot 读取 application.yml）
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败（路径：%s）：%v", configPath, err)
	}
	defer file.Close() // 函数结束后自动关闭文件，避免资源泄漏

	// 3. 解析 YAML 文件内容到 Config 结构体（核心步骤）
	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析 YAML 配置失败：%v", err)
	}

	// 4. 检查配置是否为空（可选，避免密钥没填）
	if config.Google.APIKey == "" || config.Google.SearchEngineID == "" {
		return nil, fmt.Errorf("config.yaml 中 google.api_key 或 google.search_engine_id 未配置")
	}

	return &config, nil
}
