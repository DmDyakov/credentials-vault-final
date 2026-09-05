package add

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"credentials-vault/client-cli/internal/command/add/mocks"
)

func TestAddCardCmd(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setupMock func(*mocks.MockClient)
		wantErr   bool
	}{
		{
			name: "missing brand",
			args: []string{
				"--number", "4111111111111111",
				"--holder", "IVAN IVANOV",
				"--expiry", "12/25",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "missing number",
			args: []string{
				"--brand", "visa",
				"--holder", "IVAN IVANOV",
				"--expiry", "12/25",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "missing holder",
			args: []string{
				"--brand", "visa",
				"--number", "4111111111111111",
				"--expiry", "12/25",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "missing expiry",
			args: []string{
				"--brand", "visa",
				"--number", "4111111111111111",
				"--holder", "IVAN IVANOV",
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

			cmd := newCardCmd(mockClient)
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
