package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用程序配置结构
type Config struct {
	Host      string         `yaml:"host"`
	Port      int            `yaml:"port"`
	Database  DatabaseConfig `yaml:"database"`
	Redis     RedisConfig    `yaml:"redis"`
	Hertz     HertzConfig    `yaml:"hertz"`
	Eino      EinoConfig     `yaml:"eino"`
	Interview InterviewConfig `yaml:"interview"`
	Security  SecurityConfig `yaml:"security"`
	GoogleSearch GoogleConfig `yaml:"google_search"`
	OpenAI    OpenAIConfig   `yaml:"openai"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver         string `yaml:"driver"`
	DSN            string `yaml:"dsn"`
	MaxOpenConns   int    `yaml:"max_open_conns"`
	MaxIdleConns   int    `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	DialTimeout  string `yaml:"dial_timeout"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

// HertzConfig Hertz框架配置
type HertzConfig struct {
	LogLevel    string `yaml:"log_level"`
	LogPath     string `yaml:"log_path"`
	ReadTimeout string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout string `yaml:"idle_timeout"`
}

// EinoConfig Eino框架配置
type EinoConfig struct {
	Model       string  `yaml:"model"`
	APIKey      string  `yaml:"api_key"`
	BaseURL     string  `yaml:"base_url"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
	RetryCount  int     `yaml:"retry_count"`
	RetryDelay  string  `yaml:"retry_delay"`
}

// InterviewConfig 面试系统配置
type InterviewConfig struct {
	MaxDuration    string `yaml:"max_duration"`
	QuestionTimeout string `yaml:"question_timeout"`
	MaxQuestions   int    `yaml:"max_questions"`
	MinQuestions   int    `yaml:"min_questions"`
}

// SecurityConfig 安全性配置
type SecurityConfig struct {
	JWTSecret      string   `yaml:"jwt_secret"`
	JWTExpiration  string   `yaml:"jwt_expiration"`
	CORS           CORSConfig `yaml:"cors"`
}

// GoogleConfig Google搜索配置
type GoogleConfig struct {
	APIKey       string `yaml:"api_key"`
	SearchEngineID string `yaml:"search_engine_id"`
}

// OpenAIConfig OpenAI配置
type OpenAIConfig struct {
	APIKey    string `yaml:"api_key"`
	ModelName string `yaml:"model_name"`
	BaseURL   string `yaml:"base_url"`
}

// Global 全局配置实例
var Global Config

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	Global = cfg
	log.Println("配置加载成功")
	return &cfg, nil
}