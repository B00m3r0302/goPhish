package controllers

import (
	"net/http"
	"strconv"

	"your_project/backend/models"
	"your_project/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// StartListener starts a new listener
// @Summary Start a new listener
// @Description Starts a new listener with the specified parameters
// @Tags Listeners
// @Accept  json
// @Produce  json
// @Param listener body models.ListenerRequest true "Listener details"
// @Success 201 {object} models.Listener "Listener started successfully"
// @Failure 400 {object} gin.H "Invalid input"
// @Failure 500 {object} gin.H "Failed to start listener"
// @Router /listeners/start [post]
func StartListener(c *gin.Context) {
	var listenerRequest models.ListenerRequest
	if err := c.ShouldBindJSON(&listenerRequest); err != nil {
		logger.WithError(err).Error("Invalid input for starting listener")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	listener, err := services.StartListener(listenerRequest)
	if err != nil {
		logger.WithError(err).Error("Failed to start listener")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start listener"})
		return
	}

	logger.WithFields(logrus.Fields{
		"listener_id": listener.ID,
		"client_ip":   c.ClientIP(),
	}).Info("Listener started successfully")
	c.JSON(http.StatusCreated, listener)
}

// StopListener stops an existing listener
// @Summary Stop an existing listener
// @Description Stops a listener by its ID
// @Tags Listeners
// @Produce  json
// @Param id path int true "Listener ID"
// @Success 204 "Listener stopped successfully"
// @Failure 404 {object} gin.H "Listener not found"
// @Failure 500 {object} gin.H "Failed to stop listener"
// @Router /listeners/stop/{id} [post]
func StopListener(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("Invalid listener ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listener ID"})
		return
	}

	err = services.StopListener(id)
	if err != nil {
		logger.WithError(err).Error("Failed to stop listener")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop listener"})
		return
	}

	logger.WithFields(logrus.Fields{
		"listener_id": id,
		"client_ip":   c.ClientIP(),
	}).Info("Listener stopped successfully")
	c.Status(http.StatusNoContent)
}

// GetListener retrieves an existing listener
// @Summary Retrieve an existing listener
// @Description Retrieves a listener by its ID
// @Tags Listeners
// @Produce  json
// @Param id path int true "Listener ID"
// @Success 200 {object} models.Listener "Listener retrieved successfully"
// @Failure 404 {object} gin.H "Listener not found"
// @Router /listeners/{id} [get]
func GetListener(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("Invalid listener ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listener ID"})
		return
	}

	listener, err := services.GetListener(id)
	if err != nil {
		logger.WithError(err).Error("Listener not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Listener not found"})
		return
	}

	logger.WithFields(logrus.Fields{
		"listener_id": listener.ID,
		"client_ip":   c.ClientIP(),
	}).Info("Listener retrieved successfully")
	c.JSON(http.StatusOK, listener)
}
