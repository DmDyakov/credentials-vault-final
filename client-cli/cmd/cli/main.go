// Package main реализует CLI-клиент Credentials Vault.
package main

import (
	"fmt"
	"os"

	"credentials-vault/client-cli/internal/client"
	"credentials-vault/client-cli/internal/command"
	"credentials-vault/client-cli/internal/config"
)

// run выполняет CLI.
func run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var cl *client.Client
	if cfg.Valid() {
		cl, err = client.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer cl.Close()
	}

	rootCmd := command.NewRootCmd(cl)
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
