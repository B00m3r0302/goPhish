package models

import (
	"time"
)

// Task represents a task assigned to an agent for execution
type Task struct {
	ID			uint 			'json:"id: gorm:"primaryKey"'
	AgentID 	uint 			'json:"agent_id" gorm:"not null"'
	Command 	string 			'json:"command" gorm:"not null"'
	Status 		string 			'json:"status" gorm:"not null"'
	Result 		string 			'json:"result"'
	AssignedAt 	time.Time 		'json:"assigned_at"'
	CompletedAt	time.Time 		'json:"completed_at"'
}

// MarkAsCompleted marks the task as completed and sets the completion time
func (t *Task) MarkAsCompleted(result string) {
	t.Status = "completed"
	t.Result = result 
	t.CompletedAt = time.Now()	
}

// IsPending checks if the task is still pending
func (t *Task) IsPending() bool {
	return t.Status == "pending"
}