package services

import (
	"errors"

	"github.com/B00m3r0302/goPhish/backend/models"
	"gorm.io/gorm"
)

// ListenerService provides methods to manage network listeners in the C2 framework
type ListenerService struct {
	db *gorm.DB
}

// NewListenerService creates a new instalce of ListenerService
func NewListenerService(db *gorm.DB) *ListenerService {
	return &ListenerService{db: db}
}

// CreateListener creates a new network listener with the specified configuration
func (s *ListenerService) CreateListener(address, protocol string, maxConnections, timeout int) (*models.Listener, error) {
	listener := &models.Listener{
		Address:        address,
		Protocol:       protocol,
		MaxConnections: maxConnections,
		Timeout:        timeout,
	}

	if err := s.db.Create(listener).Error; err != nil {
		return nil, err
	}

	return listener, nil
}

// StartListener starts the listener and begins accepting connections
func (s *ListenerService) StartListener(id uint) error {
	var listener models.Listener
	if err := s.db.First(&listener, id).Error; err != nil {
		return errors.New("Listener not found")
	}

	return listener.Start()
}

// StopListener stops the listener from accepting further connections
func (s *ListenerService) StopListener(id uint) error {
	var listener models.Listener
	if err := s.db.First(&listener, id).Error; err != nil {
		return errors.New("Listener not found")
	}

	return listener.Stop()
}

// GetListenerByID retrieves a listener by its ID
func (s *ListenerService) GetListenerByID(id uint) (*models.Listener, error) {
	var listener models.Listener
	if err := s.db.First(&listener, id).Error; err != nil {
		return nil, err
	}

	return &listener, nil
}
