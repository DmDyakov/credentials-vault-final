package interceptor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authpb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/server/internal/transport/grpc/interceptor/mocks"
	"credentials-vault/pkg/jwt"
)

// mockServerTransportStream - мок для grpc.ServerTransportStream
type mockServerTransportStream struct {
	grpc.ServerTransportStream
	headers metadata.MD
}

func (m *mockServerTransportStream) SendHeader(md metadata.MD) error {
	m.headers = md
	return nil
}

func (m *mockServerTransportStream) Method() string {
	return "test-method"
}

func setupTest(t *testing.T) (*AuthInterceptor, *mocks.MockJWTManager) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockJWT := mocks.NewMockJWTManager(ctrl)
	interceptor := NewAuthInterceptor(mockJWT)

	return interceptor, mockJWT
}

func TestUnary_LoginMethod(t *testing.T) {
	interceptor, mockJWT := setupTest(t)

	userID := "user-123"
	token := jwt.Token("test-token")
	expiresAt := time.Now().Add(24 * time.Hour)

	loginResp := &authpb.LoginResponse{
		User: &authpb.User{
			Id:       userID,
			Username: "testuser",
		},
	}

	mockJWT.EXPECT().
		Generate(userID).
		Return(token, expiresAt, nil)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return loginResp, nil
	}

	stream := &mockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), stream)
	ctx = metadata.NewIncomingContext(ctx, metadata.New(nil))

	resp, err := interceptor.Unary()(
		ctx,
		&authpb.LoginRequest{Username: "testuser", Password: "password123"},
		&grpc.UnaryServerInfo{FullMethod: authpb.AuthService_Login_FullMethodName},
		handler,
	)

	assert.NoError(t, err)
	assert.Equal(t, loginResp, resp)
	assert.NotNil(t, stream.headers)
	assert.Equal(t, "Bearer test-token", stream.headers.Get("authorization")[0])
}

func TestUnary_LoginMethod_GenerateTokenError(t *testing.T) {
	interceptor, mockJWT := setupTest(t)

	userID := "user-123"

	loginResp := &authpb.LoginResponse{
		User: &authpb.User{
			Id:       userID,
			Username: "testuser",
		},
	}

	mockJWT.EXPECT().
		Generate(userID).
		Return(jwt.Token(""), time.Time{}, errors.New("failed to generate token"))

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return loginResp, nil
	}

	stream := &mockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), stream)
	ctx = metadata.NewIncomingContext(ctx, metadata.New(nil))

	resp, err := interceptor.Unary()(
		ctx,
		&authpb.LoginRequest{Username: "testuser", Password: "password123"},
		&grpc.UnaryServerInfo{FullMethod: authpb.AuthService_Login_FullMethodName},
		handler,
	)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestUnary_RegisterMethod(t *testing.T) {
	interceptor, _ := setupTest(t)

	registerResp := &authpb.RegisterResponse{
		User: &authpb.User{
			Id:       "user-123",
			Username: "testuser",
		},
		Message: "User registered successfully",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return registerResp, nil
	}

	resp, err := interceptor.Unary()(
		context.Background(),
		&authpb.RegisterRequest{Username: "testuser", Password: "password123"},
		&grpc.UnaryServerInfo{FullMethod: authpb.AuthService_Register_FullMethodName},
		handler,
	)

	assert.NoError(t, err)
	assert.Equal(t, registerResp, resp)
}

func TestUnary_AuthenticatedMethod_ValidToken(t *testing.T) {
	interceptor, mockJWT := setupTest(t)

	token := jwt.Token("valid-token")
	claims := &jwt.Claims{UserID: "user-123"}

	mockJWT.EXPECT().
		Verify(token).
		Return(claims, nil)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, "user-123", md.Get("user_id")[0])
		return "response", nil
	}

	md := metadata.Pairs("authorization", "Bearer valid-token")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor.Unary()(
		ctx,
		"request",
		&grpc.UnaryServerInfo{FullMethod: "/credentials_vault.vault.v1.VaultService/ListItems"},
		handler,
	)

	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
}

func TestUnary_AuthenticatedMethod_NoMetadata(t *testing.T) {
	interceptor, _ := setupTest(t)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	resp, err := interceptor.Unary()(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{FullMethod: "/credentials_vault.vault.v1.VaultService/ListItems"},
		handler,
	)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, "metadata is not provided", st.Message())
}

func TestUnary_AuthenticatedMethod_NoToken(t *testing.T) {
	interceptor, _ := setupTest(t)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	md := metadata.New(nil)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor.Unary()(
		ctx,
		"request",
		&grpc.UnaryServerInfo{FullMethod: "/credentials_vault.vault.v1.VaultService/ListItems"},
		handler,
	)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, "authorization token is not provided", st.Message())
}

func TestUnary_AuthenticatedMethod_InvalidFormat(t *testing.T) {
	interceptor, _ := setupTest(t)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	md := metadata.Pairs("authorization", "invalid-format-token")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor.Unary()(
		ctx,
		"request",
		&grpc.UnaryServerInfo{FullMethod: "/credentials_vault.vault.v1.VaultService/ListItems"},
		handler,
	)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, "invalid authorization format", st.Message())
}

func TestUnary_AuthenticatedMethod_InvalidToken(t *testing.T) {
	interceptor, mockJWT := setupTest(t)

	token := jwt.Token("invalid-token")

	mockJWT.EXPECT().
		Verify(token).
		Return(nil, errors.New("invalid token"))

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	md := metadata.Pairs("authorization", "Bearer invalid-token")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor.Unary()(
		ctx,
		"request",
		&grpc.UnaryServerInfo{FullMethod: "/credentials_vault.vault.v1.VaultService/ListItems"},
		handler,
	)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, "invalid token", st.Message())
}

func TestUnary_AuthenticatedMethod_ExpiredToken(t *testing.T) {
	interceptor, mockJWT := setupTest(t)

	token := jwt.Token("expired-token")

	mockJWT.EXPECT().
		Verify(token).
		Return(nil, errors.New("token expired"))

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "should not be called", nil
	}

	md := metadata.Pairs("authorization", "Bearer expired-token")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor.Unary()(
		ctx,
		"request",
		&grpc.UnaryServerInfo{FullMethod: "/credentials_vault.vault.v1.VaultService/ListItems"},
		handler,
	)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, "invalid token", st.Message())
}
