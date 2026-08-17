// Package handler содержит gRPC обработчики.
package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	authpb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/internal/domain"
)

// UserService - интерфейс сервиса пользователей
type UserService interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.User, error)
}

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	userService UserService
}

func NewAuthHandler(userService UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	user, err := h.userService.Register(ctx, req.Username, req.Password)
	if err != nil {
		return nil, mapAuthError(err)
	}

	return &authpb.RegisterResponse{
		User:    domainUserToProto(user),
		Message: "User registered successfully",
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	user, err := h.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, mapAuthError(err)
	}

	return &authpb.LoginResponse{
		User: domainUserToProto(user),
	}, nil
}

// mapAuthError маппит доменные ошибки в gRPC статусы
func mapAuthError(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "username already exists")
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid username or password")
	case errors.Is(err, domain.ErrUsernameRequired),
		errors.Is(err, domain.ErrPasswordRequired),
		errors.Is(err, domain.ErrUsernameTooShort),
		errors.Is(err, domain.ErrUsernameTooLong),
		errors.Is(err, domain.ErrPasswordTooShort),
		errors.Is(err, domain.ErrPasswordTooLong):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}

// domainUserToProto конвертирует доменную модель в proto
func domainUserToProto(user *domain.User) *authpb.User {
	return &authpb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
