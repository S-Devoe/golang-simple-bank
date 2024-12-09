package api

import "github.com/gin-gonic/gin"

func (server *Server) setUpAuthRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("")
	{
		// auth
		authGroup.POST("/login", server.loginUser)
		authGroup.POST("/signup", server.createUser)
	}
}
