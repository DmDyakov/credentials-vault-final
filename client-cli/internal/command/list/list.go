// Package list содержит команды получения списка данных.
package list

import (
	"context"
	vaultpb "credentials-vault/gen/go/vault/v1"
	"fmt"

	"github.com/spf13/cobra"
)

//go:generate mockgen -source=list.go -destination=mocks/client_mock.go -package=mocks Client

// Client - интерфейс клиента CLI.
type Client interface {
	ListItems(ctx context.Context) ([]*vaultpb.VaultItem, error)
}

// NewListCmd создаёт команду list.
func NewListCmd(cl Client) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Список элементов",
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := cl.ListItems(context.Background())
			if err != nil {
				return err
			}

			if len(items) == 0 {
				fmt.Println("No items found")
				return nil
			}

			fmt.Println("ID\t\t\t\t\tTYPE\t\tCREATED")
			for _, item := range items {
				fmt.Printf("%s\t%s\t%s\n",
					item.GetId(),
					item.GetType().String(),
					item.GetCreatedAt().AsTime().Format("2006-01-02 15:04:05"),
				)
			}

			return nil
		},
	}
}
