package delivery

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/backendmaster/simple_bank/domain"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type usersHandlerDelivery struct {
	usercase domain.UsersTableUseCase
	router   *gin.Engine
}

func NewUsersHandlerDelivery(usercase domain.UsersTableUseCase) usersHandlerDelivery {
	return usersHandlerDelivery{
		usercase: usercase,
	}
}

func errResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}

func (d *usersHandlerDelivery) handlerCreateUser(ctx *gin.Context) {
	var req domain.CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	rsp, err := d.usercase.CreateUser(context.Background(), req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (d *usersHandlerDelivery) handlerLoginUser(ctx *gin.Context) {
	// bind req body
	var req domain.LoginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	rsp, err := d.usercase.LoginUser(context.Background(), req)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (d *usersHandlerDelivery) SetupRouter() {

	router := gin.Default()
	router.POST("/users", d.handlerCreateUser)
	router.POST("/users/login", d.handlerLoginUser)
	// d.router.POST("tokens/renew_access", server.renewAccessToken)

	// authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// authRoute.POST("/accounts", server.createAccount)
	// authRoute.GET("/accounts/:id", server.getAccount)
	// authRoute.GET("/accounts", server.listAccount)
	// authRoute.POST("/transfers", server.createTransfer)
	d.router = router
}

func (d *usersHandlerDelivery) Start(address string) error {
	return d.router.Run(address)
}
