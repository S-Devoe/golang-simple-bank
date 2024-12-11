package api

import (
	"net/http"

	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/S-Devoe/golang-simple-bank/util/password"
	"github.com/gin-gonic/gin"
)

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	Username    string       `json:"username"`
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, util.CreateResponse(http.StatusNotFound, nil, "User not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	passwordMatch, err := password.ComparePasswordAndHash(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, err))
		return
	}

	if !passwordMatch {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, "invalid username or password"))
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(user.Username, user.Email, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	rsp := LoginUserResponse{
		Username:    user.Username,
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, util.CreateResponse(http.StatusOK, rsp, nil))
}
