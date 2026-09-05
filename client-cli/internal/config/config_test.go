package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("config not exists", func(t *testing.T) {
		t.Run("config not exists", func(t *testing.T) {
			tmpHome := t.TempDir()
			t.Setenv("USERPROFILE", tmpHome)
			t.Setenv("HOME", tmpHome)

			cfg, err := New()

			assert.NoError(t, err)
			assert.NotNil(t, cfg)
			assert.False(t, cfg.Valid())
		})
	})

	t.Run("invalid config", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("USERPROFILE", tmpHome)
		t.Setenv("HOME", tmpHome)

		configDir := filepath.Join(tmpHome, ".credentials-vault")
		err := os.MkdirAll(configDir, 0700)
		assert.NoError(t, err)

		configPath := filepath.Join(configDir, "config.json")
		err = os.WriteFile(configPath, []byte("invalid json"), 0600)
		assert.NoError(t, err)

		_, err = New()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse config")
	})

	t.Run("missing server_address", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("USERPROFILE", tmpHome)
		t.Setenv("HOME", tmpHome)

		configDir := filepath.Join(tmpHome, ".credentials-vault")
		err := os.MkdirAll(configDir, 0700)
		assert.NoError(t, err)

		configPath := filepath.Join(configDir, "config.json")
		err = os.WriteFile(configPath, []byte(`{"ca_file":"certs/server.crt"}`), 0600)
		assert.NoError(t, err)

		cfg, err := New()
		assert.NoError(t, err)
		assert.False(t, cfg.Valid())
	})

	t.Run("valid config", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("USERPROFILE", tmpHome)
		t.Setenv("HOME", tmpHome)

		configDir := filepath.Join(tmpHome, ".credentials-vault")
		err := os.MkdirAll(configDir, 0700)
		assert.NoError(t, err)

		configPath := filepath.Join(configDir, "config.json")
		err = os.WriteFile(configPath, []byte(`{"server_address":"remote:9090","token":"test-token"}`), 0600)
		assert.NoError(t, err)

		cfg, err := New()
		assert.NoError(t, err)
		assert.Equal(t, "remote:9090", cfg.ServerAddress)
		assert.Equal(t, "test-token", cfg.Token)
	})
}

func TestSave(t *testing.T) {
	t.Run("save creates file", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("USERPROFILE", tmpHome)
		t.Setenv("HOME", tmpHome)

		cfg, err := New()
		assert.NoError(t, err)

		cfg.ServerAddress = "remote:9090"
		cfg.Token = "test-token"

		err = cfg.Save()
		assert.NoError(t, err)

		_, err = os.Stat(cfg.path)
		assert.NoError(t, err)

		data, err := os.ReadFile(cfg.path)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "remote:9090")
		assert.Contains(t, string(data), "test-token")
	})

	t.Run("save without path", func(t *testing.T) {
		cfg := &Config{
			ServerAddress: "localhost:9090",
		}

		err := cfg.Save()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config path is not set")
	})

	t.Run("save and load roundtrip", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("USERPROFILE", tmpHome)
		t.Setenv("HOME", tmpHome)

		cfg, err := New()
		assert.NoError(t, err)

		cfg.ServerAddress = "remote:9090"
		cfg.Token = "test-token"

		err = cfg.Save()
		assert.NoError(t, err)

		loaded, err := New()
		assert.NoError(t, err)
		assert.Equal(t, "remote:9090", loaded.ServerAddress)
		assert.Equal(t, "test-token", loaded.Token)
	})
}
