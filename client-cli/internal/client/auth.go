package client

import (
	"context"
	"fmt"
	"strings"

	authpb "credentials-vault/gen/go/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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
