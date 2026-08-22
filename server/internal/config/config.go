// Package config загружает конфигурацию приложения.
package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config — общая конфигурация приложения
type Config struct {
	AppEnv          string        `env:"APP_ENV" envDefault:"dev" validate:"oneof=dev staging prod"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s" validate:"gt=0"`
	GRPCServer      GRPCServerConfig
	Postgres        PostgresConfig
	JWT             JWTConfig
}

// GRPCServerConfig — конфигурация gRPC сервера
type GRPCServerConfig struct {
	Address      string        `env:"GRPC_ADDRESS" envDefault:":9090" validate:"required"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"5s" validate:"gt=0"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"5s" validate:"gt=0"`
	MaxMsgSize   int           `env:"SERVER_MAX_MSG_SIZE" envDefault:"4194304" validate:"gte=1024"`
	TLSCertFile  string        `env:"TLS_CERT_FILE" envDefault:"certs/server.crt"`
	TLSKeyFile   string        `env:"TLS_KEY_FILE" envDefault:"certs/server.key"`
}

// PostgresConfig — конфигурация подключения к PostgreSQL
type PostgresConfig struct {
	DSN string `env:"POSTGRES_DSN,required,notEmpty"`

	MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"25" validate:"gte=1"`
	MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"5" validate:"gte=0"`
	ConnMaxLifetime time.Duration `env:"POSTGRES_CONN_MAX_LIFETIME" envDefault:"5m" validate:"gt=0"`
	ConnMaxIdleTime time.Duration `env:"POSTGRES_CONN_MAX_IDLE_TIME" envDefault:"5m" validate:"gt=0"`

	ConnectTimeout time.Duration `env:"POSTGRES_CONNECT_TIMEOUT" envDefault:"10s" validate:"gt=0"`
	MaxRetries     int           `env:"POSTGRES_MAX_RETRIES" envDefault:"5" validate:"gte=1"`
	RetryInterval  time.Duration `env:"POSTGRES_RETRY_INTERVAL" envDefault:"3s" validate:"gt=0"`

	MigrateOnStart   bool          `env:"POSTGRES_MIGRATE_ON_START" envDefault:"true"`
	MigrationTimeout time.Duration `env:"POSTGRES_MIGRATION_TIMEOUT" envDefault:"30s" validate:"gt=0"`
}

// JWTConfig — конфигурация JWT токенов
type JWTConfig struct {
	Secret         string        `env:"JWT_SECRET,required" validate:"min=32"`
	AccessTokenTTL time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m" validate:"gt=0"`
	Issuer         string        `env:"JWT_ISSUER" envDefault:"credentials-vault" validate:"required"`
}

// Load — загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// IsProd — проверяет, является ли окружение production
func (c *Config) IsProd() bool {
	return c.AppEnv == "prod"
}

// IsStaging — проверяет, является ли окружение staging
func (c *Config) IsStaging() bool {
	return c.AppEnv == "staging"
}

// IsDev — проверяет, является ли окружение development
func (c *Config) IsDev() bool {
	return c.AppEnv == "dev"
}

// DSNRedacted — возвращает DSN без пароля для логирования
func (c PostgresConfig) DSNRedacted() string {
	u, err := url.Parse(c.DSN)
	if err != nil {
		return "invalid-dsn"
	}

	if u.User != nil {
		if _, has := u.User.Password(); has {
			u.User = url.UserPassword(u.User.Username(), "***")
		}
	}

	return u.String()
}
