package handler

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	authpb "credentials-vault/gen/go/auth/v1"
	vaultpb "credentials-vault/gen/go/vault/v1"
	"credentials-vault/server/internal/domain"
)

// toProtoUser конвертирует доменную модель User в proto User.
func toProtoUser(user *domain.User) *authpb.User {
	return authpb.User_builder{
		Id:        proto.String(user.ID.String()),
		Username:  proto.String(user.Username),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}.Build()
}

// toDomainVaultItemType конвертирует proto ItemType в domain ItemType.
func toDomainVaultItemType(protoType vaultpb.ItemType) domain.ItemType {
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

// toProtoVaultItemType конвертирует domain ItemType в proto ItemType.
func toProtoVaultItemType(itemType domain.ItemType) vaultpb.ItemType {
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

// toProtoVaultItem конвертирует доменную модель VaultItem в proto VaultItem.
func toProtoVaultItem(item *domain.VaultItem) *vaultpb.VaultItem {
	return vaultpb.VaultItem_builder{
		Id:            proto.String(item.ID.String()),
		Type:          toProtoVaultItemType(item.Type).Enum(),
		EncryptedData: item.EncryptedData,
		Metadata:      item.Metadata,
		CreatedAt:     timestamppb.New(item.CreatedAt),
		UpdatedAt:     timestamppb.New(item.UpdatedAt),
	}.Build()
}
