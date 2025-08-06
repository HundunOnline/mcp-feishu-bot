package config

import (
	"encoding/json"
	"fmt"
	"mcp-feishu/internal/types"
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Feishu types.FeishuConfig `json:"feishu"`
	Server ServerConfig       `json:"server"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.json"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() *Config {
	config := &Config{
		Feishu: types.FeishuConfig{
			WebhookURL:   os.Getenv("FEISHU_WEBHOOK_URL"),
			Secret:       os.Getenv("FEISHU_SECRET"),
			SecurityType: getEnvOrDefault("FEISHU_SECURITY_TYPE", "none"),
		},
		Server: ServerConfig{
			Port: getEnvAsIntOrDefault("SERVER_PORT", 3000),
			Host: getEnvOrDefault("SERVER_HOST", "localhost"),
		},
	}

	// 处理关键词
	keywords := os.Getenv("FEISHU_KEYWORDS")
	if keywords != "" {
		var keywordList []string
		if err := json.Unmarshal([]byte(keywords), &keywordList); err == nil {
			config.Feishu.Keywords = keywordList
		}
	}

	return config
}

// MergeConfig 合并配置，env配置优先，file配置作为fallback
func MergeConfig(envConfig, fileConfig *Config) *Config {
	merged := &Config{}

	// 合并飞书配置
	merged.Feishu.WebhookURL = envConfig.Feishu.WebhookURL
	if merged.Feishu.WebhookURL == "" {
		merged.Feishu.WebhookURL = fileConfig.Feishu.WebhookURL
	}

	merged.Feishu.Secret = envConfig.Feishu.Secret
	if merged.Feishu.Secret == "" {
		merged.Feishu.Secret = fileConfig.Feishu.Secret
	}

	merged.Feishu.SecurityType = envConfig.Feishu.SecurityType
	if merged.Feishu.SecurityType == "none" && fileConfig.Feishu.SecurityType != "" {
		merged.Feishu.SecurityType = fileConfig.Feishu.SecurityType
	}

	// 关键词优先使用环境变量，否则使用文件配置
	if len(envConfig.Feishu.Keywords) > 0 {
		merged.Feishu.Keywords = envConfig.Feishu.Keywords
	} else {
		merged.Feishu.Keywords = fileConfig.Feishu.Keywords
	}

	// 合并服务器配置
	merged.Server.Port = envConfig.Server.Port
	if merged.Server.Port == 3000 && fileConfig.Server.Port != 0 {
		merged.Server.Port = fileConfig.Server.Port
	}

	merged.Server.Host = envConfig.Server.Host
	if merged.Server.Host == "localhost" && fileConfig.Server.Host != "" {
		merged.Server.Host = fileConfig.Server.Host
	}

	return merged
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Feishu.WebhookURL == "" {
		return fmt.Errorf("飞书Webhook URL不能为空")
	}

	switch types.SecurityType(config.Feishu.SecurityType) {
	case types.SecurityTypeSignature:
		if config.Feishu.Secret == "" {
			return fmt.Errorf("签名校验模式下密钥不能为空")
		}
	case types.SecurityTypeKeyword:
		if len(config.Feishu.Keywords) == 0 {
			return fmt.Errorf("关键词模式下关键词列表不能为空")
		}
	case types.SecurityTypeNone:
		// 无安全设置，不需要验证
	default:
		return fmt.Errorf("不支持的安全类型: %s", config.Feishu.SecurityType)
	}

	return nil
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault 获取环境变量并转换为整数，失败时返回默认值
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		configPath = "config.json"
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	return nil
}
