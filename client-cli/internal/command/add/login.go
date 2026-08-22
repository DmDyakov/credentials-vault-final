// Package add содержит команды добавления данных.
package add

import (
	"context"

	"github.com/spf13/cobra"
)

// newLoginCmd создаёт команду add login.
func newLoginCmd(cl Client) *cobra.Command {
	var site, username, password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Добавить логин/пароль",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.AddLogin(context.Background(), site, username, password)
		},
	}

	cmd.Flags().StringVarP(&site, "site", "s", "", "Site URL")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	cobra.CheckErr(cmd.MarkFlagRequired("site"))
	cobra.CheckErr(cmd.MarkFlagRequired("username"))
	cobra.CheckErr(cmd.MarkFlagRequired("password"))

	return cmd
}
