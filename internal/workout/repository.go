package workout

import (
	"errors"
	"rideaware/pkg/database"
	"time"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateWorkout(workout *Workout) error {
	return database.DB.Create(workout).Error
}

func (r *Repository) GetWorkoutByID(id, userID uint) (*Workout, error) {
	var workout Workout
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).
		First(&workout).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workout not found")
		}
		return nil, err
	}
	return &workout, nil
}

func (r *Repository) GetUserWorkouts(userID uint) ([]Workout, error) {
	var workouts []Workout
	if err := database.DB.Where("user_id = ?", userID).
		Order("scheduled_date DESC").
		Find(&workouts).Error; err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *Repository) GetWorkoutsByDateRange(userID uint, start, end time.Time) ([]Workout, error) {
	var workouts []Workout
	if err := database.DB.Where("user_id = ? AND scheduled_date BETWEEN ? AND ?", userID, start, end).
		Order("scheduled_date ASC").
		Find(&workouts).Error; err != nil {
		return nil, err
	}
	return workouts, nil
}

func (r *Repository) GetWorkoutsByMonth(userID uint, year, month int) ([]Workout, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return r.GetWorkoutsByDateRange(userID, start, end)
}

func (r *Repository) UpdateWorkout(workout *Workout) error {
	return database.DB.Model(workout).Updates(workout).Error
}

func (r *Repository) DeleteWorkout(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).
		Delete(&Workout{}).Error
}