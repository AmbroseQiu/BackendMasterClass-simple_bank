package delivery

import (
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/backendmaster/simple_bank/db/gorm"
	repository "github.com/backendmaster/simple_bank/repository/postgresql"
	"github.com/backendmaster/simple_bank/token"
	"github.com/backendmaster/simple_bank/usercase"
	"github.com/backendmaster/simple_bank/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	userhandler usersHandlerDelivery
	router      *gin.Engine
}

func NewGormServer(config util.Config, conn *sql.DB) (*Server, error) {
	gormDb, err := gorm.NewDB(conn)
	if err != nil {
		log.Fatal().Err(err)
	}
	userRepo := repository.NewpostgresqlUserRepository(gormDb)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't not create token ")
	}
	sessionRepo := repository.NewpostgresqlSessionRepository(gormDb)
	sessionUseCase := usercase.NewSessionUseCase(sessionRepo)
	usersTableUseCase := usercase.NewusersTableUserCase(config, userRepo, sessionUseCase, tokenMaker)
	userhandler := NewUsersHandlerDelivery(usersTableUseCase)
	server := &Server{
		userhandler: userhandler,
	}
	server.SetupRouter()
	return server, nil
}

func (s *Server) SetupRouter() {

	router := gin.Default()
	router.POST("/users", s.userhandler.handlerCreateUser)
	router.POST("/users/login", s.userhandler.handlerLoginUser)
	// d.router.POST("tokens/renew_access", server.renewAccessToken)

	// authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// authRoute.POST("/accounts", server.createAccount)
	// authRoute.GET("/accounts/:id", server.getAccount)
	// authRoute.GET("/accounts", server.listAccount)
	// authRoute.POST("/transfers", server.createTransfer)
	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
