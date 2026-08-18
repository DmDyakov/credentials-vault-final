// Package auth содержит команды аутентификации.
package auth

import (
	"context"

	"github.com/spf13/cobra"
)

//go:generate mockgen -source=register.go -destination=mocks/client_mock.go -package=mocks Client

// Client - интерфейс клиента аутентификации.
type Client interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) error
}

// NewRegisterCmd создаёт команду register.
func NewRegisterCmd(cl Client) *cobra.Command {
	var username, password string

	cmd := &cobra.Command{
		Use:   "register",
		Short: "Регистрация нового пользователя",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.MarkFlagRequired("username"); err != nil {
				return err
			}
			if err := cmd.MarkFlagRequired("password"); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.Register(context.Background(), username, password)
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")

	return cmd
}
