package models

import (
	"time"
)

// Session represents a session between an operator and an agent
type Session struct {
	ID			uint			'json:"id" gorm:"primaryKey"'
	UserID 		uint			'json:"user_id" gorm:"not null"'
	AgentID 	uint 			'json:"agent_id" gorm:"not null"'
	StartTime 	time.Time 		'json:"start_time" gorm:"not null"'
	EndTime 	time.Time 		'json:";end_time"'
	Active 		bool 			'json:"active" gorm:"not null"'
}

// StartSession sets the session as active and records the start time
func (s *Session) StartSession() {
	s.Active = true 
	s.StartTime = time.Now()
}

// EndSession sets the session as inactive and records he end time
func (s *Session) EndSession() {
	s.Active = false 
	s.Endtime = time.Now()
}