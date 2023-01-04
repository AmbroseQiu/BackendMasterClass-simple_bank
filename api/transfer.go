package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	From_Accoount_ID int64  `json:"from_account_id" binding:"required,min=1"`
	To_Account_ID    int64  `json:"to_account_id" binding:"required,min=1"`
	Amount           int64  `json:"amount" binding:"required,gt=0"`
	Currency         string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.From_Accoount_ID,
		ToAccountID:   req.To_Account_ID,
		Amount:        req.Amount,
	}

	if !server.validateAccount(ctx, req.From_Accoount_ID, req.Currency) {
		return
	}
	if !server.validateAccount(ctx, req.To_Account_ID, req.Currency) {
		return
	}

	account, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("accouont %v mismatched: %v vs %v", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return false
	}

	return true
}
