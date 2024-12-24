package gapi

import (
	"context"

	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/pb"
	"github.com/S-Devoe/golang-simple-bank/util/password"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == db.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found")

		}
		return nil, status.Errorf(codes.Internal, "cannot fetch user: %v", err)
	}

	passwordMatch, err := password.ComparePasswordAndHash(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	if !passwordMatch {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, user.Email, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		user.Email,
		s.config.RefreshTokenDuration,
	)

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID.String(),
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	resp := &pb.LoginResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		SessionId:             session.ID,
		User:                  converteUser(user),
	}
	return resp, nil

}
