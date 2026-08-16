// Package app собирает и запускает сервер.
package app

import (
	"context"
	"fmt"
	"io"

	"credentials-vault/internal/config"
	"credentials-vault/internal/infrastructure/postgres"
	"credentials-vault/internal/transport/grpc"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// App управляет жизненным циклом сервера.
type App struct {
	cfg        *config.Config
	logger     *zap.Logger
	grpcServer *grpc.Server
	closers    []io.Closer
}

// New создаёт новый App.
func New(cfg *config.Config, logger *zap.Logger) (*App, error) {
	pg, err := postgres.New(cfg.Postgres, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres: %w", err)
	}

	grpcServer := grpc.NewServer(cfg, logger)

	return &App{
		cfg:        cfg,
		logger:     logger,
		grpcServer: grpcServer,
		closers:    []io.Closer{pg},
	}, nil
}

// Run запускает сервер.
func (a *App) Run(ctx context.Context) error {
	defer func() {
		for _, c := range a.closers {
			if err := c.Close(); err != nil {
				a.logger.Error("failed to close resource", zap.Error(err))
			}
		}
	}()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.grpcServer.Run(ctx)
	})

	return g.Wait()
}
