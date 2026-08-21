package command

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	vaultpb "credentials-vault/gen/go/vault/v1"

	"credentials-vault/client-cli/internal/command/mocks"
)

func TestGetCmd(t *testing.T) {
	itemID := uuid.New().String()

	tests := []struct {
		name      string
		args      []string
		setupMock func(*mocks.MockClient)
		wantErr   bool
	}{
		{
			name: "success",
			args: []string{itemID},
			setupMock: func(mockClient *mocks.MockClient) {
				item := vaultpb.VaultItem_builder{
					Id:        proto.String(itemID),
					Type:      vaultpb.ItemType_ITEM_TYPE_LOGIN.Enum(),
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
				}.Build()

				mockClient.EXPECT().
					GetItem(gomock.Any(), itemID).
					Return(item, nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			args: []string{itemID},
			setupMock: func(mockClient *mocks.MockClient) {
				mockClient.EXPECT().
					GetItem(gomock.Any(), itemID).
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

			cmd := newGetCmd(mockClient)
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
