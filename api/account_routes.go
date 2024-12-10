package api

import "github.com/gin-gonic/gin"

func (server *Server) setUpAccountRoutes(router *gin.RouterGroup) {
	accountsGroup := router.Group("/accounts").Use(authMiddleware(server.tokenMaker))
	{
		// accounts endpoint
		accountsGroup.POST("", server.createAccount)
		accountsGroup.GET("/:id", server.getAccount)
		accountsGroup.GET("", server.listAccounts)
	}
}
