package models

import (
	"time"
)

// Agent represents a compromised machine connecting to the C2 server
type Agent struct {
	ID			uint			'json:"id" gorm:"primaryKey"'
	AgentID		string			'json:"agent_id" gorm:"unique;not null"'
	Hostname	string			'json:"hostname" gorm:"not null"'
	IP 			string 			'json:"ip" gorm:"not null"'
	OS 			string 			'json:"os" gorm:"not null"'
	Arch 		string 			'json:"os" gorm:"not null"'
	LastCheckin	time.Time 		'json:"last_checkin"'
	CreatedAt 	time.Time 		'json:"created_at"'
	UpdatedAt   time.Time 		'json:"updated_at'
}

// Updatecheckin updates the last check-in time of the agent
func (a *Agent) UpdateCheckin() {
	a.LastCheckin = time.Now()
}