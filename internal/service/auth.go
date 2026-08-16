// Package service содержит бизнес-логику.
package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/internal/model"
	"credentials-vault/internal/repository"
)

//go:generate mockgen -source=auth.go -destination=mocks/users_mock.go -package=mocks UserRepository
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := validateCredentials(req.Username, req.Password); err != nil {
		return nil, err
	}

	_, err := s.repo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, status.Error(codes.Internal, "failed to check user existence")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "username already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.RegisterResponse{
		User:    modelToProto(user),
		Message: "User registered successfully",
	}, nil
}

func modelToProto(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func validateCredentials(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}
	if len(username) < 3 {
		return status.Error(codes.InvalidArgument, "username must be at least 3 characters")
	}
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if len(password) < 6 {
		return status.Error(codes.InvalidArgument, "password must be at least 6 characters")
	}
	return nil
}
