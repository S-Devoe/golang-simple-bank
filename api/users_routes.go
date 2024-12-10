package api

import "github.com/gin-gonic/gin"

func (server *Server) setUpUserRoutes(router *gin.RouterGroup) {
	usersGroup := router.Group("/users")

	{ //create user routes doesnt need token
		usersGroup.POST("", server.createUser) // create user

		authUsersGroup := usersGroup.Group("").Use(authMiddleware(server.tokenMaker))

		authUsersGroup.GET("/:username", server.getUser)       // get user by username
		authUsersGroup.DELETE("/:username", server.deleteUser) // delete user by username
	}
}
