// Package add содержит команды добавления данных.
package add

import (
	"context"

	"github.com/spf13/cobra"
)

// newCredentialsCmd создаёт команду add credentials.
func newCredentialsCmd(cl Client) *cobra.Command {
	var site, username, password string

	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Добавить учётные данные",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.AddCredentials(context.Background(), site, username, password)
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
