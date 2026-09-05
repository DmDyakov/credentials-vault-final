package list

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"credentials-vault/client-cli/internal/command/list/mocks"
	"credentials-vault/client-cli/internal/model"
)

func newTestListVaultItem() *model.ListVaultItem {
	return &model.ListVaultItem{
		ID:        uuid.New().String(),
		Type:      "ITEM_TYPE_LOGIN",
		Metadata:  map[string]string{"site": "example.com"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestListCmd(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockClient)
		wantErr   bool
	}{
		{
			name: "success with items",
			setupMock: func(mockClient *mocks.MockClient) {
				items := []*model.ListVaultItem{
					newTestListVaultItem(),
					newTestListVaultItem(),
				}

				mockClient.EXPECT().
					ListItems(gomock.Any()).
					Return(items, nil)
			},
			wantErr: false,
		},
		{
			name: "success empty",
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					ListItems(gomock.Any()).
					Return([]*model.ListVaultItem{}, nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					ListItems(gomock.Any()).
					Return(nil, errors.New("list failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockClient := mocks.NewMockClient(ctrl)
			tt.setupMock(mockClient)

			cmd := NewListCmd(mockClient)
			err := cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
