// Package client содержит gRPC клиент.
package client

import (
	"fmt"

	authpb "credentials-vault/gen/go/auth/v1"
	vaultpb "credentials-vault/gen/go/vault/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"credentials-vault/client-cli/internal/config"
)

// Client - gRPC клиент Credentials Vault.
type Client struct {
	conn   *grpc.ClientConn
	auth   authpb.AuthServiceClient
	vault  vaultpb.VaultServiceClient
	config *config.Config
}

// New создаёт новый gRPC клиент.
func New(cfg *config.Config) (*Client, error) {
	conn, err := grpc.NewClient(
		cfg.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return &Client{
		conn:   conn,
		auth:   authpb.NewAuthServiceClient(conn),
		vault:  vaultpb.NewVaultServiceClient(conn),
		config: cfg,
	}, nil
}

// Close закрывает соединение.
func (c *Client) Close() error {
	return c.conn.Close()
}
