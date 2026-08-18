package command

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	vaultpb "credentials-vault/gen/go/vault/v1"

	"credentials-vault/client-cli/internal/command/mocks"
)

func TestListCmd(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockClient)
		wantErr   bool
	}{
		{
			name: "success with items",
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					ListItems(gomock.Any()).
					Return([]*vaultpb.VaultItem{
						{
							Id:        uuid.New().String(),
							Type:      vaultpb.ItemType_ITEM_TYPE_LOGIN,
							CreatedAt: timestamppb.Now(),
						},
						{
							Id:        uuid.New().String(),
							Type:      vaultpb.ItemType_ITEM_TYPE_CARD,
							CreatedAt: timestamppb.Now(),
						},
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "success empty",
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					ListItems(gomock.Any()).
					Return([]*vaultpb.VaultItem{}, nil)
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

			cmd := newListCmd(mockClient)
			err := cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
