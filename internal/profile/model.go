package profile

import "time"

type Equipment struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	UserID uint      `gorm:"not null;index" json:"user_id"`
	Name   string    `gorm:"not null" json:"name"`
	Type   string    `gorm:"not null" json:"type"` // "bike", "shoes", "helmet", etc.
	Brand  string    `gorm:"default:''" json:"brand"`
	Model  string    `gorm:"default:''" json:"model"`
	Weight float64   `gorm:"default:0" json:"weight"` // grams
	Notes  string    `gorm:"default:''" json:"notes"`
	Active bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Stats struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	TotalRides    int   `gorm:"default:0" json:"total_rides"`
	TotalDistance float64 `gorm:"default:0" json:"total_distance"`
	TotalTime     int   `gorm:"default:0" json:"total_time"`
	AverageSpeed  float64 `gorm:"default:0" json:"average_speed"`
	MaxSpeed      float64 `gorm:"default:0" json:"max_speed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}