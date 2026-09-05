package client

import (
	"context"
	"encoding/json"
	"fmt"

	"credentials-vault/client-cli/internal/crypto"
	"credentials-vault/client-cli/internal/model"
	"credentials-vault/client-cli/internal/session"
	"credentials-vault/client-cli/internal/validate"
	vaultpb "credentials-vault/gen/go/vault/v1"

	"google.golang.org/protobuf/proto"
)

// AddCredentials добавляет учётные данные в хранилище.
func (c *Client) AddCredentials(ctx context.Context, site, username, password string) error {
	if err := validate.Credentials(site, username, password); err != nil {
		return err
	}

	key, err := session.Get()
	if err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	secret := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	encryptedData, err := crypto.Encrypt(jsonData, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}

	req := vaultpb.CreateItemRequest_builder{
		Type:          vaultpb.ItemType_ITEM_TYPE_LOGIN.Enum(),
		EncryptedData: encryptedData,
		Metadata: map[string]string{
			"site": site,
		},
	}.Build()

	resp, err := c.vault.CreateItem(c.withAuthToken(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to add credentials: %w", err)
	}

	fmt.Printf("Credentials added: %s\n", resp.GetItem().GetId())
	return nil
}

// AddCard добавляет карту в хранилище.
func (c *Client) AddCard(ctx context.Context, brand, bank, number, holder, expiry, cvv string) error {
	if err := validate.Card(brand, bank, number, holder, expiry, cvv); err != nil {
		return err
	}

	key, err := session.Get()
	if err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	secret := map[string]string{
		"number": number,
		"holder": holder,
		"expiry": expiry,
		"cvv":    cvv,
	}

	jsonData, err := json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	encryptedData, err := crypto.Encrypt(jsonData, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}

	if len(number) < 4 {
		return fmt.Errorf("card number must be at least 4 characters")
	}
	last4 := number[len(number)-4:]

	req := vaultpb.CreateItemRequest_builder{
		Type:          vaultpb.ItemType_ITEM_TYPE_CARD.Enum(),
		EncryptedData: encryptedData,
		Metadata: map[string]string{
			"brand": brand,
			"bank":  bank,
			"last4": last4,
		},
	}.Build()

	resp, err := c.vault.CreateItem(c.withAuthToken(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to add card: %w", err)
	}

	fmt.Printf("Card added: %s\n", resp.GetItem().GetId())
	return nil
}

// GetItem возвращает элемент по ID.
func (c *Client) GetItem(ctx context.Context, id string) (*model.VaultItem, error) {
	key, err := session.Get()
	if err != nil {
		return nil, fmt.Errorf("not logged in: %w", err)
	}

	req := vaultpb.GetItemRequest_builder{
		Id: proto.String(id),
	}.Build()

	resp, err := c.vault.GetItem(c.withAuthToken(ctx), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	jsonData, err := crypto.Decrypt(resp.GetItem().GetEncryptedData(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	var secret map[string]string
	if err := json.Unmarshal(jsonData, &secret); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &model.VaultItem{
		ID:        resp.GetItem().GetId(),
		Type:      resp.GetItem().GetType().String(),
		Secret:    secret,
		Metadata:  resp.GetItem().GetMetadata(),
		CreatedAt: resp.GetItem().GetCreatedAt().AsTime(),
		UpdatedAt: resp.GetItem().GetUpdatedAt().AsTime(),
	}, nil
}

// ListItems возвращает список всех элементов.
func (c *Client) ListItems(ctx context.Context) ([]*model.ListVaultItem, error) {
	req := &vaultpb.ListItemsRequest{}

	resp, err := c.vault.ListItems(c.withAuthToken(ctx), req)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	items := make([]*model.ListVaultItem, 0, len(resp.GetItems()))
	for _, pb := range resp.GetItems() {
		items = append(items, &model.ListVaultItem{
			ID:        pb.GetId(),
			Type:      pb.GetType().String(),
			Metadata:  pb.GetMetadata(),
			CreatedAt: pb.GetCreatedAt().AsTime(),
			UpdatedAt: pb.GetUpdatedAt().AsTime(),
		})
	}

	return items, nil
}
