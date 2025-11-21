package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Profile        *Profile         `gorm:"foreignKey:UserID;constraint:OnDelete:Cascade" json:"profile,omitempty"`
	PasswordResets []PasswordReset  `gorm:"foreignKey:UserID;constraint:OnDelete:Cascade" json:"password_resets,omitempty"`
	Sessions       []Session        `gorm:"foreignKey:UserID;constraint:OnDelete:Cascade" json:"sessions,omitempty"`
}

type Profile struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	FirstName      string    `gorm:"default:''" json:"first_name"`
	LastName       string    `gorm:"default:''" json:"last_name"`
	Bio            string    `gorm:"default:''" json:"bio"`
	ProfilePicture string    `gorm:"default:''" json:"profile_picture"`
	RestingHR      int       `gorm:"default:0" json:"resting_hr"`
	MaxHR          int       `gorm:"default:0" json:"max_hr"`
	FTP            int       `gorm:"default:0" json:"ftp"`
	Weight         float64   `gorm:"default:0" json:"weight"`
	TotalRides     int       `gorm:"default:0" json:"total_rides"`
	TotalDistance  float64   `gorm:"default:0" json:"total_distance"`
	TotalTime      int       `gorm:"default:0" json:"total_time"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PasswordReset struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null" json:"user_id"`
	Token     string     `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type Session struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Token      string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt  time.Time `gorm:"not null;index" json:"expires_at"`
	DeviceName string    `gorm:"default:''" json:"device_name"`
	UserAgent  string    `gorm:"default:''" json:"user_agent"`
	IPAddress  string    `gorm:"default:''" json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
}

// ===== Methods =====

// SetPassword hashes and sets the password
func (u *User) SetPassword(rawPassword string) error {
	if len(rawPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(rawPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	) == nil
}

// AfterCreate hook: automatically create profile after user insert
func (u *User) AfterCreate(tx *gorm.DB) error {
	profile := &Profile{
		UserID: u.ID,
	}
	return tx.Create(profile).Error
}

// IsPasswordResetTokenValid checks if token exists and is not expired
func (prt *PasswordReset) IsValid() bool {
	return prt.UsedAt == nil && time.Now().Before(prt.ExpiresAt)
}

// IsSessionValid checks if session is not expired
func (s *Session) IsValid() bool {
	return time.Now().Before(s.ExpiresAt)
}