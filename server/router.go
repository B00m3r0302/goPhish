package server

import (
	"your_project/api/v1/controllers"
	"your_project/server/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRouter() *gin.Engine {
	// Initialize the logger
	logger := logrus.New()

	// Initialize the router
	router := gin.Default()

	// Apply middlewares
	router.Use(middlewares.LoggingMiddleware(logger))
	router.Use(middlewares.ErrorHandlingMiddleware(logger))
	router.Use(middlewares.AuthMiddleware())

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/agents", controllers.CreateAgent)
		v1.GET("/agents/:id", controllers.GetAgent)
		v1.GET("/agents", controllers.ListAgents)
		v1.DELETE("/agents/:id", controllers.DeleteAgent)

		v1.POST("/commands", controllers.CreateCommand)
		v1.GET("/commands/:id", controllers.GetCommand)
		v1.GET("/commands", controllers.ListCommands)
		v1.DELETE("/commands/:id", controllers.DeleteCommand)

		v1.POST("/listeners", controllers.CreateListener)
		v1.GET("/listeners/:id", controllers.GetListener)
		v1.GET("/listeners", controllers.ListListeners)
		v1.DELETE("/listeners/:id", controllers.DeleteListener)

		v1.POST("/users", controllers.CreateUser)
		v1.POST("/users/authenticate", controllers.AuthenticateUser)
		v1.GET("/users/:id", controllers.GetUser)
	}

	return router
}
