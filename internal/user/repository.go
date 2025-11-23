package user

import (
	"errors"
	"log"
	"rideaware/pkg/database"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateUser(user *User) error {
	return database.DB.Create(user).Error
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id uint) (*User, error) {
	var user User
	
	// Get the user
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Manually load the profile
	var profile Profile
	if err := database.DB.Where("user_id = ?", id).First(&profile).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error loading profile: %v", err)
		}
		// Profile might not exist, that's okay
	} else {
		user.Profile = &profile
	}

	log.Printf("DEBUG: Loaded user %d, profile ID=%d, profile=%+v", id, profile.ID, user.Profile)

	return &user, nil
}

func (r *Repository) UpdateUser(user *User) error {
	// Update the user
	if err := database.DB.Model(user).Updates(user).Error; err != nil {
		return err
	}

	// Update the profile if it exists
	if user.Profile != nil {
		if err := database.DB.Model(user.Profile).Updates(user.Profile).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) UserExists(username, email string) (bool, error) {
	var count int64
	err := database.DB.Model(&User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error
	return count > 0, err
}