// Package interceptor содержит gRPC интерсепторы.
package interceptor

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authpb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/pkg/jwt"
)

//go:generate mockgen -source=auth.go -destination=mocks/jwt_manager_mock.go -package=mocks JWTManager
type JWTManager interface {
	Generate(userID string) (jwt.Token, time.Time, error)
	Verify(token jwt.Token) (*jwt.Claims, error)
}

type AuthInterceptor struct {
	jwtManager JWTManager
}

func NewAuthInterceptor(jwtManager JWTManager) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager: jwtManager,
	}
}

// Unary возвращает unary интерсептор
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		switch info.FullMethod {
		case authpb.AuthService_Login_FullMethodName:
			return i.handleLogin(ctx, req, handler)
		case authpb.AuthService_Register_FullMethodName:
			return i.handleRegister(ctx, req, handler)
		default:
			return i.handleAuthenticated(ctx, req, handler)
		}
	}
}

// handleLogin обрабатывает логин и добавляет токен в метаданные ответа
func (i *AuthInterceptor) handleLogin(
	ctx context.Context,
	req interface{},
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	loginResp, ok := resp.(*authpb.LoginResponse)
	if !ok || loginResp.User == nil {
		return resp, nil
	}

	token, _, err := i.jwtManager.Generate(loginResp.User.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	header := metadata.Pairs("authorization", "Bearer "+token.String())
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, status.Error(codes.Internal, "failed to send authorization header")
	}

	return loginResp, nil
}

// handleRegister обрабатывает регистрацию без аутентификации
func (i *AuthInterceptor) handleRegister(
	ctx context.Context,
	req interface{},
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return handler(ctx, req)
}

// handleAuthenticated проверяет токен и добавляет user_id в метаданные
func (i *AuthInterceptor) handleAuthenticated(
	ctx context.Context,
	req interface{},
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	tokenString := values[0]
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
	}
	token := jwt.Token(strings.TrimPrefix(tokenString, "Bearer "))

	claims, err := i.jwtManager.Verify(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	md.Set("user_id", claims.UserID)
	newCtx := metadata.NewIncomingContext(ctx, md)

	return handler(newCtx, req)
}
