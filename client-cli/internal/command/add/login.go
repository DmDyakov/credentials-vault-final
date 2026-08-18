// Package add содержит команды добавления данных.
package add

import (
	"context"

	"github.com/spf13/cobra"

	"credentials-vault/client-cli/internal/client"
)

// newAddLoginCmd создаёт команду add login.
func newLoginCmd(cl *client.Client) *cobra.Command {
	var site, username, password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Добавить логин/пароль",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.MarkFlagRequired("site"); err != nil {
				return err
			}
			if err := cmd.MarkFlagRequired("username"); err != nil {
				return err
			}
			if err := cmd.MarkFlagRequired("password"); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.AddLogin(context.Background(), site, username, password)
		},
	}

	cmd.Flags().StringVarP(&site, "site", "s", "", "Site URL")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")

	return cmd
}
