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

func TestRegister(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		req       *authpb.RegisterRequest
		setupMock func(*mocks.MockUserService)
		wantCode  codes.Code
		wantMsg   string
		wantUser  bool
	}{
		{
			name: "success",
			req: &authpb.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func(mockService *mocks.MockUserService) {
				mockService.EXPECT().
					Register(gomock.Any(), "testuser", "password123").
					Return(&domain.User{ID: userID, Username: "testuser"}, nil)
			},
			wantCode: codes.OK,
			wantUser: true,
		},
		{
			name: "nil request",
			req:  nil,
			setupMock: func(mockService *mocks.MockUserService) {
			},
			wantCode: codes.InvalidArgument,
			wantUser: false,
		},
		{
			name: "duplicate user",
			req: &authpb.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func(mockService *mocks.MockUserService) {
				mockService.EXPECT().
					Register(gomock.Any(), "testuser", "password123").
					Return(nil, domain.ErrUserAlreadyExists)
			},
			wantCode: codes.AlreadyExists,
			wantMsg:  "user already exists",
			wantUser: false,
		},
		{
			name: "validation error",
			req: &authpb.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func(mockService *mocks.MockUserService) {
				mockService.EXPECT().
					Register(gomock.Any(), "testuser", "password123").
					Return(nil, domain.ErrUsernameRequired)
			},
			wantCode: codes.InvalidArgument,
			wantMsg:  "username is required",
			wantUser: false,
		},
		{
			name: "internal error",
			req: &authpb.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func(mockService *mocks.MockUserService) {
				mockService.EXPECT().
					Register(gomock.Any(), "testuser", "password123").
					Return(nil, errors.New("database error"))
			},
			wantCode: codes.Internal,
			wantMsg:  "internal error",
			wantUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockUserService(ctrl)
			handler := NewAuthHandler(mockService)

			tt.setupMock(mockService)

			resp, err := handler.Register(context.Background(), tt.req)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.wantUser {
					assert.NotNil(t, resp.User)
					assert.Equal(t, userID.String(), resp.User.Id)
					assert.Equal(t, "testuser", resp.User.Username)
				}
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
				if tt.wantMsg != "" {
					assert.Equal(t, tt.wantMsg, st.Message())
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		req       *authpb.LoginRequest
		setupMock func(*mocks.MockUserService)
		wantCode  codes.Code
		wantMsg   string
		wantUser  bool
	}{
		{
			name: "success",
			req: &authpb.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func(mockService *mocks.MockUserService) {
				mockService.EXPECT().
					Login(gomock.Any(), "testuser", "password123").
					Return(&domain.User{ID: userID, Username: "testuser"}, nil)
			},
			wantCode: codes.OK,
			wantUser: true,
		},
		{
			name: "nil request",
			req:  nil,
			setupMock: func(mockService *mocks.MockUserService) {
			},
			wantCode: codes.InvalidArgument,
			wantUser: false,
		},
		{
			name: "invalid credentials",
			req: &authpb.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setupMock: func(mockService *mocks.MockUserService) {
				mockService.EXPECT().
					Login(gomock.Any(), "testuser", "wrongpassword").
					Return(nil, domain.ErrInvalidCredentials)
			},
			wantCode: codes.Unauthenticated,
			wantMsg:  "invalid username or password",
			wantUser: false,
		},
		{
			name: "internal error",
			req: &authpb.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func(mockService *mocks.MockUserService) {
				mockService.EXPECT().
					Login(gomock.Any(), "testuser", "password123").
					Return(nil, errors.New("database error"))
			},
			wantCode: codes.Internal,
			wantMsg:  "internal error",
			wantUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockUserService(ctrl)
			handler := NewAuthHandler(mockService)

			tt.setupMock(mockService)

			resp, err := handler.Login(context.Background(), tt.req)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.wantUser {
					assert.NotNil(t, resp.User)
					assert.Equal(t, userID.String(), resp.User.Id)
					assert.Equal(t, "testuser", resp.User.Username)
				}
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
				if tt.wantMsg != "" {
					assert.Equal(t, tt.wantMsg, st.Message())
				}
			}
		})
	}
}

func TestToProtoUser(t *testing.T) {
	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Username: "testuser",
	}

	protoUser := toProtoUser(user)

	assert.NotNil(t, protoUser)
	assert.Equal(t, userID.String(), protoUser.Id)
	assert.Equal(t, "testuser", protoUser.Username)
}
