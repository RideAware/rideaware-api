package models

type UserProfile struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	UserID         uint   `gorm:"not null" json:"user_id"`
	FirstName      string `gorm:"size:80;not null" json:"first_name"`
	LastName       string `gorm:"size:80;not null" json:"last_name"`
	Bio            string `gorm:"type:text" json:"bio"`
	ProfilePicture string `gorm:"size:255" json:"profile_picture"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
