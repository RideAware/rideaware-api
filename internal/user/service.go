package user

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"
	"time"

	"rideaware/internal/config"
	"rideaware/internal/email"
	"rideaware/pkg/database"
)

type Service struct {
	repo  *Repository
	email *email.Service
}

func NewService() *Service {
	return &Service{
		repo:  NewRepository(),
		email: email.NewService(),
	}
}

func (s *Service) CreateUser(username, password, email, firstName, lastName string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	if email != "" {
		if !isValidEmail(email) {
			return nil, errors.New("invalid email format")
		}
	}

	exists, err := s.repo.UserExists(username, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username or email already exists")
	}

	user := &User{
		Username: username,
		Email:    email,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	_ = s.email.SendWelcomeEmail(email, username)

	return user, nil
}

func (s *Service) VerifyUser(username, password string) (*User, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *Service) RequestPasswordReset(email string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		// Don't leak if email exists
		return nil
	}

	token, err := generateSecureToken(32)
	if err != nil {
		return err
	}

	resetToken := &PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(config.JWT.ResetTokenDuration),
	}

	if err := database.DB.Create(resetToken).Error; err != nil {
		return err
	}

	resetLink := "https://rideaware.app/reset-password?token=" + token
	return s.email.SendPasswordResetEmail(user.Email, user.Username, resetLink)
}

func (s *Service) ResetPassword(token, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var resetToken PasswordReset
	if err := database.DB.Where("token = ?", token).First(&resetToken).Error; err != nil {
		return errors.New("invalid or expired reset token")
	}

	if !resetToken.IsValid() {
		return errors.New("reset token has expired")
	}

	user, err := s.repo.GetUserByID(resetToken.UserID)
	if err != nil {
		return err
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	now := time.Now()
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(user).Update("password", user.Password).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&resetToken).Update("used_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Helper functions
func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}