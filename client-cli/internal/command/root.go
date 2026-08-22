// Package command содержит команды CLI.
package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"credentials-vault/client-cli/internal/client"
	"credentials-vault/client-cli/internal/command/add"
	"credentials-vault/client-cli/internal/command/auth"
	"credentials-vault/client-cli/internal/command/system"
)

// NewRootCmd создаёт корневую команду.
func NewRootCmd(cl *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Credentials Vault CLI",
		Long:  `CLI-клиент для управления приватными данными в Credentials Vault.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch cmd.Name() {
			case "init", "version", "help", "completion":
				return nil
			}

			if cl == nil {
				return fmt.Errorf("run 'vault init' first")
			}

			return nil
		},
	}

	cmd.AddCommand(
		auth.NewRegisterCmd(cl),
		auth.NewLoginCmd(cl),
		add.NewAddCmd(cl),
		newListCmd(cl),
		newGetCmd(cl),
		system.NewInitCmd(),
		system.NewVersionCmd(),
	)

	return cmd
}
