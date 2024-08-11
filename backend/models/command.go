package models

import (
	"time"
)

// Command represents a command sent to an agent for execution
type Command struct {
	ID			uint		'json:"id" gorm:"primaryKey"'
	AgentID		uintchan	'json:"agent_id" gorm:"not null"'
	Command 	string		'json:"command" gorm:"not null"'
	Status 		string 		'json:"status" gorm:"not null"'
	Output 		string 		'json:"output"'
	CreatedAt 	time.Time 	'json:"created_at"'
	UpdatedAt 	time.Time	'json:"updated_at"'
}

// SetStatus updates the status of the command (e.g, pending, executed, failed)
func (c *Command) SetStatus(status string) {
	c.Status = status
}

// AppendOutput appends output from the agent's command execution
func (c *Command) AppendOutput(output string) {
	c.Output += output 
}