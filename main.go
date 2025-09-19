package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/rideaware/rideaware-api/config"
	"github.com/rideaware/rideaware-api/models"
	"github.com/rideaware/rideaware-api/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db := config.InitDB()

	// Auto migrate models
	if err := db.AutoMigrate(&models.User{}, &models.UserProfile{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.Default())

	// Session middleware
	secret := os.Getenv("SECRET_KEY")
	if len(secret) < 32 {
		log.Fatal("SECRET_KEY must be at least 32 bytes")
	}

	authKey := []byte(secret)
	encKey := []byte(secret[:32])
	store := cookie.NewStore(authKey, encKey)
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 80 * 24 * 7, // 7 days
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("rideaware-session", store))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Register auth routes
	routes.RegisterAuthRoutes(r, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
