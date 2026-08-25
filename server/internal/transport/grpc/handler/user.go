// Package handler содержит gRPC обработчики.
package handler

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	authpb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/server/internal/domain"
)

//go:generate mockgen -source=user.go -destination=mocks/user_service_mock.gen.go -package=mocks UserService
type UserService interface {
	Register(ctx context.Context, username, password string, encryptionSalt []byte) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.User, error)
}

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	userService UserService
	logger      *zap.Logger
}

func NewAuthHandler(userService UserService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	user, err := h.userService.Register(ctx, req.GetUsername(), req.GetPassword(), req.GetSalt())
	if err != nil {
		return nil, mapError(err, h.logger)
	}

	return authpb.RegisterResponse_builder{
		User:    toProtoUser(user),
		Message: proto.String("User registered successfully"),
	}.Build(), nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	user, err := h.userService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, mapError(err, h.logger)
	}

	return authpb.LoginResponse_builder{
		User: toProtoUser(user),
		Salt: user.Salt,
	}.Build(), nil
}
