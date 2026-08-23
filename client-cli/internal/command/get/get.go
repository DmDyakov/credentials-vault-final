// Package get содержит команды получения данных.
package get

import (
	"context"
	vaultpb "credentials-vault/gen/go/vault/v1"
	"fmt"

	"github.com/spf13/cobra"
)

//go:generate mockgen -source=get.go -destination=mocks/client_mock.go -package=mocks Client

// Client - интерфейс клиента CLI.
type Client interface {
	GetItem(ctx context.Context, id string) (*vaultpb.VaultItem, error)
}

// NewGetCmd создаёт команду get.
func NewGetCmd(cl Client) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Получить элемент по ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := cl.GetItem(context.Background(), args[0])
			if err != nil {
				return err
			}

			fmt.Printf("ID: %s\n", item.GetId())
			fmt.Printf("Type: %s\n", item.GetType().String())
			fmt.Printf("Data: %s\n", string(item.GetEncryptedData()))
			fmt.Printf("Metadata:\n")
			for k, v := range item.GetMetadata() {
				fmt.Printf("  %s: %s\n", k, v)
			}
			fmt.Printf("Created: %s\n", item.GetCreatedAt().AsTime().Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated: %s\n", item.GetUpdatedAt().AsTime().Format("2006-01-02 15:04:05"))

			return nil
		},
	}
}
