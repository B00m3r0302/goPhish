package services

import (
	"errors"

	"github.com/B00m3r0302/goPhish/backend/models"
	"gorm.io/gorm"
)

// CommandService provides methods to manage commands in the C2 framework 
type CommandService struct {
	db *gorm.DB 
}

// NewCommandService creates a new instance of CommandService
func NewCommandService(db *gorm.DB) *CommandService {
	return &CommandService{db: db}
}

// CreateCommand creates a new command for an agent to execute 
func (s *CommandService) CreateCommand(agentID uint, command string) (*models.Command, error) {
	cmd := &models.Command{
		AgentID: agentID,
		Command: command,
		Status:  "pending",
	}

	if err := s.db.Create(cmd).Error; err != nil {
		return nil, err
	}

	return cmd, nil
}

// UpdateCommandStatus updates the status of a command
func (s *CommandService) UpdateCommandStatus(commandID uint, status string) error {
	var cmd models.Command
	if err := s.db.First(&cmd, commandID).Error; err != nil {
		return err
	}
	cmd.Status = status
	return s.db.Save(&cmd).Error
}

// AppendCommandOutput appends output to a command's result
func (s *CommandService) AppendCommandOutput(commandID uint, output string) error {
	var cmd models.Command
	if err := s.db.First(&cmd, commandID).Error; err != nil {
		return err
	}
	cmd.Output += output
	return s.db.Save(&cmd).Error
}

// GetCommandByID retrieves a command by its ID
func (s *CommandService) GetCommandByID(id uint) (*models.Command, error) {
	var cmd models.Commandif err := s.db.First(&cmd, id).Error; err != nil {
		return nil, errors.ErrUnsupported
	}
	return &cmd, nil
}