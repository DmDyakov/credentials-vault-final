// Package app собирает и запускает сервер.
package app

import (
	"context"
	"fmt"

	"credentials-vault/pkg/jwt"
	"credentials-vault/server/internal/config"
	"credentials-vault/server/internal/infrastructure/postgres"
	"credentials-vault/server/internal/repository"
	"credentials-vault/server/internal/service/user"
	"credentials-vault/server/internal/service/vault"
	"credentials-vault/server/internal/transport/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// App управляет жизненным циклом сервера.
type App struct {
	cfg        *config.Config
	logger     *zap.Logger
	grpcServer *grpc.Server
	pg         *pgxpool.Pool
}

// New создаёт новый App.
func New(cfg *config.Config, logger *zap.Logger) (*App, error) {
	pg, err := postgres.New(cfg.Postgres, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres: %w", err)
	}

	userRepo := repository.NewUserRepository(pg)
	vaultRepo := repository.NewVaultItemRepository(pg)

	userService := user.NewService(userRepo)
	vaultService := vault.NewService(vaultRepo)

	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)
	grpcServer, err := grpc.NewServer(cfg, logger, userService, vaultService, jwtManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC server: %w", err)
	}

	return &App{
		cfg:        cfg,
		logger:     logger,
		grpcServer: grpcServer,
		pg:         pg,
	}, nil
}

// Run запускает сервер.
func (a *App) Run(ctx context.Context) error {
	defer a.pg.Close()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.grpcServer.Run(ctx)
	})

	return g.Wait()
}
