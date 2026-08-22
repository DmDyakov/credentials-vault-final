package command

import (
	"context"
	vaultpb "credentials-vault/gen/go/vault/v1"
	"fmt"

	"github.com/spf13/cobra"
)

//go:generate mockgen -source=root.go -destination=mocks/client_mock.go -package=mocks Client

// Client - интерфейс клиента CLI.
type Client interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) error
	AddLogin(ctx context.Context, site, username, password string) error
	ListItems(ctx context.Context) ([]*vaultpb.VaultItem, error)
	GetItem(ctx context.Context, id string) (*vaultpb.VaultItem, error)
}

// newListCmd создаёт команду list.
func newListCmd(cl Client) *cobra.Command {
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
