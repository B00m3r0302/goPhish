package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/B00m3r0302/goPhish/backend/services"
	"github.com/B00m3r0302/goPhish/server/config"
)

// CreateCommand creates a new command
// @Summary Create a new command
// @Description Creates a new command and returns its ID
// @Tags Commands
// @Accept json
// @Produce	json
// @Param command body struct{ AgentID uint; Command string } true "Command details"
// @Success 201 {object} map[string]unit "id"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /commands/create [post]
func CreateCommand(c *gin.Context) {
	var commandRequest struct {
		AgentID uint	`json:"agent_id" binding:"required"`
		Command string	`json:"command" binding:"required"`
	}

	if err := c.ShouldBindJSON(&commandRequest); err != nil {
		logger.WithFields(logrus.Fields{
			"client_ip":	c.ClientIP(),
			"error":		err.Error(),
		}).Warn("Invalid input for command creation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	commandService := services.NewCommandService(config.DB)
	commandID, err := commandService.CreateCommand(commandRequest.AgentID, commandRequest.Command)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"agent_id":		commandRequest.AgentID,
			"command": 		commandRequest.Command,
			"client_ip": 	c.ClientIP,
			"error": 		err.Error(),
		}).Error("Failed to create command")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create command"})
		return 
	}

	logger.WithFields(logrus.Fields{
		"command_id": commandID,
		"agent_id":	  commandRequest.AgentID,
		"client_ip":  c.ClientIP(),
	}).Info("Command created successfully")
	c.JSON(http.StatusCreated, gin.H{"id": commandID})
}

// GetCommand retrieves a command by ID
// @Summary Retrieve a command
// @Description Retrieves a command by ID
// @Tags Commands
// @Produce json
// @Param id path int true "Command ID"
// @Success 200 {object} services.Command
// @Failure 404 {object} map[string]string "error"
// @Router /commands/{id} [get]
func GetCommand(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	commandService := services.NewCommandService(config.DB)
	command, err := commandService.GetCommandByID(uint(id))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"command_id":	id,
			"client_ip":	c.ClientIP(),
		}).Warn("Command not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Command not found"})
		return 
	}

	logger.WithFields(logrus.Fields{
		"command_id": command.ID,
		"agent_id":	  command.AgentID,
		"client_ip":  c.ClientIP(),	
	}).Info("Command retrieved successfully")\
	c.JSON(http.StatusOK, command)
}

// AppendCommandOutput appends output to an existing command
// @Summary Append output to a command
// @Description Appends output to an existing command
// @Tags Commands
// @Accept json
// @Param id path int true "Command ID"
// @Param output body struct{ Output string } true "Command output"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /commands/{id}/output [put]
func AppendCommandOutput(c *gin.Context) {
	id, _ := strconf.Atoi(c.Param("id"))
	var outputRequest struct {
		Output string `json:"output" binding:"required` 
	}

	if err := c.ShouldBindJSON(&outputRequest); err != nil {
		logger.WithFields(logrus.Fields{
			"command_id":	id,
			"client_ip":	c.ClientIP(),
			"error":		err.Error(),
		}).Warn("Invalid input for appending command output")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return 
	}

	commandService := services.NewCommandService(config.DB)
	if err := commandService.AppendCommandOutput(uint(id), outputRequest.Output); err != nil {
		logger.WithFields(logrus.Fields{
			"command_id":	id,
			"output":		outputRequest.Output,
			"client_ip":	c.ClientIP(),
			"error":		err.Error(),
		}).Error("Failed to append output to command")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to append output tp command"})
		return
	}

	logger.WithFields(logrus.Fields{
		"command_id":	id,
		"output": outputRequest.Output,
		"client_ip": c.ClientIP(),
	}).Info("Output appended to command successfully")
	c.Status(http.StatusNoContent)
}