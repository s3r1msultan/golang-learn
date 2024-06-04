package initializers

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func InitDotEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env files")
	}
}

func GetPort() string {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		return ":8080"
	}
	return ":" + PORT
}
