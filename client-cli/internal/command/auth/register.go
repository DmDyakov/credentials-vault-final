// Package auth содержит команды аутентификации.
package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

//go:generate mockgen -source=register.go -destination=mocks/client_mock.go -package=mocks Client

// Client - интерфейс клиента аутентификации.
type Client interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) error
	Logout() error
}

// NewRegisterCmd создаёт команду register.
func NewRegisterCmd(cl Client) *cobra.Command {
	var username string

	cmd := &cobra.Command{
		Use:   "register",
		Short: "Регистрация нового пользователя",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Password: ")
			password, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Println()

			fmt.Print("Confirm password: ")
			confirm, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("failed to read confirm password: %w", err)
			}
			fmt.Println()

			if string(password) != string(confirm) {
				return fmt.Errorf("passwords do not match")
			}

			return cl.Register(context.Background(), username, string(password))
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cobra.CheckErr(cmd.MarkFlagRequired("username"))

	return cmd
}
