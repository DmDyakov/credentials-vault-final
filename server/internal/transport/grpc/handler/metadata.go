package handler

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// getUserIDFromMetadata извлекает user_id из gRPC метаданных.
func getUserIDFromMetadata(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("user_id")
	if len(values) == 0 {
		return uuid.Nil, status.Error(codes.Unauthenticated, "user_id is not provided")
	}

	userID, err := uuid.Parse(values[0])
	if err != nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid user_id")
	}

	return userID, nil
}
