package initializers

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func InitDotEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env files")
	}
}

func GetPort() string {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		return ":4444"
	}
	return ":" + PORT
}

func GetConnectionType() string {
	CONN_TYPE := os.Getenv("CONN_TYPE")
	if CONN_TYPE == "" {
		return "tcp"
	}
	return CONN_TYPE
}
