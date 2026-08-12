// Package app собирает и запускает сервер.
package app

import (
	"context"

	"credentials-vault/internal/config"
	"credentials-vault/internal/transport/grpc"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// App управляет жизненным циклом сервера.
type App struct {
	cfg        *config.Config
	logger     *zap.Logger
	grpcServer *grpc.Server
}

// New создаёт новый App.
func New(cfg *config.Config, logger *zap.Logger) (*App, error) {
	grpcServer := grpc.NewServer(cfg, logger)

	return &App{
		cfg:        cfg,
		logger:     logger,
		grpcServer: grpcServer,
	}, nil
}

// Run запускает сервер.
func (a *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.grpcServer.Run(ctx)
	})

	return g.Wait()
}
