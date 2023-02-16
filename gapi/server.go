package gapi

import (
	"fmt"

	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/backendmaster/simple_bank/pb"
	"github.com/backendmaster/simple_bank/token"
	"github.com/backendmaster/simple_bank/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't not create token")
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker}

	return server, nil
}
