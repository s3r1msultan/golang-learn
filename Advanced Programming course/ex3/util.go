package main

import (
	"github.com/joho/godotenv"
	"strings"
)

func loadEnvVariables() error {
	return godotenv.Load()
}

func isValidRequest(message string) bool {
	validTopics := []string{"law", "regulations", "statutes", "legal", "compliance"}
	for _, topic := range validTopics {
		if strings.Contains(strings.ToLower(message), topic) {
			return true
		}
	}
	return false
}
