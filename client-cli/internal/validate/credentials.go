// Package validate содержит функции валидации данных.
package validate

import (
	"fmt"
	"strings"
)

// Credentials проверяет учётные данные.
func Credentials(site, username, password string) error {
	if site == "" {
		return fmt.Errorf("site is required")
	}

	if !strings.Contains(site, ".") {
		return fmt.Errorf("site must be a valid domain")
	}

	if username == "" {
		return fmt.Errorf("username is required")
	}

	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	return nil
}
