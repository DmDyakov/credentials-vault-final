package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/internal/domain"
	"credentials-vault/internal/transport/grpc/handler/mocks"
)

func setupTest(t *testing.T) (*AuthHandler, *mocks.MockUserService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockService := mocks.NewMockUserService(ctrl)
	handler := NewAuthHandler(mockService)

	return handler, mockService
}

func TestRegister_Success(t *testing.T) {
	handler, mockService := setupTest(t)

	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Username: "testuser",
	}

	mockService.EXPECT().
		Register(gomock.Any(), "testuser", "password123").
		Return(user, nil)

	resp, err := handler.Register(context.Background(), &authpb.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, userID.String(), resp.User.Id)
	assert.Equal(t, "testuser", resp.User.Username)
	assert.Equal(t, "User registered successfully", resp.Message)
}

func TestRegister_NilRequest(t *testing.T) {
	handler, _ := setupTest(t)

	resp, err := handler.Register(context.Background(), nil)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestRegister_DuplicateUser(t *testing.T) {
	handler, mockService := setupTest(t)

	mockService.EXPECT().
		Register(gomock.Any(), "testuser", "password123").
		Return(nil, domain.ErrUserAlreadyExists)

	resp, err := handler.Register(context.Background(), &authpb.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Equal(t, "user already exists", st.Message())
}

func TestRegister_ValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode codes.Code
		wantMsg  string
	}{
		{
			name:     "username required",
			err:      domain.ErrUsernameRequired,
			wantCode: codes.InvalidArgument,
			wantMsg:  "username is required",
		},
		{
			name:     "username too short",
			err:      domain.ErrUsernameTooShort,
			wantCode: codes.InvalidArgument,
			wantMsg:  "username must be at least 3 characters",
		},
		{
			name:     "username too long",
			err:      domain.ErrUsernameTooLong,
			wantCode: codes.InvalidArgument,
			wantMsg:  "username must be at most 64 characters",
		},
		{
			name:     "password required",
			err:      domain.ErrPasswordRequired,
			wantCode: codes.InvalidArgument,
			wantMsg:  "password is required",
		},
		{
			name:     "password too short",
			err:      domain.ErrPasswordTooShort,
			wantCode: codes.InvalidArgument,
			wantMsg:  "password must be at least 6 characters",
		},
		{
			name:     "password too long",
			err:      domain.ErrPasswordTooLong,
			wantCode: codes.InvalidArgument,
			wantMsg:  "password must be at most 72 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupTest(t)

			mockService.EXPECT().
				Register(gomock.Any(), "testuser", "password123").
				Return(nil, tt.err)

			resp, err := handler.Register(context.Background(), &authpb.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			})

			assert.Nil(t, resp)
			assert.Error(t, err)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.wantCode, st.Code())
			assert.Equal(t, tt.wantMsg, st.Message())
		})
	}
}

func TestRegister_InternalError(t *testing.T) {
	handler, mockService := setupTest(t)

	mockService.EXPECT().
		Register(gomock.Any(), "testuser", "password123").
		Return(nil, errors.New("database connection failed"))

	resp, err := handler.Register(context.Background(), &authpb.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "internal error", st.Message())
}

func TestLogin_Success(t *testing.T) {
	handler, mockService := setupTest(t)

	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Username: "testuser",
	}

	mockService.EXPECT().
		Login(gomock.Any(), "testuser", "password123").
		Return(user, nil)

	resp, err := handler.Login(context.Background(), &authpb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, userID.String(), resp.User.Id)
	assert.Equal(t, "testuser", resp.User.Username)
}

func TestLogin_NilRequest(t *testing.T) {
	handler, _ := setupTest(t)

	resp, err := handler.Login(context.Background(), nil)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestLogin_InvalidCredentials(t *testing.T) {
	handler, mockService := setupTest(t)

	mockService.EXPECT().
		Login(gomock.Any(), "testuser", "wrongpassword").
		Return(nil, domain.ErrInvalidCredentials)

	resp, err := handler.Login(context.Background(), &authpb.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Equal(t, "invalid username or password", st.Message())
}

func TestLogin_InternalError(t *testing.T) {
	handler, mockService := setupTest(t)

	mockService.EXPECT().
		Login(gomock.Any(), "testuser", "password123").
		Return(nil, errors.New("database connection failed"))

	resp, err := handler.Login(context.Background(), &authpb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestDomainUserToProto(t *testing.T) {
	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Username: "testuser",
	}

	protoUser := domainUserToProto(user)

	assert.NotNil(t, protoUser)
	assert.Equal(t, userID.String(), protoUser.Id)
	assert.Equal(t, "testuser", protoUser.Username)
}
