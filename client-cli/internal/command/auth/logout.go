package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewLogoutCmd(cl Client) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Выйти из системы",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.Logout(); err != nil {
				return err
			}
			fmt.Println("Logged out")
			return nil
		},
	}
}
