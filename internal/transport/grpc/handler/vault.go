package handler

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	vaultpb "credentials-vault/gen/go/vault/v1"
	"credentials-vault/internal/domain"
)

//go:generate mockgen -source=vault.go -destination=mocks/vault_service_mock.go -package=mocks VaultService
type VaultService interface {
	CreateItem(ctx context.Context, userID uuid.UUID, itemType domain.ItemType, encryptedData []byte, metadata map[string]string) (*domain.VaultItem, error)
	ListItems(ctx context.Context, userID uuid.UUID, itemType *domain.ItemType) ([]*domain.VaultItem, error)
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

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	itemType := protoItemTypeToDomain(req.Type)
	if itemType == "" {
		return nil, status.Error(codes.InvalidArgument, "item type is required")
	}

	item, err := h.vaultService.CreateItem(ctx, userID, itemType, req.EncryptedData, req.Metadata)
	if err != nil {
		return nil, mapVaultError(err)
	}

	return &vaultpb.CreateItemResponse{
		Item: domainVaultItemToProto(item),
	}, nil
}

// ListItems обрабатывает запрос на получение списка элементов.
func (h *VaultHandler) ListItems(ctx context.Context, req *vaultpb.ListItemsRequest) (*vaultpb.ListItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var itemType *domain.ItemType
	if req.Type != vaultpb.ItemType_ITEM_TYPE_UNSPECIFIED {
		t := protoItemTypeToDomain(req.Type)
		if t == "" {
			return nil, status.Error(codes.InvalidArgument, "invalid item type")
		}
		itemType = &t
	}

	items, err := h.vaultService.ListItems(ctx, userID, itemType)
	if err != nil {
		return nil, mapVaultError(err)
	}

	protoItems := make([]*vaultpb.VaultItem, 0, len(items))
	for _, item := range items {
		protoItems = append(protoItems, domainVaultItemToProto(item))
	}

	return &vaultpb.ListItemsResponse{
		Items: protoItems,
	}, nil
}

// getUserIDFromContext извлекает user_id из gRPC метаданных.
func getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("user_id")
	if len(values) == 0 {
		return uuid.Nil, status.Error(codes.Unauthenticated, "user_id is not provided")
	}

	userID, err := uuid.Parse(values[0])
	if err != nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid user_id")
	}

	return userID, nil
}

// protoItemTypeToDomain конвертирует proto ItemType в domain ItemType.
func protoItemTypeToDomain(protoType vaultpb.ItemType) domain.ItemType {
	switch protoType {
	case vaultpb.ItemType_ITEM_TYPE_LOGIN:
		return domain.ItemTypeLogin
	case vaultpb.ItemType_ITEM_TYPE_CARD:
		return domain.ItemTypeCard
	case vaultpb.ItemType_ITEM_TYPE_TEXT:
		return domain.ItemTypeText
	case vaultpb.ItemType_ITEM_TYPE_BINARY:
		return domain.ItemTypeBinary
	default:
		return ""
	}
}

// domainItemTypeToProto конвертирует domain ItemType в proto ItemType.
func domainItemTypeToProto(itemType domain.ItemType) vaultpb.ItemType {
	switch itemType {
	case domain.ItemTypeLogin:
		return vaultpb.ItemType_ITEM_TYPE_LOGIN
	case domain.ItemTypeCard:
		return vaultpb.ItemType_ITEM_TYPE_CARD
	case domain.ItemTypeText:
		return vaultpb.ItemType_ITEM_TYPE_TEXT
	case domain.ItemTypeBinary:
		return vaultpb.ItemType_ITEM_TYPE_BINARY
	default:
		return vaultpb.ItemType_ITEM_TYPE_UNSPECIFIED
	}
}

// domainVaultItemToProto конвертирует domain VaultItem в proto VaultItem.
func domainVaultItemToProto(item *domain.VaultItem) *vaultpb.VaultItem {
	return &vaultpb.VaultItem{
		Id:            item.ID.String(),
		Type:          domainItemTypeToProto(item.Type),
		EncryptedData: item.EncryptedData,
		Metadata:      item.Metadata,
		CreatedAt:     timestamppb.New(item.CreatedAt),
		UpdatedAt:     timestamppb.New(item.UpdatedAt),
	}
}

// mapVaultError маппит доменные ошибки в gRPC статусы.
func mapVaultError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidItemType):
		return status.Error(codes.InvalidArgument, "invalid item type")
	case errors.Is(err, domain.ErrEncryptedDataRequired):
		return status.Error(codes.InvalidArgument, "encrypted data is required")
	case errors.Is(err, domain.ErrUserIDRequired):
		return status.Error(codes.Unauthenticated, "user id is required")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
