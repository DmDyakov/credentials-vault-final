// Package client содержит gRPC клиент.
package client

import (
	"context"
	"crypto/tls"
	"fmt"

	authpb "credentials-vault/gen/go/auth/v1"
	vaultpb "credentials-vault/gen/go/vault/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"credentials-vault/client-cli/internal/config"
)

// Client — gRPC клиент Credentials Vault.
type Client struct {
	conn   *grpc.ClientConn
	auth   authpb.AuthServiceClient
	vault  vaultpb.VaultServiceClient
	config *config.Config
}

// New создаёт новый gRPC клиент.
func New(cfg *config.Config) (*Client, error) {
	var creds credentials.TransportCredentials

	if cfg.CAFile == "" {
		creds = credentials.NewTLS(&tls.Config{})
	} else {
		var err error
		creds, err = credentials.NewClientTLSFromFile(cfg.CAFile, "")
		if err != nil {
			return nil, fmt.Errorf("failed to load CA certificate: %w", err)
		}
	}

	conn, err := grpc.NewClient(
		cfg.ServerAddress,
		grpc.WithTransportCredentials(creds),
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

// withAuthToken добавляет JWT токен в gRPC метаданные.
func (c *Client) withAuthToken(ctx context.Context) context.Context {
	if c.config.Token == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.config.Token)
}
