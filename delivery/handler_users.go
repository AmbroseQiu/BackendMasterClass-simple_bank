package delivery

import (
	"context"
	"net/http"

	"github.com/backendmaster/simple_bank/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type usersHandlerDelivery struct {
	usercase domain.UsersTableUseCase
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
		if errors.Is(err, domain.ErrorUniqueViolation) {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
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
		if errors.Cause(err) == domain.ErrorUserNotFound {
			ctx.JSON(http.StatusNotFound, errResponse(err))
		} else if errors.Cause(err) == domain.ErrorInternalServerError {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
		} else if errors.Cause(err) == domain.ErrorStatusForbidden {
			ctx.JSON(http.StatusForbidden, errResponse(err))
		} else if errors.Cause(err) == domain.ErrorPermissionNowAllowed {
			ctx.JSON(http.StatusForbidden, errResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
