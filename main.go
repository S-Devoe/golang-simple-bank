package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/S-Devoe/golang-simple-bank/api"
	"github.com/S-Devoe/golang-simple-bank/config"
	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	_ "github.com/S-Devoe/golang-simple-bank/docs"
	"github.com/S-Devoe/golang-simple-bank/gapi"
	"github.com/S-Devoe/golang-simple-bank/pb"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

// @BasePath /api/v1
// @version 1.0
// @title Simple Bank API
// @description Simple Bank API documentation.
// @host localhost:8080
func main() {
	config := config.InitConfig()
	fmt.Println(config.RefreshTokenDuration)
	connection, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(connection)
	runGinServer(config, store)
	// runGrpcServer(config, store)

}

func runGrpcServer(config config.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatalf("cannot listen on %s: %v", config.HttpServerAddress, err)
	}
	log.Println("Starting gRPC server on", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server: ", err)
	}
}

func runGinServer(config config.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
