package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	database := os.Getenv("PG_DATABASE")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")

	// Try with quoted password
	dsn := fmt.Sprintf(`host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable`,
		host, port, user, password, database)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}
