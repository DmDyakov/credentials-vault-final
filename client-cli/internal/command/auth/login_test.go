package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"credentials-vault/client-cli/internal/command/auth/mocks"
)

func TestLoginCmd(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setupMock func(*mocks.MockClient)
		wantErr   bool
	}{
		{
			name: "success",
			args: []string{"--username", "testuser", "--password", "password123"},
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					Login(gomock.Any(), "testuser", "password123").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			args: []string{"--username", "testuser", "--password", "password123"},
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					Login(gomock.Any(), "testuser", "password123").
					Return(errors.New("login failed"))
			},
			wantErr: true,
		},
		{
			name:      "missing username",
			args:      []string{"--password", "password123"},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name:      "missing password",
			args:      []string{"--username", "testuser"},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockClient := mocks.NewMockClient(ctrl)
			tt.setupMock(mockClient)

			cmd := NewLoginCmd(mockClient)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
