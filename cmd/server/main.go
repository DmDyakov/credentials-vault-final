package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"credentials-vault/internal/app"
	"credentials-vault/internal/config"
	"credentials-vault/internal/logger"
	"credentials-vault/pkg/buildinfo"
	"credentials-vault/pkg/lifecycle"

	"go.uber.org/zap"
)

func main() {
	buildinfo.Print()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to sync logger: %v\n", err)
		}
	}()

	app, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal("failed to create app: %v", zap.Error(err))
	}

	if err := lifecycle.Run(ctx, app, cfg.ShutdownTimeout); err != nil {
		logger.Fatal("app terminated with error", zap.Error(err))
	}
}
