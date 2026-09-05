package validate

import (
	"fmt"
	"strings"
)

// Username проверяет имя пользователя.
func Username(username string) error {
	username = strings.TrimSpace(username)

	if username == "" {
		return fmt.Errorf("username is required")
	}

	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}

	if len(username) > 64 {
		return fmt.Errorf("username must be at most 64 characters")
	}

	return nil
}

// Password проверяет пароль.
func Password(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	if len(password) > 72 {
		return fmt.Errorf("password must be at most 72 characters")
	}

	return nil
}
