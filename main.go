package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/pakojabi/simplebank/api"
	db "github.com/pakojabi/simplebank/db/sqlc"
	"github.com/pakojabi/simplebank/gapi"
	"github.com/pakojabi/simplebank/pb"
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

	// err = runGinServer(config, store)
	err = runGrpcServer(config, store)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}

func runGinServer(config util.Config, store db.Store) error {
	server, err := api.NewServer(config, store)
	if err != nil {
		return err
	}

	return server.Start(config.HTTPServerAddress)
}

func runGrpcServer(config util.Config, store db.Store) error {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)

	// optional but allows the client discover the service
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		return err
	}

	log.Printf("starting grpc server at %s", listener.Addr().String())
	return grpcServer.Serve(listener)
}
