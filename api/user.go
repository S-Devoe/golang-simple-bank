package api

import (
	"net/http"
	"time"

	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/S-Devoe/golang-simple-bank/util/password"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	PasswordChanged time.Time `json:"password_changed_at,omitempty"` // Omit if empty

}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:        user.Username,
		FullName:        user.FullName,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt,
		PasswordChanged: user.PasswordChangedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err))
		return
	}

	hashedPassword, err := password.GeneratePasswordHash(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, util.CreateResponse(http.StatusForbidden, nil, "Username or Email already exists"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return

	}
	resp := newUserResponse(user)
	ctx.JSON(http.StatusCreated, util.CreateResponse(http.StatusCreated, resp, nil))
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required"`
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
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

	ctx.JSON(http.StatusOK, util.CreateResponse(http.StatusOK, newUserResponse(user), nil))
}

type messageResponse struct {
	Message string `json:"message"`
}

func (s *Server) deleteUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.CreateResponse(http.StatusBadRequest, nil, err))
		return
	}

	// First, check if the user exists
	_, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, util.CreateResponse(http.StatusNotFound, nil, "User not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	// Now, proceed to delete the user
	err = s.store.DeleteUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.CreateResponse(http.StatusInternalServerError, nil, err))
		return
	}

	// Successfully deleted
	resp := messageResponse{
		Message: "User deleted",
	}
	ctx.JSON(http.StatusOK, util.CreateResponse(http.StatusOK, resp, nil))
}
