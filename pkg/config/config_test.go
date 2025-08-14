package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Run("should load default config when no file or env vars are provided", func(t *testing.T) {
		os.Unsetenv("PORT")
		os.Unsetenv("CONFIG_DIR")
		cfg, err := NewConfig()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if cfg.Port != "8080" {
			t.Errorf("expected port 8080, got %s", cfg.Port)
		}
	})

	t.Run("should load config from env vars", func(t *testing.T) {
		os.Unsetenv("PORT")
		os.Unsetenv("CONFIG_DIR")
		os.Setenv("PORT", "8888")
		defer os.Unsetenv("PORT")

		cfg, err := NewConfig()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if cfg.Port != "8888" {
			t.Errorf("expected port 8888, got %s", cfg.Port)
		}
	})

	t.Run("should load config from file", func(t *testing.T) {
		os.Unsetenv("PORT")
		os.Unsetenv("CONFIG_DIR")
		configDir := "./.test_config"
		configFile := filepath.Join(configDir, "config.json")

		os.MkdirAll(configDir, 0o755)
		defer os.RemoveAll(configDir)

		file, err := os.Create(configFile)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		file.WriteString(`{"port":"9999"}`)
		file.Close()

		os.Setenv("CONFIG_DIR", configDir)
		defer os.Unsetenv("CONFIG_DIR")

		cfg, err := NewConfig()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if cfg.Port != "9999" {
			t.Errorf("expected port 9999, got %s", cfg.Port)
		}
	})
}