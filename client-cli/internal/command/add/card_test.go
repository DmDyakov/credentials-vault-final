package add

import (
	"errors"
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
			name: "success",
			args: []string{
				"--brand", "visa",
				"--bank", "sberbank",
				"--number", "4111111111111111",
				"--holder", "IVAN IVANOV",
				"--expiry", "12/25",
				"--cvv", "123",
			},
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					AddCard(gomock.Any(), "visa", "sberbank", "4111111111111111", "IVAN IVANOV", "12/25", "123").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			args: []string{
				"--brand", "visa",
				"--number", "4111111111111111",
				"--holder", "IVAN IVANOV",
				"--expiry", "12/25",
				"--cvv", "123",
			},
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					AddCard(gomock.Any(), "visa", "", "4111111111111111", "IVAN IVANOV", "12/25", "123").
					Return(errors.New("add failed"))
			},
			wantErr: true,
		},
		{
			name: "missing brand",
			args: []string{
				"--number", "4111111111111111",
				"--holder", "IVAN IVANOV",
				"--expiry", "12/25",
				"--cvv", "123",
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
				"--cvv", "123",
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
				"--cvv", "123",
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
				"--cvv", "123",
			},
			setupMock: func(mockClient *mocks.MockClient) {},
			wantErr:   true,
		},
		{
			name: "missing cvv",
			args: []string{
				"--brand", "visa",
				"--number", "4111111111111111",
				"--holder", "IVAN IVANOV",
				"--expiry", "12/25",
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
