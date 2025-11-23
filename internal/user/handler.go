package user

import (
	"encoding/json"
	"log"
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

	user, err := h.service.GetUserByID(claims.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	log.Printf("DEBUG GetProfile: User ID=%d, Profile=%+v", user.ID, user.Profile)

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
		RestingHR int     `json:"resting_hr"`
		Weight    float64 `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	user, err := h.service.GetUserByID(claims.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	log.Printf("DEBUG UpdateProfile: Before - Profile=%+v", user.Profile)

	if user.Profile != nil {
		user.Profile.FirstName = req.FirstName
		user.Profile.LastName = req.LastName
		user.Profile.Bio = req.Bio
		user.Profile.FTP = req.FTP
		user.Profile.MaxHR = req.MaxHR
		user.Profile.RestingHR = req.RestingHR
		user.Profile.Weight = req.Weight

		log.Printf("DEBUG UpdateProfile: After - Profile=%+v", user.Profile)

		if err := h.service.UpdateUser(user); err != nil {
			log.Printf("DEBUG UpdateProfile: Error updating - %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to update profile"})
			return
		}

		user, err = h.service.GetUserByID(claims.UserID)
		if err != nil {
			log.Printf("DEBUG UpdateProfile: Error reloading - %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to load profile"})
			return
		}
	}

	log.Printf("DEBUG UpdateProfile: Final - Profile=%+v", user.Profile)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetProfileResponse{
		User:    user,
		Profile: user.Profile,
	})
}