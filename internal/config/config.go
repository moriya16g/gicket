package config

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// ProjectConfig はリポジトリ共有設定 (.gicket/config.yml)
type ProjectConfig struct {
	Version         string   `yaml:"version"`
	DefaultPriority string   `yaml:"default_priority,omitempty"`
	DefaultLabels   []string `yaml:"default_labels,omitempty"`
	LabelCandidates []string `yaml:"label_candidates,omitempty"`
}

// UserConfig はユーザー個人設定 (~/.config/gicket/config.yml)
type UserConfig struct {
	Lang   string `yaml:"lang,omitempty"`
	Author string `yaml:"author,omitempty"`
}

// Config は統合された設定
type Config struct {
	Project ProjectConfig
	User    UserConfig
}

// Load はプロジェクト設定とユーザー設定を読み込んで統合する
func Load(gicketRoot string) *Config {
	cfg := &Config{
		Project: ProjectConfig{
			Version:         "1",
			DefaultPriority: "medium",
		},
	}

	// プロジェクト設定
	projectPath := filepath.Join(gicketRoot, ".gicket", "config.yml")
	if data, err := os.ReadFile(projectPath); err == nil {
		yaml.Unmarshal(data, &cfg.Project)
	}

	// ユーザー設定
	userPath := userConfigPath()
	if userPath != "" {
		if data, err := os.ReadFile(userPath); err == nil {
			yaml.Unmarshal(data, &cfg.User)
		}
	}

	return cfg
}

// SaveProject はプロジェクト設定を保存する
func SaveProject(gicketRoot string, cfg *ProjectConfig) error {
	configPath := filepath.Join(gicketRoot, ".gicket", "config.yml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// DefaultPriority は設定されたデフォルト優先度を返す
func (c *Config) DefaultPriority() string {
	if c.Project.DefaultPriority != "" {
		return c.Project.DefaultPriority
	}
	return "medium"
}

// Author は設定された author を返す（空の場合は空文字）
func (c *Config) Author() string {
	return c.User.Author
}

// Lang は設定された言語を返す（空の場合は空文字）
func (c *Config) LangSetting() string {
	return c.User.Lang
}

func userConfigPath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "gicket", "config.yml")
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, ".config", "gicket", "config.yml")
	}

	// XDG_CONFIG_HOME or ~/.config
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return filepath.Join(xdg, "gicket", "config.yml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "gicket", "config.yml")
}
