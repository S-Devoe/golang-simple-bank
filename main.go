package main

import (
	"context"
	"fmt"
	"log"

	"github.com/S-Devoe/golang-simple-bank/api"
	"github.com/S-Devoe/golang-simple-bank/config"
	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	config := config.InitConfig()
	fmt.Println(config.RefreshTokenDuration)
	connection, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(connection)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
