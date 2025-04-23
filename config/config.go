package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 系统配置
type Config struct {
	Push struct {
		CdrURL    string `yaml:"cdr_url"`
		StatusURL string `yaml:"status_url"`
		Workers   int    `yaml:"workers"` // 并发推送的工作协程数量
	} `yaml:"push"`

	Account struct {
		ID          string `yaml:"id"`
		ServiceType int    `yaml:"service_type"`
	} `yaml:"account"`

	Retry struct {
		Times  int   `yaml:"times"`
		Delays []int `yaml:"delays"`
	} `yaml:"retry"`

	Interval struct {
		CDR     int `yaml:"cdr"`
		Status  int `yaml:"status"`
		NewCall int `yaml:"new_call"`
	} `yaml:"interval"`
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证必要的配置项
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate 验证配置是否完整
func (c *Config) validate() error {
	if c.Push.CdrURL == "" {
		return fmt.Errorf("CDR推送地址未配置")
	}
	if c.Push.StatusURL == "" {
		return fmt.Errorf("状态推送地址未配置")
	}
	if c.Account.ID == "" {
		return fmt.Errorf("账号ID未配置")
	}
	if len(c.Retry.Delays) == 0 {
		return fmt.Errorf("重试间隔未配置")
	}
	return nil
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	// 优先使用环境变量
	if path := os.Getenv("CDR_CONFIG_PATH"); path != "" {
		return path
	}

	// 默认使用当前目录下的config/config.yaml
	return filepath.Join("config", "config.yaml")
}
