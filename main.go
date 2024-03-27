package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jpnsantoss/simplebank/api"
	db "github.com/jpnsantoss/simplebank/db/sqlc"
)

const (
	dbSource     = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAdress = "0.0.0.0:8080"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAdress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
