package add

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"credentials-vault/client-cli/internal/command/add/mocks"
)

func TestAddCredentialsCmd(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setupMock func(*mocks.MockClient)
		wantErr   bool
	}{
		{
			name: "missing site",
			args: []string{
				"--username", "user",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "missing username",
			args: []string{
				"--site", "example.com",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "empty site",
			args: []string{
				"--site", "",
				"--username", "user",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "empty username",
			args: []string{
				"--site", "example.com",
				"--username", "",
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

			cmd := newCredentialsCmd(mockClient)
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
