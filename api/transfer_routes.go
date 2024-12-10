package api

import "github.com/gin-gonic/gin"

func (server *Server) setUpTransferRoutes(router *gin.RouterGroup) {
	transferGroup := router.Group("/transfer").Use(authMiddleware(server.tokenMaker))
	{
		// transfers endpoints
		transferGroup.POST("/transfers", server.createTransfer)
	}
}
