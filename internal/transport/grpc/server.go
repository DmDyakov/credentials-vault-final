// Package grpc реализует gRPC-сервер Credentials Vault.
package grpc

import (
	"context"
	"fmt"
	"net"

	"credentials-vault/internal/config"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	*grpc.Server
	addr   string
	logger *zap.Logger
}

// NewServer создаёт gRPC-сервер.
func NewServer(cfg *config.Config, logger *zap.Logger) *Server {
	s := grpc.NewServer()

	if !cfg.IsProd() {
		reflection.Register(s)
		logger.Info("gRPC reflection enabled (dev mode)")
	}

	return &Server{
		Server: s,
		addr:   cfg.GRPCAddress,
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
