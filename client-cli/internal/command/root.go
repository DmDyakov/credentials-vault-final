// Package command содержит команды CLI.
package command

import (
	"github.com/spf13/cobra"

	"credentials-vault/client-cli/internal/client"
	"credentials-vault/client-cli/internal/command/add"
	"credentials-vault/client-cli/internal/command/auth"
)

// NewRootCmd создаёт корневую команду.
func NewRootCmd(client *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Credentials Vault CLI",
		Long:  `CLI-клиент для управления приватными данными в Credentials Vault.`,
	}

	cmd.AddCommand(
		auth.NewRegisterCmd(client),
		auth.NewLoginCmd(client),
		add.NewAddCmd(client),
	)

	return cmd
}
