package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"rideaware/internal/auth"
	"rideaware/internal/config"
	"rideaware/internal/equipment"
	"rideaware/internal/middlewares"
	"rideaware/internal/user"
	"rideaware/pkg/database"
)

func main() {
	godotenv.Load()

	// Initialize database
	database.Init()
	defer database.Close()

	// Run migrations
	if err := database.Migrate(
		&user.User{},
		&user.Profile{},
		&user.PasswordReset{},
		&user.Session{},
		&equipment.Equipment{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize JWT config
	config.InitJWT()

	r := chi.NewRouter()

	// Logging middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(loggingMiddleware)
	r.Use(middleware.Recoverer)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept", "Authorization", "Content-Type",
		},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	// Routes
	setupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"[%s] %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			r.Header.Get("User-Agent"),
		)
		next.ServeHTTP(w, r)
	})
}

func setupRoutes(r *chi.Mux) {
	// Public routes
	r.Get("/health", healthCheck)

	// Auth routes
	authHandler := auth.NewHandler()
	r.Post("/signup", authHandler.Signup)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.Post("/password-reset/request", authHandler.RequestPasswordReset)
	r.Post("/password-reset/confirm", authHandler.ConfirmPasswordReset)
	r.Post("/refresh-token", authHandler.RefreshToken)

	// Protected routes
	authMiddleware := middlewares.NewAuthMiddleware()
	r.Route("/protected", func(r chi.Router) {
		r.Use(authMiddleware.ProtectedRoute)

		// User routes
		userHandler := user.NewHandler()
		r.Get("/profile", userHandler.GetProfile)
		r.Put("/profile", userHandler.UpdateProfile)

		// Equipment routes
		equipmentHandler := equipment.NewHandler()
		r.Post("/equipment", equipmentHandler.CreateEquipment)
		r.Get("/equipment", equipmentHandler.GetEquipment)
		r.Put("/equipment", equipmentHandler.UpdateEquipment)
		r.Delete("/equipment", equipmentHandler.DeleteEquipment)

		// Training zones
		r.Get("/zones", equipmentHandler.GetTrainingZones)
	})

	log.Println("✅ Routes registered successfully")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("📊 Health check called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}