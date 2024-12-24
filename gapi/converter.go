package gapi

import (
	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converteUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         timestamppb.New(user.CreatedAt),
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
	}
}
