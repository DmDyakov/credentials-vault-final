package auth

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateRegisterCredentials(username, password string) error {
	username = strings.TrimSpace(username)

	if username == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}
	if len(username) < 3 {
		return status.Error(codes.InvalidArgument, "username must be at least 3 characters")
	}
	if len(username) > 64 {
		return status.Error(codes.InvalidArgument, "username must be at most 64 characters")
	}
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if len(password) < 6 {
		return status.Error(codes.InvalidArgument, "password must be at least 6 characters")
	}
	if len(password) > 72 {
		return status.Error(codes.InvalidArgument, "password must be at most 72 characters")
	}

	return nil
}

func validateLoginCredentials(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
