package api

import "github.com/gin-gonic/gin"

func (server *Server) setUpUserRoutes(router *gin.RouterGroup) {
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("", server.createUser)             // create user
		usersGroup.GET("/:username", server.getUser)       // get user by username
		usersGroup.DELETE("/:username", server.deleteUser) // delete user by username
	}
}
