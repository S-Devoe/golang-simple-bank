package gapi

import (
	"fmt"

	"github.com/S-Devoe/golang-simple-bank/config"
	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/pb"
	"github.com/S-Devoe/golang-simple-bank/token"
)

// server struct will serve all gRPC requests for the banking service
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	config     config.Config
	pb.UnimplementedSimpleBankServer
}

// Newserver creates a new gRPC server.
func NewServer(config config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil

}
