// Package config содержит конфигурацию клиента.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config - конфигурация клиента.
type Config struct {
	path          string `json:"-"`
	ServerAddress string `json:"server_address"`
	Token         string `json:"token,omitempty"`
	CAFile        string `json:"ca_file,omitempty"`
}

// New создаёт конфигурацию, загружая её из файла.
func New() (*Config, error) {
	cfg := &Config{}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cfg.path = filepath.Join(home, ".credentials-vault", "config.json")

	data, err := os.ReadFile(cfg.path)
	if err != nil {
		if os.IsNotExist(err) {
			if saveErr := cfg.Save(); saveErr != nil {
				return nil, fmt.Errorf("failed to create config: %w", saveErr)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// Save сохраняет конфигурацию в файл.
func (c *Config) Save() error {
	if c.path == "" {
		return fmt.Errorf("config path is not set")
	}

	if err := os.MkdirAll(filepath.Dir(c.path), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Valid проверяет, что конфиг заполнен.
func (c *Config) Valid() bool {
	return c != nil && c.ServerAddress != ""
}
