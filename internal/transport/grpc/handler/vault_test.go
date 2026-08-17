package handler

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	vaultpb "credentials-vault/gen/go/vault/v1"
	"credentials-vault/internal/domain"
	"credentials-vault/internal/transport/grpc/handler/mocks"
)

func contextWithUserID(userID uuid.UUID) context.Context {
	md := metadata.Pairs("user_id", userID.String())
	return metadata.NewIncomingContext(context.Background(), md)
}

func TestCreateItem(t *testing.T) {
	userID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		req       *vaultpb.CreateItemRequest
		setupMock func(*mocks.MockVaultService)
		wantCode  codes.Code
		wantItem  bool
	}{
		{
			name: "success",
			ctx:  contextWithUserID(userID),
			req: &vaultpb.CreateItemRequest{
				Type:          vaultpb.ItemType_ITEM_TYPE_LOGIN,
				EncryptedData: []byte("encrypted"),
				Metadata:      map[string]string{"site": "example.com"},
			},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					CreateItem(gomock.Any(), userID, domain.ItemTypeLogin, []byte("encrypted"), gomock.Any()).
					Return(&domain.VaultItem{ID: itemID, UserID: userID, Type: domain.ItemTypeLogin}, nil)
			},
			wantCode: codes.OK,
			wantItem: true,
		},
		{
			name: "no user id",
			ctx:  context.Background(),
			req: &vaultpb.CreateItemRequest{
				Type:          vaultpb.ItemType_ITEM_TYPE_LOGIN,
				EncryptedData: []byte("encrypted"),
			},
			setupMock: func(mockService *mocks.MockVaultService) {},
			wantCode:  codes.Unauthenticated,
			wantItem:  false,
		},
		{
			name: "invalid type",
			ctx:  contextWithUserID(userID),
			req: &vaultpb.CreateItemRequest{
				Type:          vaultpb.ItemType_ITEM_TYPE_UNSPECIFIED,
				EncryptedData: []byte("encrypted"),
			},
			setupMock: func(mockService *mocks.MockVaultService) {},
			wantCode:  codes.InvalidArgument,
			wantItem:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockVaultService(ctrl)
			handler := NewVaultHandler(mockService)

			tt.setupMock(mockService)

			resp, err := handler.CreateItem(tt.ctx, tt.req)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.wantItem {
					assert.NotNil(t, resp.Item)
					assert.NotEmpty(t, resp.Item.Id)
				}
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
			}
		})
	}
}

func TestGetItem(t *testing.T) {
	userID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		req       *vaultpb.GetItemRequest
		setupMock func(*mocks.MockVaultService)
		wantCode  codes.Code
	}{
		{
			name: "success",
			ctx:  contextWithUserID(userID),
			req:  &vaultpb.GetItemRequest{Id: itemID.String()},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					GetItem(gomock.Any(), itemID, userID).
					Return(&domain.VaultItem{ID: itemID, UserID: userID}, nil)
			},
			wantCode: codes.OK,
		},
		{
			name: "not found",
			ctx:  contextWithUserID(userID),
			req:  &vaultpb.GetItemRequest{Id: itemID.String()},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					GetItem(gomock.Any(), itemID, userID).
					Return(nil, domain.ErrVaultItemNotFound)
			},
			wantCode: codes.NotFound,
		},
		{
			name:      "invalid id",
			ctx:       contextWithUserID(userID),
			req:       &vaultpb.GetItemRequest{Id: "invalid-uuid"},
			setupMock: func(mockService *mocks.MockVaultService) {},
			wantCode:  codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockVaultService(ctrl)
			handler := NewVaultHandler(mockService)

			tt.setupMock(mockService)

			resp, err := handler.GetItem(tt.ctx, tt.req)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Item)
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
			}
		})
	}
}

func TestListItems(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		req       *vaultpb.ListItemsRequest
		setupMock func(*mocks.MockVaultService)
		wantCode  codes.Code
		wantCount int
	}{
		{
			name: "success",
			ctx:  contextWithUserID(userID),
			req:  &vaultpb.ListItemsRequest{},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					ListItems(gomock.Any(), userID, nil).
					Return([]*domain.VaultItem{
						{ID: uuid.New(), UserID: userID},
						{ID: uuid.New(), UserID: userID},
					}, nil)
			},
			wantCode:  codes.OK,
			wantCount: 2,
		},
		{
			name: "empty",
			ctx:  contextWithUserID(userID),
			req:  &vaultpb.ListItemsRequest{},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					ListItems(gomock.Any(), userID, nil).
					Return([]*domain.VaultItem{}, nil)
			},
			wantCode:  codes.OK,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockVaultService(ctrl)
			handler := NewVaultHandler(mockService)

			tt.setupMock(mockService)

			resp, err := handler.ListItems(tt.ctx, tt.req)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Items, tt.wantCount)
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	userID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		req       *vaultpb.UpdateItemRequest
		setupMock func(*mocks.MockVaultService)
		wantCode  codes.Code
	}{
		{
			name: "success",
			ctx:  contextWithUserID(userID),
			req: &vaultpb.UpdateItemRequest{
				Id:            itemID.String(),
				EncryptedData: []byte("new data"),
			},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					UpdateItem(gomock.Any(), itemID, userID, []byte("new data"), gomock.Any()).
					Return(&domain.VaultItem{ID: itemID, UserID: userID}, nil)
			},
			wantCode: codes.OK,
		},
		{
			name: "not found",
			ctx:  contextWithUserID(userID),
			req: &vaultpb.UpdateItemRequest{
				Id:            itemID.String(),
				EncryptedData: []byte("new data"),
			},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					UpdateItem(gomock.Any(), itemID, userID, []byte("new data"), gomock.Any()).
					Return(nil, domain.ErrVaultItemNotFound)
			},
			wantCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockVaultService(ctrl)
			handler := NewVaultHandler(mockService)

			tt.setupMock(mockService)

			resp, err := handler.UpdateItem(tt.ctx, tt.req)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Item)
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	userID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		req       *vaultpb.DeleteItemRequest
		setupMock func(*mocks.MockVaultService)
		wantCode  codes.Code
	}{
		{
			name: "success",
			ctx:  contextWithUserID(userID),
			req:  &vaultpb.DeleteItemRequest{Id: itemID.String()},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					DeleteItem(gomock.Any(), itemID, userID).
					Return(nil)
			},
			wantCode: codes.OK,
		},
		{
			name: "not found",
			ctx:  contextWithUserID(userID),
			req:  &vaultpb.DeleteItemRequest{Id: itemID.String()},
			setupMock: func(mockService *mocks.MockVaultService) {
				mockService.EXPECT().
					DeleteItem(gomock.Any(), itemID, userID).
					Return(domain.ErrVaultItemNotFound)
			},
			wantCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockVaultService(ctrl)
			handler := NewVaultHandler(mockService)

			tt.setupMock(mockService)

			resp, err := handler.DeleteItem(tt.ctx, tt.req)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Message)
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
			}
		})
	}
}
