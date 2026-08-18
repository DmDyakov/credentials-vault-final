package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"credentials-vault/pkg/buildinfo"
	"credentials-vault/pkg/lifecycle"
	"credentials-vault/server/internal/app"
	"credentials-vault/server/internal/config"
	"credentials-vault/server/internal/logger"
)

func run() error {
	buildinfo.Print()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger, err := logger.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			fmt.Fprintf(os.Stderr, "failed to sync logger: %v\n", syncErr)
		}
	}()

	app, err := app.New(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	if err := lifecycle.Run(ctx, app, cfg.ShutdownTimeout); err != nil {
		return fmt.Errorf("app terminated with error: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
