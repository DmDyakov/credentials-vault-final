// Package service содержит бизнес-логику.
package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/internal/model"
	"credentials-vault/internal/repository"
	"credentials-vault/pkg/jwt"
)

//go:generate mockgen -source=auth.go -destination=mocks/users_mock.go -package=mocks UserRepository
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	repo       UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(repo UserRepository, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := validateRegisterCredentials(req.Username, req.Password); err != nil {
		return nil, err
	}

	_, err := s.repo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, status.Errorf(codes.Internal, "failed to check user existence: %v", err)
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

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := validateLoginCredentials(req.Username, req.Password); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invalid username or password")
		}
		return nil, status.Error(codes.Internal, "failed to find user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	token, expiresAt, err := s.jwtManager.Generate(user.ID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{
		AccessToken: token,
		ExpiresAt:   timestamppb.New(expiresAt),
		User:        modelToProto(user),
	}, nil
}
