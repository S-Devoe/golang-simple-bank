package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/token"
	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// create account
type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// ctx.JSON(http.StatusBadRequest, errorResponse(err))
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, util.CreateResponse(http.StatusForbidden, nil, err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}
	ctx.JSON(http.StatusCreated, util.CreateResponse(http.StatusCreated, account, nil))
}

// get account by id
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.CreateResponse(http.StatusNotFound, nil, "Account not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, "Account doesn't belong to this authenticated user"))
		return
	}
	ctx.JSON(http.StatusOK, util.CreateResponse(http.StatusOK, account, nil))
}

// get accounts
type listAccountsRequest struct {
	Limit int32 `query:"limit"`
	Page  int32 `query:"page"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest

	pagination, err := util.ParsePaginationQuery(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err.Error()))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
		Owner:  authPayload.Username,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	totalItems, err := server.store.CountAccounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}
	ctx.JSON(http.StatusOK, util.CreatePaginatedResponse(http.StatusOK, accounts, pagination.Page, pagination.Limit, totalItems, nil))
}
