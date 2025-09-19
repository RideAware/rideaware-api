package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	database := os.Getenv("PG_DATABASE")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")

	dsn := fmt.Sprintf(`host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable`,
		host, port, user, password, database)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB from gorm:", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	return db
}
