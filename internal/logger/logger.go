// Package logger предоставляет функции для логирования.
package logger

import (
	"go.uber.org/zap"

	"credentials-vault/internal/config"
)

// New создаёт логгер в зависимости от окружения.
func New(cfg *config.Config) (*zap.Logger, error) {
	if cfg.IsProd() {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
