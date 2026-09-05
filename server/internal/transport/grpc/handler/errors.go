package handler

import (
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"credentials-vault/server/internal/domain"
)

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
	domain.ErrSaltRequired:       codes.InvalidArgument,

	// Vault errors
	domain.ErrVaultItemNotFound:     codes.NotFound,
	domain.ErrVaultItemIDRequired:   codes.InvalidArgument,
	domain.ErrInvalidItemType:       codes.InvalidArgument,
	domain.ErrEncryptedDataRequired: codes.InvalidArgument,
	domain.ErrUserIDRequired:        codes.Unauthenticated,
}

func mapError(err error, logger *zap.Logger) error {
	for domainErr, code := range errorToCode {
		if errors.Is(err, domainErr) {
			return status.Error(code, domainErr.Error())
		}
	}
	logger.Error("internal error", zap.Error(err))

	return status.Error(codes.Internal, "internal error")
}
