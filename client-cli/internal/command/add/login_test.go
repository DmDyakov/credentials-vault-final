package add

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"credentials-vault/client-cli/internal/command/add/mocks"
)

func TestAddLoginCmd(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setupMock func(*mocks.MockClient)
		wantErr   bool
	}{
		{
			name: "success",
			args: []string{
				"--site", "example.com",
				"--username", "user",
				"--password", "pass",
			},
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					AddLogin(gomock.Any(), "example.com", "user", "pass").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			args: []string{
				"--site", "example.com",
				"--username", "user",
				"--password", "pass",
			},
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					AddLogin(gomock.Any(), "example.com", "user", "pass").
					Return(errors.New("add failed"))
			},
			wantErr: true,
		},
		{
			name: "missing site",
			args: []string{
				"--username", "user",
				"--password", "pass",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "missing username",
			args: []string{
				"--site", "example.com",
				"--password", "pass",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "missing password",
			args: []string{
				"--site", "example.com",
				"--username", "user",
			},
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

			cmd := newLoginCmd(mockClient)
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
