// Package system содержит системные команды.
package system

import (
	"fmt"

	"github.com/spf13/cobra"

	clientconfig "credentials-vault/client-cli/internal/config"
)

// NewInitCmd создаёт команду init.
func NewInitCmd() *cobra.Command {
	var serverAddress, caFile string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Создать конфигурационный файл",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.MarkFlagRequired("server"); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := clientconfig.New()
			if err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}

			cfg.ServerAddress = serverAddress
			cfg.CAFile = caFile

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}

			fmt.Println("Config created at ~/.credentials-vault/config.json")

			return nil
		},
	}

	cmd.Flags().StringVarP(&serverAddress, "server", "s", "", "Server address (e.g., localhost:9090)")
	cmd.Flags().StringVar(&caFile, "ca-file", "", "Path to CA certificate file (dev only)")

	return cmd
}
