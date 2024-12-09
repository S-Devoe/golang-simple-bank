package api

import "github.com/gin-gonic/gin"

func (server *Server) setUpTransferRoutes(router *gin.RouterGroup) {
	transferGroup := router.Group("/transfer")
	{
		// transfers endpoints
		transferGroup.POST("/transfers", server.createTransfer)
	}
}
