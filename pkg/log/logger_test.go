package logger

import (
	"downloader/pkg/config"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Run("should create a logger with valid config", func(t *testing.T) {
		logDir := "./.test_logs"
		cfg := &config.Config{LogDir: logDir}

		logger := NewLogger(cfg)

		if logger == nil {
			t.Error("expected logger to be not nil, got nil")
		}

		// Clean up the created log directory
		defer os.RemoveAll(logDir)
	})

	t.Run("should return console logger when log directory cannot be created", func(t *testing.T) {
		// Create a file with the same name as the log directory to cause an error
		logDir := "./.test_logs"
		filePath := filepath.Join(logDir, "app.log")
		os.MkdirAll(logDir, 0o755)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		file.Close()

		cfg := &config.Config{LogDir: filePath}

		logger := NewLogger(cfg)

		if logger == nil {
			t.Error("expected logger to be not nil, got nil")
		}

		// Clean up the created file and directory
		defer os.RemoveAll(logDir)
	})
}
