package main

import (
	"final/db"
	"final/initializers"
	"final/middlewares"
	"final/routers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	err := db.Connect()
	if err != nil {
		initializers.LogError("connection to db", err, nil)
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Use(middlewares.LoggerMiddleware)
	routers.HomeRouter(r.PathPrefix("/").Subrouter())
	routers.SupportRouter(r.PathPrefix("/support").Subrouter())
	routers.MenuRouter(r.PathPrefix("/menu").Subrouter())
	routers.ProfileRouter(r.PathPrefix("/profile").Subrouter())
	routers.AuthRouter(r.PathPrefix("/auth").Subrouter())
	routers.VerificationRouter(r.PathPrefix("/verify").Subrouter())
	routers.AdminRouter(r.PathPrefix("/admin").Subrouter())

	PORT := initializers.GetPort()
	err = http.ListenAndServe(PORT, r)
	if err != nil {
		initializers.LogError("starting the server", err, nil)
	}
	log.Info("Listening to port: " + PORT)
}

func init() {
	initializers.InitLogger()
	initializers.InitDotEnv()
}
