package handler

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	vaultpb "credentials-vault/gen/go/vault/v1"
	"credentials-vault/internal/domain"
)

//go:generate mockgen -source=vault.go -destination=mocks/vault_service_mock.go -package=mocks VaultService
type VaultService interface {
	CreateItem(ctx context.Context, userID uuid.UUID, itemType domain.ItemType, encryptedData []byte, metadata map[string]string) (*domain.VaultItem, error)
	GetItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.VaultItem, error)
	ListItems(ctx context.Context, userID uuid.UUID, itemType *domain.ItemType) ([]*domain.VaultItem, error)
	UpdateItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, encryptedData []byte, metadata map[string]string) (*domain.VaultItem, error)
	DeleteItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

// VaultHandler - gRPC обработчик хранилища.
type VaultHandler struct {
	vaultpb.UnimplementedVaultServiceServer
	vaultService VaultService
}

// NewVaultHandler создаёт новый обработчик хранилища.
func NewVaultHandler(vaultService VaultService) *VaultHandler {
	return &VaultHandler{
		vaultService: vaultService,
	}
}

// CreateItem обрабатывает запрос на создание элемента.
func (h *VaultHandler) CreateItem(ctx context.Context, req *vaultpb.CreateItemRequest) (*vaultpb.CreateItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userID, err := getUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	itemType := toDomainVaultItemType(req.Type)
	if itemType == "" {
		return nil, status.Error(codes.InvalidArgument, "item type is required")
	}

	item, err := h.vaultService.CreateItem(ctx, userID, itemType, req.EncryptedData, req.Metadata)
	if err != nil {
		return nil, mapError(err)
	}

	return &vaultpb.CreateItemResponse{
		Item: toProtoVaultItem(item),
	}, nil
}

// GetItem обрабатывает запрос на получение элемента по ID.
func (h *VaultHandler) GetItem(ctx context.Context, req *vaultpb.GetItemRequest) (*vaultpb.GetItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userID, err := getUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item id")
	}

	item, err := h.vaultService.GetItem(ctx, itemID, userID)
	if err != nil {
		return nil, mapError(err)
	}

	return &vaultpb.GetItemResponse{
		Item: toProtoVaultItem(item),
	}, nil
}

// UpdateItem обрабатывает запрос на обновление элемента.
func (h *VaultHandler) UpdateItem(ctx context.Context, req *vaultpb.UpdateItemRequest) (*vaultpb.UpdateItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userID, err := getUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item id")
	}

	item, err := h.vaultService.UpdateItem(ctx, itemID, userID, req.EncryptedData, req.Metadata)
	if err != nil {
		return nil, mapError(err)
	}

	return &vaultpb.UpdateItemResponse{
		Item: toProtoVaultItem(item),
	}, nil
}

// DeleteItem обрабатывает запрос на мягкое удаление элемента.
func (h *VaultHandler) DeleteItem(ctx context.Context, req *vaultpb.DeleteItemRequest) (*vaultpb.DeleteItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userID, err := getUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item id")
	}

	if err := h.vaultService.DeleteItem(ctx, itemID, userID); err != nil {
		return nil, mapError(err)
	}

	return &vaultpb.DeleteItemResponse{
		Message: "Item deleted successfully",
	}, nil
}

// ListItems обрабатывает запрос на получение списка элементов.
func (h *VaultHandler) ListItems(ctx context.Context, req *vaultpb.ListItemsRequest) (*vaultpb.ListItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userID, err := getUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	var itemType *domain.ItemType
	if req.Type != vaultpb.ItemType_ITEM_TYPE_UNSPECIFIED {
		t := toDomainVaultItemType(req.Type)
		if t == "" {
			return nil, status.Error(codes.InvalidArgument, "invalid item type")
		}
		itemType = &t
	}

	items, err := h.vaultService.ListItems(ctx, userID, itemType)
	if err != nil {
		return nil, mapError(err)
	}

	protoItems := make([]*vaultpb.VaultItem, 0, len(items))
	for _, item := range items {
		protoItems = append(protoItems, toProtoVaultItem(item))
	}

	return &vaultpb.ListItemsResponse{
		Items: protoItems,
	}, nil
}
