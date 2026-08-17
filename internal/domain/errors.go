// Package domain содержит доменные модели и ошибки.
package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUsernameRequired   = errors.New("username is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrUsernameTooShort   = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong    = errors.New("username must be at most 64 characters")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
	ErrPasswordTooLong    = errors.New("password must be at most 72 characters")
)

var (
	ErrVaultItemNotFound = errors.New("vault item not found")
)
