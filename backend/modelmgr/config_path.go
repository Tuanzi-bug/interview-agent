package modelmgr

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolveConfigPath 解析配置文件路径。
// 如果提供的路径为空，则按照预定义的优先级顺序查找配置目录。
func ResolveConfigPath(configPath string) (string, error) {
	// 如果已经提供了路径，直接使用
	if configPath != "" {
		// 如果是绝对路径，直接返回
		if filepath.IsAbs(configPath) {
			return configPath, nil
		}
		// 如果是相对路径，转换为绝对路径
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve config path: %w", err)
		}
		return absPath, nil
	}

	// 如果没有提供路径，按优先级查找
	searchPaths := []string{
		os.Getenv("MODELMGR_CONFIG_PATH"), // 1. 环境变量
		"./config",                        // 2. 当前目录下的 config
		"../config",                       // 3. 上级目录的 config（对 examples/ark 等示例很有用）
		"../../config",                    // 4. 上两级目录的 config
		"./modelmgr/config",               // 5. 当前目录下的 modelmgr/config
	}

	// 尝试查找项目根目录下的配置（通过 go.mod 定位）
	if projectRoot := findProjectRoot(); projectRoot != "" {
		searchPaths = append(searchPaths,
			filepath.Join(projectRoot, "backend/modelmgr/examples/config"), // 示例配置目录
			filepath.Join(projectRoot, "backend/modelmgr/config"),
			filepath.Join(projectRoot, "modelmgr/examples/config"),
			filepath.Join(projectRoot, "modelmgr/config"),
			filepath.Join(projectRoot, "config"),
		)
	}

	for _, path := range searchPaths {
		if path == "" {
			continue
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		// 检查路径是否存在
		if info, err := os.Stat(absPath); err == nil && info.IsDir() {
			return absPath, nil
		}
	}

	return "", fmt.Errorf("no valid config directory found. Please set MODELMGR_CONFIG_PATH environment variable or provide ConfigPath")
}

// findProjectRoot 查找项目根目录（包含 go.mod 的目录）。
func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// 向上查找，最多查找 10 层
	for i := 0; i < 10; i++ {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// 已经到根目录
			break
		}
		dir = parent
	}

	return ""
}

// GetDefaultConfigPath 获取默认配置路径（用于文档和示例）。
func GetDefaultConfigPath() string {
	if envPath := os.Getenv("MODELMGR_CONFIG_PATH"); envPath != "" {
		return envPath
	}

	if projectRoot := findProjectRoot(); projectRoot != "" {
		configPath := filepath.Join(projectRoot, "backend/modelmgr/config")
		if info, err := os.Stat(configPath); err == nil && info.IsDir() {
			return configPath
		}
	}
	return "./config"
}
