// Package config 管理 gongfeng-cli 的凭据加载与持久化，支持四级优先链：
// CLI flags > 环境变量 > 当前目录 .gongfeng.json > 用户主目录 ~/.gongfeng.json
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 表示本地持久化的配置数据，存储于 .gongfeng.json
type Config struct {
	Token      string `json:"token,omitempty"`
	ProjectID  string `json:"project_id,omitempty"`
	BaseURL    string `json:"base_url,omitempty"`
}

// LoadConfig 按优先级加载配置：环境变量 > ./.gongfeng.json > ~/.gongfeng.json
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// 尝试从 ~/.gongfeng.json 加载
	homePath, err := GetConfigPath(false)
	if err == nil {
		if homeCfg, e := readConfigFile(homePath); e != nil {
			return nil, e
		} else if homeCfg != nil {
			cfg = homeCfg
		}
	}

	// 尝试从 ./.gongfeng.json 加载（优先级高于 home）
	localPath, _ := GetConfigPath(true)
	if localCfg, e := readConfigFile(localPath); e != nil {
		return nil, e
	} else if localCfg != nil {
		if localCfg.Token != "" {
			cfg.Token = localCfg.Token
		}
		if localCfg.ProjectID != "" {
			cfg.ProjectID = localCfg.ProjectID
		}
		if localCfg.BaseURL != "" {
			cfg.BaseURL = localCfg.BaseURL
		}
	}

	// 环境变量优先级最高
	if v := os.Getenv("GONGFENG_TOKEN"); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv("GONGFENG_PROJECT_ID"); v != "" {
		cfg.ProjectID = v
	}
	if v := os.Getenv("GONGFENG_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}

	return cfg, nil
}

// SaveConfig 将配置写入指定路径的 JSON 文件，自动创建父目录，文件权限 0600
func SaveConfig(cfg *Config, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0600)
}

// GetConfigPath 返回配置文件路径：local=true 返回 ./.gongfeng.json，否则返回 ~/.gongfeng.json
func GetConfigPath(local bool) (string, error) {
	if local {
		return ".gongfeng.json", nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gongfeng.json"), nil
}

// SaveProjectID 将 project_id 保存到当前目录的 .gongfeng.json，保留已有的其他字段
func SaveProjectID(projectID string) error {
	path, _ := GetConfigPath(true)
	cfg := &Config{}
	if existing, err := readConfigFile(path); err != nil {
		return err
	} else if existing != nil {
		cfg = existing
	}
	cfg.ProjectID = projectID
	return SaveConfig(cfg, path)
}

// readConfigFile 读取并解析指定路径的 .gongfeng.json 配置文件
// 文件不存在时返回 (nil, nil)，解析失败时返回错误
func readConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
