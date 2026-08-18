// Package add содержит команды добавления данных.
package add

import (
	"github.com/spf13/cobra"

	"credentials-vault/client-cli/internal/client"
)

// NewAddCmd создаёт команду add.
func NewAddCmd(cl *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Добавление данных",
	}

	cmd.AddCommand(
		newLoginCmd(cl),
	)

	return cmd
}
