// Package grpc реализует gRPC-сервер Credentials Vault.
package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	authpb "credentials-vault/gen/go/auth/v1"
	vaultpb "credentials-vault/gen/go/vault/v1"
	"credentials-vault/internal/config"
	"credentials-vault/internal/domain"
	"credentials-vault/internal/transport/grpc/handler"
	"credentials-vault/internal/transport/grpc/interceptor"
	"credentials-vault/pkg/jwt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserService interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.User, error)
}

type VaultService interface {
	CreateItem(ctx context.Context, userID uuid.UUID, itemType domain.ItemType, encryptedData []byte, metadata map[string]string) (*domain.VaultItem, error)
	GetItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.VaultItem, error)
	ListItems(ctx context.Context, userID uuid.UUID, itemType *domain.ItemType) ([]*domain.VaultItem, error)
	UpdateItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, encryptedData []byte, metadata map[string]string) (*domain.VaultItem, error)
	DeleteItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type JWTManager interface {
	Generate(userID string) (jwt.Token, time.Time, error)
	Verify(token jwt.Token) (*jwt.Claims, error)
}

type Server struct {
	*grpc.Server
	addr   string
	logger *zap.Logger
}

// NewServer создаёт gRPC-сервер.
func NewServer(cfg *config.Config, logger *zap.Logger, userService UserService, vaultService VaultService, jwtManager JWTManager) *Server {
	authInterceptor := interceptor.NewAuthInterceptor(jwtManager)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	authHandler := handler.NewAuthHandler(userService)
	authpb.RegisterAuthServiceServer(s, authHandler)

	vaultHandler := handler.NewVaultHandler(vaultService)
	vaultpb.RegisterVaultServiceServer(s, vaultHandler)

	if cfg.IsDev() {
		reflection.Register(s)
		logger.Info("gRPC reflection enabled (dev mode)")
	}

	return &Server{
		Server: s,
		addr:   cfg.GRPCServer.Address,
		logger: logger,
	}
}

// Run запускает сервер.
func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("gRPC server failed to listen: %w", err)
	}

	s.logger.Info("gRPC server started", zap.String("addr", s.addr))

	go func() {
		<-ctx.Done()
		s.logger.Info("Shutting down gRPC server...")
		s.GracefulStop()
	}()

	if err := s.Serve(listener); err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("gRPC server error: %w", err)
	}

	s.logger.Info("gRPC server stopped gracefully")
	return nil
}
