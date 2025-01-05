package config

import (
	"fmt"
	"log"
	"os"

	"taas-api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Fetch environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Build DSN string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Card{}, &models.Booking{}, &models.VideoControl{}, &models.Users_ref{}, &models.TalentRegistration{}, &models.ServiceCard{}, &models.AvailableTimeSlots{}, &models.BookingRequests{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	DB = db
	fmt.Println("Database connected and schema migrated")
}
