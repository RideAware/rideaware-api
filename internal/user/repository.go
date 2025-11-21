package user

import (
	"errors"
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
	if err := database.DB.Preload("Profile").Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *User) error {
	return database.DB.Save(user).Error
}

func (r *Repository) UserExists(username, email string) (bool, error) {
	var count int64
	err := database.DB.Model(&User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error
	return count > 0, err
}