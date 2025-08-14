package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("should load default config when no file or env vars are provided", func(t *testing.T) {
		os.Unsetenv("PORT")
		os.Unsetenv("LOG_DIR")
		os.Unsetenv("VIDEO_DIR")
		os.Unsetenv("WEBHOOK")
		os.Unsetenv("CONFIG_DIR")

		cfg, err := NewConfig()

		assert.NoError(t, err)
		assert.Equal(t, "8080", cfg.GetPort())
		assert.Equal(t, "./.logs", cfg.GetLogDir())
		assert.Equal(t, "./videos", cfg.GetVideoDir())
		assert.Equal(t, "http://host.docker.internal:5677/webhook/downloader-yt", cfg.GetURLWebhook())
	})

	t.Run("should load config from env vars", func(t *testing.T) {
		os.Setenv("PORT", "8888")
		os.Setenv("LOG_DIR", "/custom/logs")
		os.Setenv("VIDEO_DIR", "/custom/videos")
		os.Setenv("WEBHOOK", "http://custom-webhook")
		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("LOG_DIR")
			os.Unsetenv("VIDEO_DIR")
			os.Unsetenv("WEBHOOK")
		}()

		cfg := &Config{}
		cfg.setDefaults()

		assert.Equal(t, "8888", cfg.GetPort())
		assert.Equal(t, "/custom/logs", cfg.GetLogDir())
		assert.Equal(t, "/custom/videos", cfg.GetVideoDir())
		assert.Equal(t, "http://custom-webhook", cfg.GetURLWebhook())
	})

	t.Run("should load config from file with database config", func(t *testing.T) {
		configDir := t.TempDir()
		configFile := filepath.Join(configDir, "config.json")

		testConfig := `{
			"port": "9999",
			"log_dir": "./file_logs",
			"video_dir": "./file_videos",
			"db": {
				"port": "5432",
				"url": "db.example.com",
				"usr": "dbuser",
				"psw": "dbpass"
			}
		}`

		os.WriteFile(configFile, []byte(testConfig), 0o644)
		os.Setenv("CONFIG_DIR", configDir)
		defer os.Unsetenv("CONFIG_DIR")

		cfg, err := NewConfig()

		assert.NoError(t, err)
		assert.Equal(t, "9999", cfg.GetPort())
		assert.Equal(t, "./file_logs", cfg.GetLogDir())
		assert.Equal(t, "./file_videos", cfg.GetVideoDir())
		db := cfg.GetDbConfig()
		assert.Equal(t, "5432", db.Port)
		assert.Equal(t, "db.example.com", db.URL)
		assert.Equal(t, "dbuser", db.User)
		assert.Equal(t, "dbpass", db.Psw)
	})

	t.Run("should handle invalid config file", func(t *testing.T) {
		configDir := t.TempDir()
		configFile := filepath.Join(configDir, "config.json")

		os.WriteFile(configFile, []byte("{invalid json}"), 0o644)
		os.Setenv("CONFIG_DIR", configDir)
		defer os.Unsetenv("CONFIG_DIR")

		_, err := NewConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load config")
	})
}

func TestConfigMethods(t *testing.T) {
	cfg := &Config{
		Port:       "3000",
		LogDir:     "./logs",
		VideoDir:   "./videos",
		ConfigDir:  "./config",
		URLWebhook: "http://webhook",
		Database: ConfigDatabase{
			Port: "3306",
			URL:  "localhost",
			User: "user",
			Psw:  "pass",
		},
	}

	assert.Equal(t, "3000", cfg.GetPort())
	assert.Equal(t, "./logs", cfg.GetLogDir())
	assert.Equal(t, "./videos", cfg.GetVideoDir())
	assert.Equal(t, "./config", cfg.GetConfigDir())
	assert.Equal(t, "http://webhook", cfg.GetURLWebhook())
	db := cfg.GetDbConfig()
	assert.Equal(t, "3306", db.Port)
	assert.Equal(t, "localhost", db.URL)
	assert.Equal(t, "user", db.User)
	assert.Equal(t, "pass", db.Psw)

	// Test empty config
	emptyCfg := &Config{}
	assert.Equal(t, "", emptyCfg.GetPort())
	assert.Equal(t, "", emptyCfg.GetLogDir())
	assert.Equal(t, "", emptyCfg.GetVideoDir())
	assert.Equal(t, "", emptyCfg.GetURLWebhook())
}

func TestSaveConfig(t *testing.T) {
	configDir := t.TempDir()
	os.Setenv("CONFIG_DIR", configDir)
	defer os.Unsetenv("CONFIG_DIR")

	cfg := &Config{
		Port:      "3000",
		ConfigDir: configDir,
	}

	err := cfg.saveConfig()
	assert.NoError(t, err)

	// Verify file was created
	configFile := filepath.Join(configDir, "config.json")
	_, err = os.Stat(configFile)
	assert.NoError(t, err)

	// Test save error (invalid directory)
	invalidCfg := &Config{ConfigDir: "/invalid/path"}
	err = invalidCfg.saveConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create config file")
}
