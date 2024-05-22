package main

import (
	_ "github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	if err := loadEnvVariables(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	router := setupRouter(callOpenAI)
	log.Fatal(http.ListenAndServe(":8080", router))
}
