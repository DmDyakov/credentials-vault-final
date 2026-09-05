// Package add содержит команды добавления данных.
package add

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// newCredentialsCmd создаёт команду add credentials.
func newCredentialsCmd(cl Client) *cobra.Command {
	var site, username string

	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Добавить учётные данные",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Password: ")
			password, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Println()

			return cl.AddCredentials(context.Background(), site, username, string(password))
		},
	}

	cmd.Flags().StringVarP(&site, "site", "s", "", "Site URL")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cobra.CheckErr(cmd.MarkFlagRequired("site"))
	cobra.CheckErr(cmd.MarkFlagRequired("username"))

	return cmd
}
