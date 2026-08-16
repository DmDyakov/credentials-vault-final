package auth

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/internal/model"
	"credentials-vault/internal/repository"
	"credentials-vault/internal/service/auth/mocks"
	"credentials-vault/pkg/jwt"
)

func setupTest(t *testing.T) (*AuthService, *mocks.MockUserRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := jwt.New("test-secret-123456789012345678901234567890", 24*time.Hour)
	service := NewAuthService(mockRepo, jwtManager)

	return service, mockRepo
}

func TestRegister_Success(t *testing.T) {
	service, mockRepo := setupTest(t)

	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "testuser").
		Return(nil, repository.ErrUserNotFound)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *model.User) error {
			user.ID = uuid.New()
			return nil
		})

	resp, err := service.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp.User)
	assert.NotEmpty(t, resp.User.Id)
	assert.Equal(t, "testuser", resp.User.Username)
	assert.Equal(t, "User registered successfully", resp.Message)
}

func TestRegister_DuplicateUser(t *testing.T) {
	service, mockRepo := setupTest(t)

	existingUser := &model.User{ID: uuid.New(), Username: "testuser"}
	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "testuser").
		Return(existingUser, nil)

	_, err := service.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
}

func TestRegister_Validation(t *testing.T) {
	service, _ := setupTest(t)

	tests := []struct {
		name     string
		username string
		password string
		wantCode codes.Code
	}{
		{
			name:     "empty username",
			username: "",
			password: "password123",
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "short username",
			username: "ab",
			password: "password123",
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "long username",
			username: strings.Repeat("a", 65),
			password: "password123",
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "short password",
			username: "testuser",
			password: "123",
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "long password",
			username: "testuser",
			password: strings.Repeat("p", 73),
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Register(context.Background(), &pb.RegisterRequest{
				Username: tt.username,
				Password: tt.password,
			})

			assert.Error(t, err)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.wantCode, st.Code())
		})
	}
}

func TestLogin_Success(t *testing.T) {
	service, mockRepo := setupTest(t)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := &model.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "testuser").
		Return(user, nil)

	resp, err := service.Login(context.Background(), &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, "testuser", resp.User.Username)
}

func TestLogin_WrongPassword(t *testing.T) {
	service, mockRepo := setupTest(t)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := &model.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "testuser").
		Return(user, nil)

	_, err := service.Login(context.Background(), &pb.LoginRequest{
		Username: "testuser",
		Password: "wrong-password",
	})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestLogin_UserNotFound(t *testing.T) {
	service, mockRepo := setupTest(t)

	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "unknown").
		Return(nil, repository.ErrUserNotFound)

	_, err := service.Login(context.Background(), &pb.LoginRequest{
		Username: "unknown",
		Password: "password123",
	})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestLogin_Validation(t *testing.T) {
	service, _ := setupTest(t)

	tests := []struct {
		name     string
		username string
		password string
		wantCode codes.Code
	}{
		{
			name:     "empty username",
			username: "",
			password: "password123",
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(context.Background(), &pb.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			})

			assert.Error(t, err)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.wantCode, st.Code())
		})
	}
}
