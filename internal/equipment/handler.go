package equipment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rideaware/internal/config"
	"rideaware/internal/middleware"
	"rideaware/internal/user"
)

type Handler struct {
	service     *Service
	userService *user.Service
}

func NewHandler() *Handler {
	return &Handler{
		service:     NewService(),
		userService: user.NewService(),
	}
}

// CreateEquipment POST /api/equipment
func (h *Handler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	var req struct {
		Name   string  `json:"name"`
		Type   string  `json:"type"`
		Brand  string  `json:"brand"`
		Model  string  `json:"model"`
		Weight float64 `json:"weight"`
		Notes  string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	equipment, err := h.service.CreateEquipment(
		claims.UserID,
		req.Name,
		req.Type,
		req.Brand,
		req.Model,
		req.Weight,
		req.Notes,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(equipment)
}

// GetEquipment GET /api/equipment
func (h *Handler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	equipment, err := h.service.GetUserEquipment(claims.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch equipment"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

// UpdateEquipment PUT /api/equipment
func (h *Handler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid equipment id"})
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	equipment, err := h.service.UpdateEquipment(uint(id), claims.UserID, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

// DeleteEquipment DELETE /api/equipment
func (h *Handler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid equipment id"})
		return
	}

	if err := h.service.DeleteEquipment(uint(id), claims.UserID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// GetTrainingZones GET /api/zones
func (h *Handler) GetTrainingZones(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	profile, err := h.userService.GetUserByID(claims.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	// Check if profile exists
	if profile == nil || profile.Profile == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user profile not found"})
		return
	}

	zones := map[string]interface{}{}

	if profile.Profile.MaxHR > 0 {
		zones["hr_zones"] = h.service.CalculateHRZones(profile.Profile.MaxHR, profile.Profile.RestingHR)
	}

	if profile.Profile.FTP > 0 {
		zones["power_zones"] = h.service.CalculatePowerZones(profile.Profile.FTP)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}