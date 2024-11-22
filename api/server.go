package api

import (
	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// server struct will serve all requests for the banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// Newserver creates a new http server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server

}

// Start runs the http server on a specific address
func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

// errorResponse will return an error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
