package client

import (
	"context"
	"fmt"

	vaultpb "credentials-vault/gen/go/vault/v1"

	"google.golang.org/grpc/metadata"
)

// AddLogin добавляет логин/пароль в хранилище.
func (c *Client) AddLogin(ctx context.Context, site, username, password string) error {
	data := fmt.Sprintf("username=%s\npassword=%s", username, password)

	metadata := map[string]string{
		"site": site,
	}

	resp, err := c.vault.CreateItem(c.withToken(ctx), &vaultpb.CreateItemRequest{
		Type:          vaultpb.ItemType_ITEM_TYPE_LOGIN,
		EncryptedData: []byte(data),
		Metadata:      metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to add login: %w", err)
	}

	fmt.Printf("Login added: %s\n", resp.Item.Id)
	return nil
}

// ListItems возвращает список элементов.
func (c *Client) ListItems(ctx context.Context) ([]*vaultpb.VaultItem, error) {
	resp, err := c.vault.ListItems(c.withToken(ctx), &vaultpb.ListItemsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	return resp.Items, nil
}

// GetItem возвращает элемент по ID.
func (c *Client) GetItem(ctx context.Context, id string) (*vaultpb.VaultItem, error) {
	resp, err := c.vault.GetItem(c.withToken(ctx), &vaultpb.GetItemRequest{
		Id: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return resp.Item, nil
}

// withToken добавляет токен в контекст.
func (c *Client) withToken(ctx context.Context) context.Context {
	if c.config.Token == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.config.Token)
}
