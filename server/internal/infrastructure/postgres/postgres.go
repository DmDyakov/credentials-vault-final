// Package postgres реализует подключение к PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"credentials-vault/server/internal/config"
	"credentials-vault/server/migrations"
)

// New создаёт подключение к PostgreSQL и применяет миграции.
func New(cfg config.PostgresConfig, logger *zap.Logger) (*sql.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("postgres: DSN is required")
	}

	// Применяем миграции через временное подключение
	if err := runMigrations(cfg, logger); err != nil {
		return nil, fmt.Errorf("postgres: failed to run migrations: %w", err)
	}

	// Основное подключение для приложения
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to open connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := pingWithRetry(db, cfg, logger); err != nil {
		db.Close()
		return nil, fmt.Errorf("postgres: failed to establish connection: %w", err)
	}

	logger.Info("postgres connection established")

	return db, nil
}

// runMigrations применяет миграции через отдельное подключение.
func runMigrations(cfg config.PostgresConfig, logger *zap.Logger) error {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to open migration connection: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping postgres for migrations: %w", err)
	}

	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	errCh := make(chan error, 1)
	go func() {
		errCh <- m.Up()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		if err == migrate.ErrNoChange {
			logger.Info("database schema is up to date")
		} else {
			logger.Info("migrations applied successfully")
		}
		return nil
	case <-time.After(cfg.MigrationTimeout):
		return fmt.Errorf("migration timeout after %s", cfg.MigrationTimeout)
	}
}

// pingWithRetry проверяет подключение с повторными попытками.
func pingWithRetry(db *sql.DB, cfg config.PostgresConfig, logger *zap.Logger) error {
	var lastErr error

	for attempt := 1; attempt <= cfg.MaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
		err := db.PingContext(ctx)
		cancel()

		if err == nil {
			return nil
		}

		lastErr = err
		logger.Warn("failed to ping postgres",
			zap.Int("attempt", attempt),
			zap.Int("max_retries", cfg.MaxRetries),
			zap.Error(err),
		)

		if attempt < cfg.MaxRetries {
			time.Sleep(cfg.RetryInterval)
		}
	}

	return fmt.Errorf("after %d attempts: %w", cfg.MaxRetries, lastErr)
}
