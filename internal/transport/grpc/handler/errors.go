package handler

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"credentials-vault/internal/domain"
)

// errorToCode - маппинг доменных ошибок в gRPC коды.
var errorToCode = map[error]codes.Code{
	// Auth errors
	domain.ErrUserAlreadyExists:  codes.AlreadyExists,
	domain.ErrInvalidCredentials: codes.Unauthenticated,
	domain.ErrUsernameRequired:   codes.InvalidArgument,
	domain.ErrPasswordRequired:   codes.InvalidArgument,
	domain.ErrUsernameTooShort:   codes.InvalidArgument,
	domain.ErrUsernameTooLong:    codes.InvalidArgument,
	domain.ErrPasswordTooShort:   codes.InvalidArgument,
	domain.ErrPasswordTooLong:    codes.InvalidArgument,

	// Vault errors
	domain.ErrVaultItemNotFound:     codes.NotFound,
	domain.ErrVaultItemIDRequired:   codes.InvalidArgument,
	domain.ErrInvalidItemType:       codes.InvalidArgument,
	domain.ErrEncryptedDataRequired: codes.InvalidArgument,
	domain.ErrUserIDRequired:        codes.Unauthenticated,
}

// mapError маппит доменные ошибки в gRPC статусы.
func mapError(err error) error {
	for domainErr, code := range errorToCode {
		if errors.Is(err, domainErr) {
			return status.Error(code, domainErr.Error())
		}
	}
	return status.Error(codes.Internal, "internal error")
}
