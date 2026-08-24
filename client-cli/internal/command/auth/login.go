// Package auth содержит команды аутентификации.
package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewLoginCmd создаёт команду login.
func NewLoginCmd(cl Client) *cobra.Command {
	var username string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Вход в систему",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Password: ")
			password, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Println()

			return cl.Login(context.Background(), username, string(password))
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cobra.CheckErr(cmd.MarkFlagRequired("username"))

	return cmd
}
