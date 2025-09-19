package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/rideaware/rideaware-api/services"
)

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {
	userService := services.NewUserService(db)

	auth := r.Group("/auth")
	{
		auth.POST("/signup", signup(userService))
		auth.POST("/login", login(userService))
		auth.POST("/logout", logout())
	}
}

func signup(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		user, err := userService.CreateUser(req.Username, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "User created successfully",
			"username": user.Username,
			"email":    user.Email,
		})
	}
}

func login(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.VerifyUser(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set session
		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user_id": user.ID,
		})
	}
}

func logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}
