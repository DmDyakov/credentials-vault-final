package user

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"credentials-vault/internal/domain"
	"credentials-vault/internal/service/user/mocks"
)

func setupTest(t *testing.T) (*UserService, *mocks.MockUserRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewService(mockRepo)

	return service, mockRepo
}

func TestRegister_Success(t *testing.T) {
	service, mockRepo := setupTest(t)

	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "testuser").
		Return(nil, domain.ErrUserNotFound)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *domain.User) error {
			user.ID = uuid.New()
			return nil
		})

	user, err := service.Register(context.Background(), "testuser", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "testuser", user.Username)
}

func TestRegister_DuplicateUser(t *testing.T) {
	service, mockRepo := setupTest(t)

	existingUser := &domain.User{ID: uuid.New(), Username: "testuser"}
	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "testuser").
		Return(existingUser, nil)

	_, err := service.Register(context.Background(), "testuser", "password123")

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
}

func TestRegister_Validation(t *testing.T) {
	service, _ := setupTest(t)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{
			name:     "empty username",
			username: "",
			password: "password123",
			wantErr:  domain.ErrUsernameRequired,
		},
		{
			name:     "short username",
			username: "ab",
			password: "password123",
			wantErr:  domain.ErrUsernameTooShort,
		},
		{
			name:     "long username",
			username: strings.Repeat("a", 65),
			password: "password123",
			wantErr:  domain.ErrUsernameTooLong,
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			wantErr:  domain.ErrPasswordRequired,
		},
		{
			name:     "short password",
			username: "testuser",
			password: "123",
			wantErr:  domain.ErrPasswordTooShort,
		},
		{
			name:     "long password",
			username: "testuser",
			password: strings.Repeat("p", 73),
			wantErr:  domain.ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Register(context.Background(), tt.username, tt.password)

			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLogin_Success(t *testing.T) {
	service, mockRepo := setupTest(t)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := &domain.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "testuser").
		Return(user, nil)

	loginUser, err := service.Login(context.Background(), "testuser", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, loginUser)
	assert.Equal(t, "testuser", loginUser.Username)
}

func TestLogin_WrongPassword(t *testing.T) {
	service, mockRepo := setupTest(t)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := &domain.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "testuser").
		Return(user, nil)

	_, err := service.Login(context.Background(), "testuser", "wrong-password")

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_UserNotFound(t *testing.T) {
	service, mockRepo := setupTest(t)

	mockRepo.EXPECT().
		FindByUsername(gomock.Any(), "unknown").
		Return(nil, domain.ErrUserNotFound)

	_, err := service.Login(context.Background(), "unknown", "password123")

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_Validation(t *testing.T) {
	service, _ := setupTest(t)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{
			name:     "empty username",
			username: "",
			password: "password123",
			wantErr:  domain.ErrUsernameRequired,
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			wantErr:  domain.ErrPasswordRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(context.Background(), tt.username, tt.password)

			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
