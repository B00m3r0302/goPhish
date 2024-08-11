package models

import (
	"time"
)

// User represents and administrator or operator within the C2 framework
type User struct {
	ID			uint		'json:"id" gorm:"primaryKey"'
	Username	string		'json:"username" gorm:"unique;not null"'
	Password	string		'json:"password" gorm:"not null"'
	Role		string		'json:"role" gorm:"not null"'
	CreatedAt	time.Time	'json:"created_at"'
	UpdatedAt	time.Time	'json:"updated_at"'
}

// HashPassword hashes the user's password before storing it in the database
func (u *User) HashPassword(password string) error {
	// Implement a password hashing logic e.g, using Bcrypt
	return nil
}

// CheckPassword checks the provided password against the stored hash
func (u *User) CheckPassword(password string) bool {
	// Implement password verification logic e.g using bcrypt
	return true
}