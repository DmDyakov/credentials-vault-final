package vault

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"credentials-vault/internal/domain"
	"credentials-vault/internal/service/vault/mocks"
)

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name          string
		userID        uuid.UUID
		itemType      domain.ItemType
		encryptedData []byte
		metadata      map[string]string
		setupMock     func(*mocks.MockVaultRepository)
		wantErr       error
		wantItem      bool
	}{
		{
			name:          "success",
			userID:        uuid.New(),
			itemType:      domain.ItemTypeLogin,
			encryptedData: []byte("encrypted"),
			metadata:      map[string]string{"site": "example.com"},
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, item *domain.VaultItem) error {
						item.ID = uuid.New()
						return nil
					})
			},
			wantErr:  nil,
			wantItem: true,
		},
		{
			name:          "invalid type",
			userID:        uuid.New(),
			itemType:      domain.ItemType("INVALID"),
			encryptedData: []byte("encrypted"),
			setupMock:     func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:       domain.ErrInvalidItemType,
			wantItem:      false,
		},
		{
			name:          "empty data",
			userID:        uuid.New(),
			itemType:      domain.ItemTypeLogin,
			encryptedData: []byte{},
			setupMock:     func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:       domain.ErrEncryptedDataRequired,
			wantItem:      false,
		},
		{
			name:          "nil user id",
			userID:        uuid.Nil,
			itemType:      domain.ItemTypeLogin,
			encryptedData: []byte("encrypted"),
			setupMock:     func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:       domain.ErrUserIDRequired,
			wantItem:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockRepo := mocks.NewMockVaultRepository(ctrl)
			service := NewService(mockRepo)

			tt.setupMock(mockRepo)

			item, err := service.CreateItem(context.Background(), tt.userID, tt.itemType, tt.encryptedData, tt.metadata)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantItem {
				assert.NotNil(t, item)
				assert.NotEmpty(t, item.ID)
			} else {
				assert.Nil(t, item)
			}
		})
	}
}

func TestGetItem(t *testing.T) {
	userID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		userID    uuid.UUID
		setupMock func(*mocks.MockVaultRepository)
		wantErr   error
		wantItem  bool
	}{
		{
			name:   "success",
			id:     itemID,
			userID: userID,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					FindByID(gomock.Any(), itemID, userID).
					Return(&domain.VaultItem{ID: itemID, UserID: userID}, nil)
			},
			wantErr:  nil,
			wantItem: true,
		},
		{
			name:   "not found",
			id:     itemID,
			userID: userID,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					FindByID(gomock.Any(), itemID, userID).
					Return(nil, domain.ErrVaultItemNotFound)
			},
			wantErr:  domain.ErrVaultItemNotFound,
			wantItem: false,
		},
		{
			name:      "nil id",
			id:        uuid.Nil,
			userID:    userID,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:   domain.ErrVaultItemIDRequired,
			wantItem:  false,
		},
		{
			name:      "nil user id",
			id:        itemID,
			userID:    uuid.Nil,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:   domain.ErrUserIDRequired,
			wantItem:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockRepo := mocks.NewMockVaultRepository(ctrl)
			service := NewService(mockRepo)

			tt.setupMock(mockRepo)

			item, err := service.GetItem(context.Background(), tt.id, tt.userID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, item)
			}
		})
	}
}

