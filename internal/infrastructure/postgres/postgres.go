package postgres

import (
	"context"
	"credentials-vault/internal/config"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// New — создает подключение к PostgreSQL
func New(cfg config.PostgresConfig, logger *zap.Logger) (*sql.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("postgres: DSN is required")
	}

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

	logger.Info("applying postgres migrations")
	if err := runMigrations(db, cfg.MigrationTimeout, logger); err != nil {
		db.Close()
		return nil, fmt.Errorf("postgres: failed to run migrations: %w", err)
	}

	return db, nil
}

// pingWithRetry — проверка подключения с повторными попытками
func pingWithRetry(db *sql.DB, cfg config.PostgresConfig, logger *zap.Logger) error {
	var lastErr error

	for attempt := 1; attempt <= cfg.MaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
		err := db.PingContext(ctx)
		cancel()

		if err == nil {
			logger.Info("postgres connection established",
				zap.Int("attempt", attempt),
			)
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

// runMigrations — применяет миграции
func runMigrations(db *sql.DB, timeout time.Duration, logger *zap.Logger) error {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	source, err := iofs.New(migrationsFS, "migrations")
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
	case <-time.After(timeout):
		return fmt.Errorf("migration timeout after %s", timeout)
	}
}
