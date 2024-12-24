package api

import (
	"fmt"

	"github.com/S-Devoe/golang-simple-bank/config"
	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	_ "github.com/S-Devoe/golang-simple-bank/docs"
	"github.com/S-Devoe/golang-simple-bank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// server struct will serve all requests for the banking service
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     config.Config
}

// Newserver creates a new http server and setup routing
func NewServer(config config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setUpRouter()

	return server, nil

}

// Start runs the http server on a specific address
func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	//add swagger
	router.GET("/docs/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/docs/" {
			c.Redirect(302, "/docs/index.html")
			return
		}
	},
		ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		server.setUpUserRoutes(api)
		server.setUpAccountRoutes(api)
		server.setUpAuthRoutes(api)
		server.setUpTransferRoutes(api)

	}

	server.router = router
}

// errorResponse will return an error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
