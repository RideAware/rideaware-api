package workout

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Workout struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	UserID         uint            `gorm:"not null;index" json:"user_id"`
	Title          string          `gorm:"not null" json:"title"`
	Description    string          `gorm:"default:''" json:"description"`
	Type           string          `gorm:"default:''" json:"type"`
	Status         string          `gorm:"default:'planned'" json:"status"`
	ScheduledDate  time.Time       `gorm:"index" json:"scheduled_date"`
	Duration       int             `gorm:"default:0" json:"duration"`
	Distance       float64         `gorm:"default:0" json:"distance"`
	ElevGain       int             `gorm:"default:0" json:"elev_gain"`
	AvgPower       int             `gorm:"default:0" json:"avg_power"`
	AvgHR          int             `gorm:"default:0" json:"avg_hr"`
	MaxPower       int             `gorm:"default:0" json:"max_power"`
	MaxHR          int             `gorm:"default:0" json:"max_hr"`
	CaloriesBurned int             `gorm:"default:0" json:"calories_burned"`
	FileType       string          `gorm:"default:''" json:"file_type"`
	FileURL        string          `gorm:"default:''" json:"file_url"`
	WorkoutData    WorkoutDataJSON `gorm:"type:jsonb" json:"workout_data,omitempty"`
	Notes          string          `json:"notes"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type WorkoutDataJSON struct {
	Name          string            `json:"name"`
	Author        string            `json:"author"`
	TotalDuration int               `json:"total_duration"`
	Segments      []WorkoutSegment  `json:"segments"`
}

type WorkoutSegment struct {
	Type      string  `json:"type"`
	Duration  int     `json:"duration"`
	PowerLow  float64 `json:"power_low"`
	PowerHigh float64 `json:"power_high"`
	Power     float64 `json:"power"`
	Cadence   int     `json:"cadence"`
}

// Scan implements sql.Scanner interface
func (w *WorkoutDataJSON) Scan(value interface{}) error {
	if value == nil {
		*w = WorkoutDataJSON{}
		return nil
	}
	bytes := value.([]byte)
	return json.Unmarshal(bytes, &w)
}

// Value implements driver.Valuer interface
func (w WorkoutDataJSON) Value() (driver.Value, error) {
	return json.Marshal(w)
}

func (Workout) TableName() string {
	return "workouts"
}