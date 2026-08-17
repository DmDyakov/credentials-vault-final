package user

import (
	"strings"

	"credentials-vault/internal/domain"
)

func validateRegisterCredentials(username, password string) error {
	username = strings.TrimSpace(username)

	if username == "" {
		return domain.ErrUsernameRequired
	}
	if len(username) < 3 {
		return domain.ErrUsernameTooShort
	}
	if len(username) > 64 {
		return domain.ErrUsernameTooLong
	}
	if password == "" {
		return domain.ErrPasswordRequired
	}
	if len(password) < 6 {
		return domain.ErrPasswordTooShort
	}
	if len(password) > 72 {
		return domain.ErrPasswordTooLong
	}

	return nil
}

func validateLoginCredentials(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return domain.ErrUsernameRequired
	}
	if password == "" {
		return domain.ErrPasswordRequired
	}

	return nil
}
