package gapi

import (
	"context"

	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/pb"
	"github.com/S-Devoe/golang-simple-bank/util/password"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := password.GeneratePasswordHash(req.GetPassword())
	if err != nil {
		// return util.CreateResponse(900, nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)),note: I need to create another util response for GRPC
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)
	}
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
		HashedPassword: hashedPassword,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {

			return nil, status.Errorf(codes.AlreadyExists, "username or email already exists")
		}

		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}
	resp := &pb.CreateUserResponse{
		User: converteUser(user),
	}
	return resp, nil

}
