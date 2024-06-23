package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Print("Error loading .env file")
		}
		log.Print("Loaded Environment Variables from .env File")
	} else if os.IsNotExist(err) {
		log.Print(".env file does not exist")
	} else {
		log.Print("An error occurred while checking for .env file")
	}
}
