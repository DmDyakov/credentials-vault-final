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
	ErrSaltRequired       = errors.New("salt is required")
)

var (
	ErrVaultItemNotFound     = errors.New("vault item not found")
	ErrInvalidItemType       = errors.New("invalid item type")
	ErrEncryptedDataRequired = errors.New("encrypted data is required")
	ErrUserIDRequired        = errors.New("user id is required")
	ErrVaultItemIDRequired   = errors.New("vault item id is required")
)
