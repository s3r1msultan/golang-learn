package main

import (
	"log"
	"net/http"

	_ "github.com/julienschmidt/httprouter"
)

func main() {
	if err := loadEnvVariables(); err != nil {

		log.Fatal("Error loading .env file:", err)
	}
	router := setupRouter(callOpenAI)
	log.Fatal(http.ListenAndServe(":3001", router))
}
