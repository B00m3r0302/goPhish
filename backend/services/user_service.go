package services

import (
	"errors"

	"github.com/B00m3r0302/goPhish/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService provides the methods to manage users in the C2 framework
type UserService struct {
	db *gorm.db
}

// NewUserService creates a new instance of UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user with the specified username, password and role
func (s *UserService) CreateUser(username, password, role string) (*models.User), error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser checks the username and password, returning the user if valid
func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("Invalid username or password")
	}

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("Invalid username or password")
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (s*UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserRole updates the role of an existing user
func (s *UserService) UpdateUserRole(id uint, role string) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}
	user.Role = role
	return s.db.Save(&user).Error
}