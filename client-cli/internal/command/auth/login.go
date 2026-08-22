// Package auth содержит команды аутентификации.
package auth

import (
	"context"

	"github.com/spf13/cobra"
)

// NewLoginCmd создаёт команду login.
func NewLoginCmd(cl Client) *cobra.Command {
	var username, password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Вход в систему",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.Login(context.Background(), username, password)
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	cobra.CheckErr(cmd.MarkFlagRequired("username"))
	cobra.CheckErr(cmd.MarkFlagRequired("password"))

	return cmd
}
