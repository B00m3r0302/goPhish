package services

import {
	"github.com/B00m3r0302/goPhish/backend/models"
	"gorm.io/gorm"
}

// SessionService provides methods to manage sessions in the C2 framework
type SessionService struct {
	db *gorm.DB 
}

// NewSessionService creates a new instance of SessionService
func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

// StartSession starts a new session between a user and an agent
func (s *SessionService) StartSession(userID, agentID uint) (*models.Session, error) {
	session := &models.Session{
		UserID:  userID,
		AgentID: agentID,
		Active:  true,
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, err
	} 

	return session, nil
}

// EndSession ends an active session
func (s *SessionService) EndSession(sessionID uint) error {
	var session models.Session
	if err := s.db.First(&session, sessionID).Error; err != nil {
		return err
	}

	session.EndSession()
	return s.db.Save(&session).Error
}

// GetActiveSessions retrieves all active sessions
func (s *SessionService) GetActiveSessions() ([]models.Session, error) {
	var sessions []models.Session
	if err := s.db.Where("active = ?", true).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}