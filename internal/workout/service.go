package workout

import (
	"errors"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

func (s *Service) CreateWorkout(userID uint, title string, scheduledDate time.Time, duration int) (*Workout, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	workout := &Workout{
		UserID:        userID,
		Title:         title,
		ScheduledDate: scheduledDate,
		Duration:      duration,
		Status:        "planned",
	}

	if err := s.repo.CreateWorkout(workout); err != nil {
		return nil, err
	}

	return workout, nil
}

func (s *Service) GetUserWorkouts(userID uint) ([]Workout, error) {
	return s.repo.GetUserWorkouts(userID)
}

func (s *Service) GetWorkoutsByMonth(userID uint, year, month int) ([]Workout, error) {
	return s.repo.GetWorkoutsByMonth(userID, year, month)
}

func (s *Service) UpdateWorkoutStatus(id, userID uint, status string) (*Workout, error) {
	if status != "planned" && status != "completed" && status != "skipped" {
		return nil, errors.New("invalid status")
	}

	workout, err := s.repo.GetWorkoutByID(id, userID)
	if err != nil {
		return nil, err
	}

	workout.Status = status
	if status == "completed" {
		workout.UpdatedAt = time.Now()
	}

	if err := s.repo.UpdateWorkout(workout); err != nil {
		return nil, err
	}

	return workout, nil
}

func (s *Service) UpdateWorkoutWithMetrics(id, userID uint, distance float64, avgPower, avgHR int) (*Workout, error) {
	workout, err := s.repo.GetWorkoutByID(id, userID)
	if err != nil {
		return nil, err
	}

	workout.Distance = distance
	workout.AvgPower = avgPower
	workout.AvgHR = avgHR
	workout.Status = "completed"

	if err := s.repo.UpdateWorkout(workout); err != nil {
		return nil, err
	}

	return workout, nil
}

func (s *Service) DeleteWorkout(id, userID uint) error {
	return s.repo.DeleteWorkout(id, userID)
}