package services

import (
	"errors"
	"log"
	"strings"

	"github.com/rideaware/rideaware-api/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	if len(username) < 3 || len(password) < 8 {
		return nil, errors.New("username must be at least 3 characters and password must be at least 8 characters")
	}

	// Basic email validation
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email format")
	}

	// Check if user exists (by username or email)
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this username or email already exists")
	}

	// Create new user
	user := models.User{
		Username: username,
		Email:    email,
	}
	if err := user.SetPassword(password); err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("could not create user")
	}

	if err := s.db.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, errors.New("could not create user")
	}

	return &user, nil
}

func (s *UserService) VerifyUser(username, password string) (*models.User, error) {
	var user models.User
	// Allow login with either username or email
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		log.Printf("User not found: %s", username)
		return nil, errors.New("invalid username or password")
	}

	if !user.CheckPassword(password) {
		log.Printf("Invalid password for user: %s", username)
		return nil, errors.New("invalid username or password")
	}

	log.Printf("User verified: %s", username)
	return &user, nil
}
