package services

import (
	"errors"
	"time"

	"github.com/B00m3r0302/goPhish/backend/models"
	"gorm.io/gorm"
)

// AgentService provides methods to manage agents in the C2 Framework
type AgentService struct {
	db *gorm.DB
}

// NewAgentService creates a new instance of AgentService
func NewAgentService(db *gorm.DB) *AgentService {
	return &AgentService{db: db}
}

// RegisterAgent registers a new agent with the provided details
func (s *AgentService) RegisterAgent(agentID, hostname, ip, os, arch string) (*models.Agent, error) {
	agent := &models.Agent{
		AgentID:  agentID,
		Hostname: hostname,
		IP:       ip,
		OS:       os,
		Arch:     arch,
	}

	if err := s.db.Create(agent).Error; err != nil {
		return nil, err
	}

	return agent, nil
}

// CheckAgent updates the last check-in time for an agent
func (s *AgentService) CheckinAgent(agentID string) error {
	var agent models.Agent
	if err := s.db.Where("agent_id = ?", agentID).First(&agent).Error; err != nil {
		return errors.New("Agent not found")
	}

	agent.LastCheckin = time.Now()
	return s.db.Save(&agent).Error
}

// GetAgentByID retrieves an agent by their ID
func (s *AgentService) GetAgentByID(id uint) (*models.Agent, error) {
	var agent models.Agent
	if err := s.db.First(&agent, id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetAgentByAgentID retrieves and agent by teir AgentID
func (s *AgentService) GetAgentByAgentID(agentID string) (*models.Agent, error) {
	var agent models.Agent
	if err := s.db.Where("agent_id = ?", agentID).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}
