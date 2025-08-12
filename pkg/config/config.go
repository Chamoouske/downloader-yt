package config

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Port     string `json:"port"`
	LogDir   string `json:"log_dir"`
	VideoDir string `json:"video_dir"`
}

var AppConfig Config

func LoadConfig() error {
	AppConfig.Port = os.Getenv("PORT")
	if AppConfig.Port == "" {
		AppConfig.Port = "8080"
	}

	AppConfig.LogDir = os.Getenv("LOG_DIR")
	if AppConfig.LogDir == "" {
		AppConfig.LogDir = "./logs"
	}

	AppConfig.VideoDir = os.Getenv("VIDEO_DIR")
	if AppConfig.VideoDir == "" {
		AppConfig.VideoDir = "./videos"
	}

	return nil
}

func init() {
	LoadConfig()
	if err := os.MkdirAll(AppConfig.LogDir, 0o755); err != nil {
		slog.Error(fmt.Sprintf("error creating log directory: %s", err))
	}

	if err := os.MkdirAll(AppConfig.VideoDir, 0o755); err != nil {
		slog.Error(fmt.Sprintf("error creating video directory: %s", err))
	}
}

func GetConfig() Config {
	return AppConfig
}
