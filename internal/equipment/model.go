package equipment

import "time"

type Equipment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null" json:"type"` // "bike", "shoes", "helmet", etc.
	Brand     string    `gorm:"default:''" json:"brand"`
	Model     string    `gorm:"default:''" json:"model"`
	Weight    float64   `gorm:"default:0" json:"weight"` // grams
	Notes     string    `gorm:"default:''" json:"notes"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TrainingZone struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Min   int    `json:"min"`
	Max   int    `json:"max"`
	Color string `json:"color"`
}

type HRZones struct {
	Zone1 TrainingZone `json:"zone_1"` // Recovery
	Zone2 TrainingZone `json:"zone_2"` // Endurance
	Zone3 TrainingZone `json:"zone_3"` // Tempo
	Zone4 TrainingZone `json:"zone_4"` // Threshold
	Zone5 TrainingZone `json:"zone_5"` // VO2 Max
}

type PowerZones struct {
	Zone1 TrainingZone `json:"zone_1"` // Active Recovery
	Zone2 TrainingZone `json:"zone_2"` // Endurance
	Zone3 TrainingZone `json:"zone_3"` // Sweet Spot
	Zone4 TrainingZone `json:"zone_4"` // Threshold
	Zone5 TrainingZone `json:"zone_5"` // VO2 Max
	Zone6 TrainingZone `json:"zone_6"` // Anaerobic
	Zone7 TrainingZone `json:"zone_7"` // Neuromuscular
}