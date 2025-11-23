package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"rideaware/internal/config"
	"rideaware/internal/user"
)

type Handler struct {
	userService *user.Service
}

func NewHandler() *Handler {
	return &Handler{
		userService: user.NewService(),
	}
}

type SignupRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	log.Println("📝 Signup request received")
	
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Signup decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	log.Printf("📝 Signup attempt for user: %s (email: %s)", req.Username, req.Email)

	newUser, err := h.userService.CreateUser(req.Username, req.Password, req.Email, req.FirstName, req.LastName)
	if err != nil {
		log.Printf("❌ Signup error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("✅ User created: %s (ID: %d)", newUser.Username, newUser.ID)

	accessToken, _ := config.GenerateAccessToken(newUser.ID, newUser.Email, newUser.Username)
	refreshToken, _ := config.GenerateRefreshToken(newUser.ID, newUser.Email, newUser.Username)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		UserID:       newUser.ID,
		Username:     newUser.Username,
		Email:        newUser.Email,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("🔐 Login request received")
	
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Login decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	log.Printf("🔐 Login attempt for user: %s", req.Username)

	user, err := h.userService.VerifyUser(req.Username, req.Password)
	if err != nil {
		log.Printf("❌ Login error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("✅ Login successful for user: %s (ID: %d)", user.Username, user.ID)

	accessToken, _ := config.GenerateAccessToken(user.ID, user.Email, user.Username)
	refreshToken, _ := config.GenerateRefreshToken(user.ID, user.Email, user.Username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
	})
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log.Println("🔄 Refresh token request received")
	
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Refresh token decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	log.Println("🔄 Verifying refresh token...")

	// Verify refresh token and get user
	claims, err := config.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		log.Printf("❌ Refresh token verify error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid refresh token"})
		return
	}

	log.Printf("✅ Refresh token valid for user ID: %d", claims.UserID)

	// Generate new access token
	newAccessToken, _ := config.GenerateAccessToken(claims.UserID, claims.Email, claims.Username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": newAccessToken,
		"expires_in":   900,
	})
}

func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	log.Println("🔑 Password reset request received")
	
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Password reset decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	log.Printf("🔑 Password reset requested for email: %s", req.Email)

	err := h.userService.RequestPasswordReset(req.Email)
	if err != nil {
		log.Printf("❌ Password reset error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("✅ Password reset email sent to: %s", req.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If email exists, reset link has been sent",
	})
}

func (h *Handler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	log.Println("🔑 Password reset confirm request received")
	
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Password reset confirm decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	log.Println("🔑 Confirming password reset...")

	if err := h.userService.ResetPassword(req.Token, req.NewPassword); err != nil {
		log.Printf("❌ Password reset confirm error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Println("✅ Password reset successful")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successful",
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("👋 Logout request received")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}