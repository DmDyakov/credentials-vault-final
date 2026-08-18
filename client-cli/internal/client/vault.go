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

// withToken добавляет токен в контекст.
func (c *Client) withToken(ctx context.Context) context.Context {
	if c.config.Token == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.config.Token)
}
