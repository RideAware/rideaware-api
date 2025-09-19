package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null;size:80" json:"username"`
	Email    string `gorm:"unique;not null;size:255" json:"email"` // Add this line
	Password string `gorm:"not null;size:255" json:"-"`

	Profile *UserProfile `gorm:"constraint:OnDelete:CASCADE;" json:"profile,omitempty"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	profile := UserProfile{
		UserID:         u.ID,
		FirstName:      "",
		LastName:       "",
		Bio:            "",
		ProfilePicture: "",
	}
	return tx.Create(&profile).Error
}