func TestListItems(t *testing.T) {
	userID := uuid.New()
	loginType := domain.ItemTypeLogin
	invalidType := domain.ItemType("INVALID")

	tests := []struct {
		name      string
		userID    uuid.UUID
		itemType  *domain.ItemType
		setupMock func(*mocks.MockVaultRepository)
		wantErr   error
		wantCount int
	}{
		{
			name:     "success without filter",
			userID:   userID,
			itemType: nil,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					FindByUserID(gomock.Any(), userID, nil).
					Return([]*domain.VaultItem{
						{ID: uuid.New(), UserID: userID},
						{ID: uuid.New(), UserID: userID},
					}, nil)
			},
			wantErr:   nil,
			wantCount: 2,
		},
		{
			name:     "success with filter",
			userID:   userID,
			itemType: &loginType,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					FindByUserID(gomock.Any(), userID, &loginType).
					Return([]*domain.VaultItem{
						{ID: uuid.New(), UserID: userID, Type: loginType},
					}, nil)
			},
			wantErr:   nil,
			wantCount: 1,
		},
		{
			name:     "empty list",
			userID:   userID,
			itemType: nil,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					FindByUserID(gomock.Any(), userID, nil).
					Return([]*domain.VaultItem{}, nil)
			},
			wantErr:   nil,
			wantCount: 0,
		},
		{
			name:      "invalid filter",
			userID:    userID,
			itemType:  &invalidType,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:   domain.ErrInvalidItemType,
			wantCount: 0,
		},
		{
			name:      "nil user id",
			userID:    uuid.Nil,
			itemType:  nil,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:   domain.ErrUserIDRequired,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockRepo := mocks.NewMockVaultRepository(ctrl)
			service := NewService(mockRepo)

			tt.setupMock(mockRepo)

			items, err := service.ListItems(context.Background(), tt.userID, tt.itemType)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.Len(t, items, tt.wantCount)
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	userID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name          string
		id            uuid.UUID
		userID        uuid.UUID
		encryptedData []byte
		metadata      map[string]string
		setupMock     func(*mocks.MockVaultRepository)
		wantErr       error
		wantItem      bool
	}{
		{
			name:          "success",
			id:            itemID,
			userID:        userID,
			encryptedData: []byte("new data"),
			metadata:      map[string]string{"updated": "true"},
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr:  nil,
			wantItem: true,
		},
		{
			name:          "not found",
			id:            itemID,
			userID:        userID,
			encryptedData: []byte("new data"),
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(domain.ErrVaultItemNotFound)
			},
			wantErr:  domain.ErrVaultItemNotFound,
			wantItem: false,
		},
		{
			name:          "nil id",
			id:            uuid.Nil,
			userID:        userID,
			encryptedData: []byte("new data"),
			setupMock:     func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:       domain.ErrVaultItemIDRequired,
			wantItem:      false,
		},
		{
			name:          "empty data",
			id:            itemID,
			userID:        userID,
			encryptedData: []byte{},
			setupMock:     func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:       domain.ErrEncryptedDataRequired,
			wantItem:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockRepo := mocks.NewMockVaultRepository(ctrl)
			service := NewService(mockRepo)

			tt.setupMock(mockRepo)

			item, err := service.UpdateItem(context.Background(), tt.id, tt.userID, tt.encryptedData, tt.metadata)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, item)
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	userID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		userID    uuid.UUID
		setupMock func(*mocks.MockVaultRepository)
		wantErr   error
	}{
		{
			name:   "success",
			id:     itemID,
			userID: userID,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					SoftDelete(gomock.Any(), itemID, userID).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "not found",
			id:     itemID,
			userID: userID,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {
				mockRepo.EXPECT().
					SoftDelete(gomock.Any(), itemID, userID).
					Return(domain.ErrVaultItemNotFound)
			},
			wantErr: domain.ErrVaultItemNotFound,
		},
		{
			name:      "nil id",
			id:        uuid.Nil,
			userID:    userID,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:   domain.ErrVaultItemIDRequired,
		},
		{
			name:      "nil user id",
			id:        itemID,
			userID:    uuid.Nil,
			setupMock: func(mockRepo *mocks.MockVaultRepository) {},
			wantErr:   domain.ErrUserIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockRepo := mocks.NewMockVaultRepository(ctrl)
			service := NewService(mockRepo)

			tt.setupMock(mockRepo)

			err := service.DeleteItem(context.Background(), tt.id, tt.userID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
