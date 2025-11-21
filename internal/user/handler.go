package user

import (
	"encoding/json"
	"net/http"

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

type GetProfileResponse struct {
	User    *User    `json:"user"`
	Profile *Profile `json:"profile"`
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	user, err := h.service.repo.GetUserByID(claims.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetProfileResponse{
		User:    user,
		Profile: user.Profile,
	})
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserContextKey).(*config.CustomClaims)

	var req struct {
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Bio       string  `json:"bio"`
		FTP       int     `json:"ftp"`
		MaxHR     int     `json:"max_hr"`
		Weight    float64 `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	user, err := h.service.repo.GetUserByID(claims.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	// Update profile
	if user.Profile != nil {
		user.Profile.FirstName = req.FirstName
		user.Profile.LastName = req.LastName
		user.Profile.Bio = req.Bio
		user.Profile.FTP = req.FTP
		user.Profile.MaxHR = req.MaxHR
		user.Profile.Weight = req.Weight

		if err := h.service.repo.UpdateUser(user); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to update profile"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetProfileResponse{
		User:    user,
		Profile: user.Profile,
	})
}