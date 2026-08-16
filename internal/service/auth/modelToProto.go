package auth

import (
	pb "credentials-vault/gen/go/auth/v1"
	"credentials-vault/internal/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func modelToProto(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
