package get

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"credentials-vault/client-cli/internal/command/get/mocks"
	"credentials-vault/client-cli/internal/model"
)

func TestGetCmd(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setupMock func(*mocks.MockClient)
		wantErr   bool
	}{
		{
			name: "success",
			args: []string{"test-id"},
			setupMock: func(mockClient *mocks.MockClient) {
				item := &model.VaultItem{
					ID:   "test-id",
					Type: "ITEM_TYPE_LOGIN",
					Secret: map[string]string{
						"username": "bob",
						"password": "secret",
					},
					Metadata: map[string]string{
						"site": "example.com",
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				mockClient.EXPECT().
					GetItem(gomock.Any(), "test-id").
					Return(item, nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			args: []string{"test-id"},
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					GetItem(gomock.Any(), "test-id").
					Return(nil, errors.New("get failed"))
			},
			wantErr: true,
		},
		{
			name:      "no args",
			args:      []string{},
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

			cmd := NewGetCmd(mockClient)
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
