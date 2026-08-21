package client

import (
	"context"
	"fmt"

	vaultpb "credentials-vault/gen/go/vault/v1"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// AddLogin добавляет логин/пароль в хранилище.
func (c *Client) AddLogin(ctx context.Context, site, username, password string) error {
	data := fmt.Sprintf("username=%s\npassword=%s", username, password)

	itemMetadata := map[string]string{
		"site": site,
	}

	req := vaultpb.CreateItemRequest_builder{
		Type:          vaultpb.ItemType_ITEM_TYPE_LOGIN.Enum(),
		EncryptedData: []byte(data),
		Metadata:      itemMetadata,
	}.Build()

	resp, err := c.vault.CreateItem(c.withToken(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to add login: %w", err)
	}

	fmt.Printf("Login added: %s\n", resp.GetItem().GetId())
	return nil
}

// ListItems возвращает список элементов.
func (c *Client) ListItems(ctx context.Context) ([]*vaultpb.VaultItem, error) {
	req := &vaultpb.ListItemsRequest{}

	resp, err := c.vault.ListItems(c.withToken(ctx), req)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	return resp.GetItems(), nil
}

// GetItem возвращает элемент по ID.
func (c *Client) GetItem(ctx context.Context, id string) (*vaultpb.VaultItem, error) {
	req := vaultpb.GetItemRequest_builder{
		Id: proto.String(id),
	}.Build()

	resp, err := c.vault.GetItem(c.withToken(ctx), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return resp.GetItem(), nil
}

// withToken добавляет токен в контекст.
func (c *Client) withToken(ctx context.Context) context.Context {
	if c.config.Token == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.config.Token)
}
