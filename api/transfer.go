package api

import (
	"fmt"
	"net/http"

	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/token"
	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64   `json:"to_account_id" binding:"required,min=1"`
	Currency      string  `json:"currency" binding:"required,currency"`
	Amount        float64 `json:"amount" binding:"required,min=1"`
}

type transferSuccessResponse struct {
	Message string `json:"message"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err))
		return
	}
	fromAccount, valid := server.validAccount(ctx, req.FromAccountId, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, "Account doesn not belong to this authenticated user"))
		return
	}
	_, valid = server.validAccount(ctx, req.ToAccountId, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	if result.FromAccount.Balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, "Insufficient balance"))
		return
	}

	res := &transferSuccessResponse{
		Message: fmt.Sprintf("Transfer of %.2f %s successful.", result.Transfer.Amount, req.Currency),
	}

	ctx.JSON(http.StatusCreated, util.CreateResponse(http.StatusCreated, res, nil))
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, util.CreateResponse(http.StatusNotFound, nil, "Account not found"))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return account, false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, "Invalid currency for the account"))
		return account, false
	}
	return account, true
}
