// Package get содержит команды получения данных.
package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"credentials-vault/client-cli/internal/model"
)

//go:generate mockgen -source=get.go -destination=mocks/client_mock.go -package=mocks Client

// Client - интерфейс клиента CLI.
type Client interface {
	GetItem(ctx context.Context, id string) (*model.VaultItem, error)
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

			fmt.Printf("ID: %s\n", item.ID)
			fmt.Printf("Type: %s\n", item.Type)

			fmt.Println("\nMetadata:")
			for k, v := range item.Metadata {
				fmt.Printf("  %s: %s\n", k, v)
			}

			fmt.Println("\nSecret:")
			for k, v := range item.Secret {
				fmt.Printf("  %s: %s\n", k, v)
			}

			return nil
		},
	}
}
