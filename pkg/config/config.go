package config

import (
	"downloader/pkg/utils"
	"encoding/json"
	"fmt"
	"log/slog"
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

var appConfig Config
var dbConfig ConfigDatabase

func GetConfig() Config {
	return appConfig
}

func GetDbConfig() ConfigDatabase {
	return dbConfig
}

func LoadConfig() error {
	readConfig()
	if appConfig.ConfigDir == "" {
		appConfig.ConfigDir = utils.GetEnvOrDefault("CONFIG_DIR", "./.config")
	}

	if appConfig.Port == "" {
		appConfig.Port = utils.GetEnvOrDefault("PORT", "8080")
	}

	if appConfig.LogDir == "" {
		appConfig.LogDir = utils.GetEnvOrDefault("LOG_DIR", "./.logs")
	}

	if appConfig.VideoDir == "" {
		appConfig.VideoDir = utils.GetEnvOrDefault("VIDEO_DIR", "./videos")
	}

	if appConfig.URLWebhook == "" {
		appConfig.URLWebhook = utils.GetEnvOrDefault("WEBHOOK", "http://host.docker.internal:5677/webhook/downloader-yt")
	}
	return nil
}

func init() {
	LoadConfig()
	if err := os.MkdirAll(appConfig.LogDir, 0o755); err != nil {
		slog.Error(fmt.Sprintf("error creating log directory: %s", err))
	}

	if err := os.MkdirAll(appConfig.VideoDir, 0o755); err != nil {
		slog.Error(fmt.Sprintf("error creating video directory: %s", err))
	}

	if err := os.MkdirAll(appConfig.ConfigDir, 0o755); err != nil {
		slog.Error(fmt.Sprintf("error creating config directory: %s", err))
	}
	saveConfig()
}

func readConfig() error {
	file, err := os.ReadFile(filepath.Join(utils.GetEnvOrDefault("CONFIG_DIR", "./.config"), "config.json"))
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	err = json.Unmarshal(file, &appConfig)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}

func saveConfig() error {
	file, err := os.Create(filepath.Join(appConfig.ConfigDir, "config.json"))
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(appConfig); err != nil {
		return fmt.Errorf("failed to encode config to JSON: %w", err)
	}

	return nil
}
