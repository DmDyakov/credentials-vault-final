package client

import (
	"context"
	"fmt"
	"strings"

	"credentials-vault/client-cli/internal/crypto"
	"credentials-vault/client-cli/internal/session"
	"credentials-vault/client-cli/internal/validate"
	authpb "credentials-vault/gen/go/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// Register регистрирует нового пользователя.
func (c *Client) Register(ctx context.Context, username, password string) error {
	if err := validate.Username(username); err != nil {
		return err
	}

	if err := validate.Password(password); err != nil {
		return err
	}

	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	req := authpb.RegisterRequest_builder{
		Username: proto.String(username),
		Password: proto.String(password),
		Salt:     salt,
	}.Build()

	resp, err := c.auth.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	key, err := crypto.DeriveKey(password, salt)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	if err := session.Save(key, session.DefaultTTL); err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	fmt.Printf("User registered: %s\n", resp.GetUser().GetUsername())
	return nil
}

// Login выполняет вход.
func (c *Client) Login(ctx context.Context, username, password string) error {
	if err := validate.Username(username); err != nil {
		return err
	}

	if err := validate.Password(password); err != nil {
		return err
	}

	var header metadata.MD

	req := authpb.LoginRequest_builder{
		Username: proto.String(username),
		Password: proto.String(password),
	}.Build()

	resp, err := c.auth.Login(ctx, req, grpc.Header(&header))
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

	key, err := crypto.DeriveKey(password, resp.GetSalt())
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	if err := session.Save(key, session.DefaultTTL); err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	fmt.Printf("Logged in as: %s\n", resp.GetUser().GetUsername())
	return nil
}

// Logout удаляет ключ из keychain.
func (c *Client) Logout() error {
	return session.Delete()
}
