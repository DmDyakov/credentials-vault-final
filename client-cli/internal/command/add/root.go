// Package add содержит команды добавления данных.
package add

import (
	"context"

	"github.com/spf13/cobra"
)

//go:generate mockgen -source=root.go -destination=mocks/cli_client_mock.go -package=mocks Client

// Client - интерфейс клиента CLI.
type Client interface {
	AddLogin(ctx context.Context, site, username, password string) error
}

// NewAddCmd создаёт команду add.
func NewAddCmd(cl Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Добавление данных",
	}

	cmd.AddCommand(
		newLoginCmd(cl),
	)

	return cmd
}
