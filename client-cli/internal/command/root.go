// Package command содержит команды CLI.
package command

import (
	"context"

	"github.com/spf13/cobra"

	"credentials-vault/client-cli/internal/command/add"
	"credentials-vault/client-cli/internal/command/auth"
	"credentials-vault/client-cli/internal/command/info"

	vaultpb "credentials-vault/gen/go/vault/v1"
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

// NewRootCmd создаёт корневую команду.
func NewRootCmd(cl Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Credentials Vault CLI",
		Long:  `CLI-клиент для управления приватными данными в Credentials Vault.`,
	}

	cmd.AddCommand(
		auth.NewRegisterCmd(cl),
		auth.NewLoginCmd(cl),
		add.NewAddCmd(cl),
		newListCmd(cl),
		newGetCmd(cl),
		info.NewVersionCmd(),
	)

	return cmd
}
