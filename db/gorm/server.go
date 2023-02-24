package gorm

import (
	"fmt"

	"github.com/backendmaster/simple_bank/pb"
	"github.com/backendmaster/simple_bank/token"
	"github.com/backendmaster/simple_bank/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HttpServer struct {
	config     util.Config
	store      *gorm.DB
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewHttpServer(config util.Config, store *gorm.DB) (*HttpServer, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't not create token")
	}
	server := &HttpServer{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker}

	server.setupRouter()

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	return server, nil
}

func (server *HttpServer) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	// router.POST("tokens/renew_access", server.renewAccessToken)

	// authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// authRoute.POST("/accounts", server.createAccount)
	// authRoute.GET("/accounts/:id", server.getAccount)
	// authRoute.GET("/accounts", server.listAccount)
	// authRoute.POST("/transfers", server.createTransfer)
	server.router = router
}

func (server *HttpServer) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}

type GrpcServer struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      *gorm.DB
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewGrpcServer(config util.Config, store *gorm.DB) (*GrpcServer, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't not create token")
	}
	server := &GrpcServer{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker}

	return server, nil
}
