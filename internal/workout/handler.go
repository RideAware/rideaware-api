package workout

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"rideaware/internal/config"
	"rideaware/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

// CreateWorkout POST /api/protected/workouts
func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateWorkout handler called")
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)
	
	if claims == nil {
		log.Printf("Claims is nil")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	log.Printf("UserID: %d", claims.UserID)

	var req struct {
		Title        string           `json:"title"`
		Description  string           `json:"description"`
		Type         string           `json:"type"`
		ScheduledDate string           `json:"scheduled_date"`
		Duration     int              `json:"duration"`
		Notes        string           `json:"notes"`
		WorkoutData  *WorkoutDataJSON `json:"workout_data"`
		FileType     string           `json:"file_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	log.Printf("CreateWorkout request: %+v", req)

	if req.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "title is required"})
		return
	}

	if req.ScheduledDate == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "scheduled_date is required"})
		return
	}

	// Parse scheduled date
	scheduledDate, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		log.Printf("Date parse error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid date format"})
		return
	}

	// Set default duration if not provided
	if req.Duration <= 0 {
		req.Duration = 60
	}

	// Initialize WorkoutData if nil
	workoutData := req.WorkoutData
	if workoutData == nil {
		workoutData = &WorkoutDataJSON{
			Segments: []WorkoutSegment{},
		}
	}

	// Create workout
	workout := &Workout{
		UserID:        claims.UserID,
		Title:         req.Title,
		Description:   req.Description,
		Type:          req.Type,
		Status:        "planned",
		ScheduledDate: scheduledDate,
		Duration:      req.Duration,
		Notes:         req.Notes,
		FileType:      req.FileType,
		WorkoutData:   *workoutData,
	}

	log.Printf("Creating workout: %+v", workout)

	if err := h.service.repo.CreateWorkout(workout); err != nil {
		log.Printf("CreateWorkout error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create workout"})
		return
	}

	log.Printf("Workout created successfully with ID: %d", workout.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

// GetWorkouts GET /api/protected/workouts
func (h *Handler) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	workouts, err := h.service.GetUserWorkouts(claims.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch workouts"})
		return
	}

	if workouts == nil {
		workouts = []Workout{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

// GetWorkoutsByMonth GET /api/protected/workouts/month?year=2025&month=11
func (h *Handler) GetWorkoutsByMonth(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid year"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid month"})
		return
	}

	workouts, err := h.service.GetWorkoutsByMonth(claims.UserID, year, month)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch workouts"})
		return
	}

	if workouts == nil {
		workouts = []Workout{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

// UpdateWorkout PUT /api/protected/workouts
func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid workout id"})
		return
	}

	var req struct {
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		Type           string  `json:"type"`
		Status         string  `json:"status"`
		Duration       int     `json:"duration"`
		Distance       float64 `json:"distance"`
		ElevGain       int     `json:"elev_gain"`
		AvgPower       int     `json:"avg_power"`
		AvgHR          int     `json:"avg_hr"`
		MaxPower       int     `json:"max_power"`
		MaxHR          int     `json:"max_hr"`
		CaloriesBurned int     `json:"calories_burned"`
		Notes          string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	workout, err := h.service.repo.GetWorkoutByID(uint(id), claims.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "workout not found"})
		return
	}

	if req.Title != "" {
		workout.Title = req.Title
	}
	if req.Description != "" {
		workout.Description = req.Description
	}
	if req.Type != "" {
		workout.Type = req.Type
	}
	if req.Status != "" {
		workout.Status = req.Status
	}
	if req.Duration > 0 {
		workout.Duration = req.Duration
	}
	if req.Distance > 0 {
		workout.Distance = req.Distance
	}
	if req.ElevGain > 0 {
		workout.ElevGain = req.ElevGain
	}
	if req.AvgPower > 0 {
		workout.AvgPower = req.AvgPower
	}
	if req.AvgHR > 0 {
		workout.AvgHR = req.AvgHR
	}
	if req.MaxPower > 0 {
		workout.MaxPower = req.MaxPower
	}
	if req.MaxHR > 0 {
		workout.MaxHR = req.MaxHR
	}
	if req.CaloriesBurned > 0 {
		workout.CaloriesBurned = req.CaloriesBurned
	}
	if req.Notes != "" {
		workout.Notes = req.Notes
	}

	if err := h.service.repo.UpdateWorkout(workout); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to update workout"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// DeleteWorkout DELETE /api/protected/workouts
func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid workout id"})
		return
	}

	if err := h.service.DeleteWorkout(uint(id), claims.UserID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// GetWorkoutTypes GET /api/protected/workout-types
func (h *Handler) GetWorkoutTypes(w http.ResponseWriter, r *http.Request) {
	types := []map[string]interface{}{
		{"id": 1, "name": "Recovery", "color": "#4285F4", "icon": "🔵"},
		{"id": 2, "name": "Endurance", "color": "#34A853", "icon": "🟢"},
		{"id": 3, "name": "Tempo", "color": "#FBBC04", "icon": "🟡"},
		{"id": 4, "name": "Threshold", "color": "#EA4335", "icon": "🔴"},
		{"id": 5, "name": "VO2 Max", "color": "#A61C00", "icon": "⭐"},
		{"id": 6, "name": "Strength", "color": "#800080", "icon": "💪"},
		{"id": 7, "name": "Race", "color": "#FF1744", "icon": "🏁"},
		{"id": 8, "name": "Rest", "color": "#CCCCCC", "icon": "😴"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

// UploadWorkoutFile POST /api/protected/workouts/upload
func (h *Handler) UploadWorkoutFile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "file too large"})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "no file provided"})
		return
	}
	defer file.Close()

	fileContent := make([]byte, handler.Size)
	file.Read(fileContent)

	// Parse ZWO file
	parsedData, err := ParseZWO(fileContent)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Get scheduled date from form
	scheduledDateStr := r.FormValue("scheduled_date")
	scheduledDate, err := time.Parse("2006-01-02", scheduledDateStr)
	if err != nil {
		scheduledDate = time.Now()
	}

	// Create workout with parsed data
	workout := &Workout{
		UserID:        claims.UserID,
		Title:         parsedData.Name,
		Description:   parsedData.Description,
		Type:          "imported",
		Status:        "planned",
		ScheduledDate: scheduledDate,
		Duration:      parsedData.TotalDuration,
		FileType:      "zwo",
		WorkoutData: WorkoutDataJSON{
			Name:          parsedData.Name,
			Author:        parsedData.Author,
			TotalDuration: parsedData.TotalDuration,
			Segments:      parsedData.Segments,
		},
	}

	if err := h.service.repo.CreateWorkout(workout); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create workout"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}