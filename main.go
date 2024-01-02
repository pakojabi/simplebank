package main

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"

	"github.com/pakojabi/simplebank/api"
	db "github.com/pakojabi/simplebank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:postgres123@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	var err error
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	
	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}