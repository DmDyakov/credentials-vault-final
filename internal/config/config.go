package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	AppEnv          string
	GRPCAddress     string
	DatabaseDSN     string
	JWTSecret       string
	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:          getEnv("APP_ENV", "dev"),
		GRPCAddress:     getEnv("GRPC_ADDRESS", ":9090"),
		DatabaseDSN:     os.Getenv("DATABASE_DSN"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}

	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func (c *Config) IsProd() bool {
	return c.AppEnv == "prod"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
