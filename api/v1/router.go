package server

import (
	"your_project/api/v1/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router and registers all the routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// API Version 1 group
	v1 := router.Group("/api/v1")
	{
		// Agent routes
		v1.POST("/agents", controllers.CreateAgent)
		v1.GET("/agents/:id", controllers.GetAgent)
		v1.GET("/agents", controllers.ListAgents)
		v1.DELETE("/agents/:id", controllers.DeleteAgent)

		// Command routes
		v1.POST("/commands", controllers.CreateCommand)
		v1.GET("/commands/:id", controllers.GetCommand)
		v1.GET("/commands", controllers.ListCommands)
		v1.DELETE("/commands/:id", controllers.DeleteCommand)

		// Listener routes
		v1.POST("/listeners", controllers.CreateListener)
		v1.GET("/listeners/:id", controllers.GetListener)
		v1.GET("/listeners", controllers.ListListeners)
		v1.DELETE("/listeners/:id", controllers.DeleteListener)

		// User routes
		v1.POST("/users", controllers.CreateUser)
		v1.POST("/users/authenticate", controllers.AuthenticateUser)
		v1.GET("/users/:id", controllers.GetUser)
	}

	return router
}
