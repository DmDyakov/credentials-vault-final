// Package auth содержит команды аутентификации.
package auth

import (
	"context"

	"github.com/spf13/cobra"

	"credentials-vault/client-cli/internal/client"
)

// NewLoginCmd создаёт команду login.
func NewLoginCmd(cl *client.Client) *cobra.Command {
	var username, password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Вход в систему",
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
			return cl.Login(context.Background(), username, password)
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")

	return cmd
}
