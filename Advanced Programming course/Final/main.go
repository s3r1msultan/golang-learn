package main

import (
	"final/controllers"
	"final/db"
	"final/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
)

var dbClient *mongo.Client

func main() {
	initLogger()
	initDotEnv()
	dbClient, _ = db.Connect()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", controllers.HomePageHandler).Methods(http.MethodGet)
	r.HandleFunc("/support", controllers.SupportPageHandler).Methods(http.MethodGet, http.MethodPost)
	routes.MenuRouter(r)
	routes.ProfileRouter(r)
	routes.AuthRouter(r)
	routes.VerificationRouter(r)

	PORT := getPort()
	err := http.ListenAndServe(PORT, r)
	if err != nil {
		log.Fatal("Error starting the server on the port " + PORT)
	}
	log.Info("Listening to port: " + PORT)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func getPort() string {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		return ":8080"
	}
	return ":" + PORT
}

func initDotEnv() {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		log.Fatal("Error loading .env files")
	}
}

func initLogger() {
	logFile := "log.txt"
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to create logfile" + logFile)
		panic(err)
	}
	//defer f.Close()
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
}
