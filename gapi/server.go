package gapi

import (
	"fmt"

	db "github.com/pakojabi/simplebank/db/sqlc"
	"github.com/pakojabi/simplebank/pb"
	"github.com/pakojabi/simplebank/token"
	"github.com/pakojabi/simplebank/util"
)

// Server serves gRPC requests
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server instance
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
