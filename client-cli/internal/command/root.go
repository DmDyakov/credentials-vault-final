// Package command содержит команды CLI.
package command

import (
	"github.com/spf13/cobra"

	"credentials-vault/client-cli/internal/client"
)

// NewRootCmd создаёт корневую команду.
func NewRootCmd(cl *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Credentials Vault CLI",
		Long:  `CLI-клиент для управления приватными данными в Credentials Vault.`,
	}

	cmd.AddCommand(
		newRegisterCmd(cl),
		newLoginCmd(cl),
	)

	return cmd
}
