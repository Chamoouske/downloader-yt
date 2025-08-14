package config

import (
	"downloader/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Port       string         `json:"port"`
	LogDir     string         `json:"log_dir"`
	VideoDir   string         `json:"video_dir"`
	ConfigDir  string         `json:"config_dir"`
	URLWebhook string         `json:"url_webhook"`
	Database   ConfigDatabase `json:"db"`
}

type ConfigDatabase struct {
	Port string `json:"port"`
	URL  string `json:"url"`
	User string `json:"usr"`
	Psw  string `json:"psw"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	cfg.ConfigDir = utils.GetEnvOrDefault("CONFIG_DIR", "./.config")

	if err := cfg.readConfig(); err != nil {
		return nil, err
	}
	cfg.setDefaults()

	if err := cfg.createDirs(); err != nil {
		return nil, err
	}
	return cfg, cfg.saveConfig()
}

func (c *Config) GetPort() string {
	return c.Port
}

func (c *Config) GetLogDir() string {
	return c.LogDir
}

func (c *Config) GetVideoDir() string {
	return c.VideoDir
}

func (c *Config) GetConfigDir() string {
	return c.ConfigDir
}

func (c *Config) GetURLWebhook() string {
	return c.URLWebhook
}

func (c *Config) GetDbConfig() ConfigDatabase {
	return c.Database
}

func (c *Config) setDefaults() {
	if c.Port == "" {
		c.Port = utils.GetEnvOrDefault("PORT", "8080")
	}
	if c.LogDir == "" {
		if c.ConfigDir != "./.config" && c.ConfigDir != "" { // Check if ConfigDir is not default or empty
			c.LogDir = utils.GetEnvOrDefault("LOG_DIR", filepath.Join(c.ConfigDir, "..", ".logs"))
		} else {
			c.LogDir = utils.GetEnvOrDefault("LOG_DIR", "./.logs")
		}
	}
	if c.VideoDir == "" {
		if c.ConfigDir != "./.config" && c.ConfigDir != "" { // Check if ConfigDir is not default or empty
			c.VideoDir = utils.GetEnvOrDefault("VIDEO_DIR", filepath.Join(c.ConfigDir, "..", "videos"))
		} else {
			c.VideoDir = utils.GetEnvOrDefault("VIDEO_DIR", "./videos")
		}
	}
	if c.URLWebhook == "" {
		c.URLWebhook = utils.GetEnvOrDefault("WEBHOOK", "http://host.docker.internal:5677/webhook/downloader-yt")
	}
}

func (c *Config) readConfig() error {
	file, err := os.ReadFile(filepath.Join(c.ConfigDir, "config.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}
	err = json.Unmarshal(file, c)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

func (c *Config) createDirs() error {
	if err := os.MkdirAll(c.LogDir, 0o755); err != nil {
		return fmt.Errorf("error creating log directory: %w", err)
	}
	if err := os.MkdirAll(c.VideoDir, 0o755); err != nil {
		return fmt.Errorf("error creating video directory: %w", err)
	}
	if err := os.MkdirAll(c.ConfigDir, 0o755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}
	return nil
}

func (c *Config) saveConfig() error {
	file, err := os.Create(filepath.Join(c.ConfigDir, "config.json"))
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config to JSON: %w", err)
	}

	return nil
}
