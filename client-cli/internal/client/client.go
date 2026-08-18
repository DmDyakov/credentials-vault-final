// Package client содержит gRPC клиент.
package client

import (
	"context"
	"fmt"
	"strings"

	authpb "credentials-vault/gen/go/auth/v1"
	vaultpb "credentials-vault/gen/go/vault/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

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

// ContextWithToken добавляет токен в контекст.
func (c *Client) ContextWithToken(ctx context.Context) context.Context {
	if c.config.Token == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.config.Token)
}

// Register регистрирует нового пользователя.
func (c *Client) Register(ctx context.Context, username, password string) error {
	resp, err := c.auth.Register(ctx, &authpb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	fmt.Printf("User registered: %s\n", resp.User.Username)
	return nil
}

// Login выполняет вход и сохраняет токен.
func (c *Client) Login(ctx context.Context, username, password string) error {
	var header metadata.MD

	resp, err := c.auth.Login(ctx, &authpb.LoginRequest{
		Username: username,
		Password: password,
	}, grpc.Header(&header))
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	values := header.Get("authorization")
	if len(values) == 0 {
		return fmt.Errorf("authorization token not found in response")
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	c.config.Token = token

	if err := c.config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Logged in as: %s\n", resp.User.Username)
	return nil
}
