package initializers

import (
	"fmt"
	"os"

	"github.com/nathanroberts55/beatbattle/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	// Get PostgreSQL connection parameters from environment variables
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	// Construct PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", postgresHost, postgresUser, postgresPassword, postgresDB)

	// Connect to PostgreSQL database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
