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

func TestRegister(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		username  string
		password  string
		setupMock func(*mocks.MockUserRepository)
		wantErr   error
		wantUser  bool
	}{
		{
			name:     "success",
			username: "testuser",
			password: "password123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					FindByUsername(gomock.Any(), "testuser").
					Return(nil, domain.ErrUserNotFound)

				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, user *domain.User) error {
						user.ID = userID
						return nil
					})
			},
			wantErr:  nil,
			wantUser: true,
		},
		{
			name:     "duplicate user",
			username: "testuser",
			password: "password123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					FindByUsername(gomock.Any(), "testuser").
					Return(&domain.User{ID: userID, Username: "testuser"}, nil)
			},
			wantErr:  domain.ErrUserAlreadyExists,
			wantUser: false,
		},
		{
			name:     "empty username",
			username: "",
			password: "password123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
			},
			wantErr:  domain.ErrUsernameRequired,
			wantUser: false,
		},
		{
			name:     "short username",
			username: "ab",
			password: "password123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
			},
			wantErr:  domain.ErrUsernameTooShort,
			wantUser: false,
		},
		{
			name:     "long username",
			username: strings.Repeat("a", 65),
			password: "password123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
			},
			wantErr:  domain.ErrUsernameTooLong,
			wantUser: false,
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
			},
			wantErr:  domain.ErrPasswordRequired,
			wantUser: false,
		},
		{
			name:     "short password",
			username: "testuser",
			password: "123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
			},
			wantErr:  domain.ErrPasswordTooShort,
			wantUser: false,
		},
		{
			name:     "long password",
			username: "testuser",
			password: strings.Repeat("p", 73),
			setupMock: func(mockRepo *mocks.MockUserRepository) {
			},
			wantErr:  domain.ErrPasswordTooLong,
			wantUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockRepo := mocks.NewMockUserRepository(ctrl)
			service := NewService(mockRepo)

			tt.setupMock(mockRepo)

			user, err := service.Register(context.Background(), tt.username, tt.password)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				if tt.wantUser {
					assert.NotNil(t, user)
					assert.NotEmpty(t, user.ID)
					assert.Equal(t, tt.username, user.Username)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	userID := uuid.New()

	tests := []struct {
		name      string
		username  string
		password  string
		setupMock func(*mocks.MockUserRepository)
		wantErr   error
		wantUser  bool
	}{
		{
			name:     "success",
			username: "testuser",
			password: "password123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					FindByUsername(gomock.Any(), "testuser").
					Return(&domain.User{
						ID:       userID,
						Username: "testuser",
						Password: string(hashedPassword),
					}, nil)
			},
			wantErr:  nil,
			wantUser: true,
		},
		{
			name:     "wrong password",
			username: "testuser",
			password: "wrong-password",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					FindByUsername(gomock.Any(), "testuser").
					Return(&domain.User{
						ID:       userID,
						Username: "testuser",
						Password: string(hashedPassword),
					}, nil)
			},
			wantErr:  domain.ErrInvalidCredentials,
			wantUser: false,
		},
		{
			name:     "user not found",
			username: "unknown",
			password: "password123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					FindByUsername(gomock.Any(), "unknown").
					Return(nil, domain.ErrUserNotFound)
			},
			wantErr:  domain.ErrInvalidCredentials,
			wantUser: false,
		},
		{
			name:     "empty username",
			username: "",
			password: "password123",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
			},
			wantErr:  domain.ErrUsernameRequired,
			wantUser: false,
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			setupMock: func(mockRepo *mocks.MockUserRepository) {
			},
			wantErr:  domain.ErrPasswordRequired,
			wantUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockRepo := mocks.NewMockUserRepository(ctrl)
			service := NewService(mockRepo)

			tt.setupMock(mockRepo)

			user, err := service.Login(context.Background(), tt.username, tt.password)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				if tt.wantUser {
					assert.NotNil(t, user)
					assert.Equal(t, tt.username, user.Username)
				}
			}
		})
	}
}
