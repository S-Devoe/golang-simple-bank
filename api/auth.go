package api

import (
	"fmt"
	"net/http"
	"time"

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
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
	SessionID             string       `json:"session_id"`
}

// Login godoc
// @Summary Login
// @Description Login a user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginUserRequest true "Login User Request"
// @Success 200 {object} util.Response{data=LoginUserResponse} "Success"
// @Failure 400 {object} util.Response "Bad Request"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 404 {object} util.Response "Not Found"
// @Failure 500 {object} util.Response "Internal Server Error"
// @Router /login [post]
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

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, user.Email, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		user.Email,
		s.config.RefreshTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID.String(),
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		fmt.Println(err)
		return
	}

	rsp := LoginUserResponse{
		User:                  newUserResponse(user),
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, util.CreateResponse(http.StatusOK, rsp, nil))
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// RenewAccessToken godoc
// @Summary Renew Access Token
// @Description Renew a user's access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param renewRequest body renewAccessTokenRequest true "Request body to renew access/refresh tokens"
// @Success 200 {object} util.Response{data=renewAccessTokenResponse} "Success"
// @Failure 400 {object} util.Response "Bad Request"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 404 {object} util.Response "Session Not Found"
// @Failure 500 {object} util.Response "Internal Server Error"
// @Router /token/renew [post]
func (s *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err))
		return
	}
	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, "Expired or invalid refresh token, please login again"))
		return
	}

	session, err := s.store.GetSession(ctx, string(refreshPayload.ID.String()))
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, util.CreateResponse(http.StatusNotFound, nil, err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, "Expired or invalid refresh token, please login again"))
		return
	}

	if session.Username != refreshPayload.Username {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, "Incorrect session user"))
		return
	}
	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, "mismathced session token"))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, util.CreateResponse(http.StatusUnauthorized, nil, "session expired"))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.Username,
		refreshPayload.Email,
		s.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, util.CreateResponse(http.StatusOK, rsp, nil))
}
