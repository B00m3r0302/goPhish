package controllers

import (
	"net/http"
	"strconv"

	"your_project/backend/models"
	"your_project/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var userLogger = logrus.New()

// CreateUser creates a new user
// @Summary Create a new user
// @Description Creates a new user with the specified details
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body models.UserRequest true "User details"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} gin.H "Invalid input"
// @Failure 500 {object} gin.H "Failed to create user"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var userRequest models.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		userLogger.WithError(err).Error("Invalid input for creating user")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := services.CreateUser(userRequest)
	if err != nil {
		userLogger.WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userLogger.WithFields(logrus.Fields{
		"user_id":   user.ID,
		"username":  user.Username,
		"client_ip": c.ClientIP(),
	}).Info("User created successfully")
	c.JSON(http.StatusCreated, user)
}

// AuthenticateUser authenticates a user
// @Summary Authenticate a user
// @Description Authenticates a user with the specified credentials
// @Tags Users
// @Accept  json
// @Produce  json
// @Param credentials body models.AuthRequest true "User credentials"
// @Success 200 {object} models.AuthResponse "Authentication successful"
// @Failure 401 {object} gin.H "Authentication failed"
// @Router /users/authenticate [post]
func AuthenticateUser(c *gin.Context) {
	var authRequest models.AuthRequest
	if err := c.ShouldBindJSON(&authRequest); err != nil {
		userLogger.WithError(err).Error("Invalid input for authentication")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	authResponse, err := services.AuthenticateUser(authRequest)
	if err != nil {
		userLogger.WithError(err).Error("Authentication failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userLogger.WithFields(logrus.Fields{
		"user_id":   authResponse.UserID,
		"username":  authRequest.Username,
		"client_ip": c.ClientIP(),
	}).Info("User authenticated successfully")
	c.JSON(http.StatusOK, authResponse)
}

// GetUser retrieves an existing user by ID
// @Summary Retrieve an existing user
// @Description Retrieves a user by their ID
// @Tags Users
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User "User retrieved successfully"
// @Failure 404 {object} gin.H "User not found"
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		userLogger.WithError(err).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := services.GetUser(id)
	if err != nil {
		userLogger.WithError(err).Error("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userLogger.WithFields(logrus.Fields{
		"user_id":   user.ID,
		"username":  user.Username,
		"client_ip": c.ClientIP(),
	}).Info("User retrieved successfully")
	c.JSON(http.StatusOK, user)
}
