package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/pakojabi/simplebank/api"
	db "github.com/pakojabi/simplebank/db/sqlc"
	"github.com/pakojabi/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	log.Default().Printf("Connected to %s using driver %s", config.DBSource, config.DBDriver)
	
	store := db.NewStore(conn)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}	