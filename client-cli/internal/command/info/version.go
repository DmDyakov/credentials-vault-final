// Package info содержит информационные команды.
package info

import (
	"fmt"

	"github.com/spf13/cobra"

	"credentials-vault/pkg/buildinfo"
)

// NewVersionCmd создаёт команду version.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Информация о версии",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", buildinfo.Version)
			fmt.Printf("Build date: %s\n", buildinfo.Date)
			fmt.Printf("Commit: %s\n", buildinfo.Commit)
		},
	}
}
