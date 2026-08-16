package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/internal/model"
	"credentials-vault/internal/repository"
	"credentials-vault/internal/service/mocks"
)

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo)

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

func TestRegister_EmptyPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo)

	_, err := service.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Password: "",
	})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestRegister_ShortUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo)

	_, err := service.Register(context.Background(), &pb.RegisterRequest{
		Username: "ab",
		Password: "password123",
	})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestRegister_ShortPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo)

	_, err := service.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Password: "123",
	})

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}
