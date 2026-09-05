package handler

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	vaultpb "credentials-vault/gen/go/vault/v1"
	"credentials-vault/server/internal/domain"
)

//go:generate mockgen -source=vault.go -destination=mocks/vault_service_mock.gen.go -package=mocks VaultService
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
	logger       *zap.Logger
}

// NewVaultHandler создаёт новый обработчик хранилища.
func NewVaultHandler(vaultService VaultService, logger *zap.Logger) *VaultHandler {
	return &VaultHandler{
		vaultService: vaultService,
		logger:       logger,
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

	itemType := toDomainVaultItemType(req.GetType())
	if itemType == "" {
		return nil, status.Error(codes.InvalidArgument, "item type is required")
	}

	item, err := h.vaultService.CreateItem(ctx, userID, itemType, req.GetEncryptedData(), req.GetMetadata())
	if err != nil {
		return nil, mapError(err, h.logger)
	}

	return vaultpb.CreateItemResponse_builder{
		Item: toProtoVaultItem(item),
	}.Build(), nil
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

	itemID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item id")
	}

	item, err := h.vaultService.GetItem(ctx, itemID, userID)
	if err != nil {
		return nil, mapError(err, h.logger)
	}

	return vaultpb.GetItemResponse_builder{
		Item: toProtoVaultItem(item),
	}.Build(), nil
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

	itemID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item id")
	}

	item, err := h.vaultService.UpdateItem(ctx, itemID, userID, req.GetEncryptedData(), req.GetMetadata())
	if err != nil {
		return nil, mapError(err, h.logger)
	}

	return vaultpb.UpdateItemResponse_builder{
		Item: toProtoVaultItem(item),
	}.Build(), nil
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

	itemID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item id")
	}

	if err := h.vaultService.DeleteItem(ctx, itemID, userID); err != nil {
		return nil, mapError(err, h.logger)
	}

	return vaultpb.DeleteItemResponse_builder{
		Message: proto.String("Item deleted successfully"),
	}.Build(), nil
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
	if req.GetType() != vaultpb.ItemType_ITEM_TYPE_UNSPECIFIED {
		t := toDomainVaultItemType(req.GetType())
		if t == "" {
			return nil, status.Error(codes.InvalidArgument, "invalid item type")
		}
		itemType = &t
	}

	items, err := h.vaultService.ListItems(ctx, userID, itemType)
	if err != nil {
		return nil, mapError(err, h.logger)
	}

	protoItems := make([]*vaultpb.VaultItem, 0, len(items))
	for _, item := range items {
		protoItems = append(protoItems, toProtoVaultItem(item))
	}

	return vaultpb.ListItemsResponse_builder{
		Items: protoItems,
	}.Build(), nil
}
