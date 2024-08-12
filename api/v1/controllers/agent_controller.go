package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/B00m3r0302/goPhish/backend/services"
	"github.com/B00m3r0302/goPhish/server/config"
)

var logger = logrus.new()

// RegisterAgent registers a new agent
// @Summary Register a new agent
// @Description Registers a new agent and returns its ID
// @Tags Agents
// @Accept json
// @Produce json
// @Param agent body struct{ Name sting; Type string } true "Agent Details"
// @Success 201 {object} map[string]uint "id"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @ router /agents/register [post]
func RegisterAgent(c *gin.Context) {
	var agentRequest struct {
		Name string 'json:"name" binding:"required"'
		Type string 'json:"type" binding:"required"'
	}
	
	if err := c.ShouldBindJSON(&agentRequest); err != nil {
		logger.WithFields(logrus.Fields{
			"client_ip":	c.ClientIP(),
			"error":		err.Error(),
		}).Warn("Invalid input for agent registration")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	agentService := services.NewAgentService(config.DB)
	agentID, err := agentService.RegisterAgent(agentRequest.Name, agentRequest.Type)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"agent_name":	agentRequest.Name,
			"agent_type":	agentRequest.Type,
			"client_ip":	c.ClientIP(),
			"error":		err.Error(),
		}).Error("Failed to register agent")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register agent"})
		return
	}

	logger.WithFields(logrus.Fields{
		"agent_id":		agentID,
		"agent_name":	agentRequest.Name,
		"client_ip": 	c.ClientIP(),
	}).Info("Agent registered successfully")
	c.JSON(http.StatusCreated, gin.H{"id": agentID})
}

// CheckInAgent allows an agent to check in 
// @Summary Agent check-in
// @Description Allows an agent to check in by updating it's status
// @Tags Agents
// @Accept json
// @Produce json
// @Param id path int true "Agent ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /agents/checkin/{id} [post]
func CheckInAgent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := agentService.CheckInAgent(uint(id)); err != nil {
		logger.WithFields(logrus.Fields{
			"agent_id":		id,
			"client_ip":	c.ClientIP(),
			"error":		err.Error(),
		}).Warn("Agent not found or failed to check in")
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found or failed to check in"})
		return
	}

	logger.WithFields(logrus.Fields{
		"agent_id":		id,
		"client_ip":	c.ClientIP(),
	}).Info("Agent checked in successfully")
	c.Status(http.StatusNoContent)
}

// GetAgent retrieves an agent by ID
// @Summary Retrieve an agent
// @Description Retrieves an agent by ID
// @Tags Agents
// @Produce json
// @Param id path int true "Agent ID"
// @Success 200 {object} services.Agent
// @Failure 404 {object} map[string]string "error"
// @Router /agents/{id} [get]
func GetAgent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	agentService := services.NewAgentService(config.DB)
	agent, err := agentService.GetAgentByID(uint(id))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"agent_id":		id,
			"client_ip":	c.ClientIP(),
			"error":		err.Error(),
		}).Warn("Agent not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	logger.WithFields(logrus.Fields{
		"agent_id":		id,
		"agent_name":	agent.Name,
		"client_ip":	c.ClientIP(),
	}).Info("Agent retrieved succcessfully")
	c.JSON(http.StatusOK, agent)
}

