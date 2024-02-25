package routes

import (
	"final/controllers"
	"github.com/gorilla/mux"
)

func VerificationRouter(r *mux.Router) {
	router := r.PathPrefix("/verify").Subrouter()
	router.StrictSlash(true)
	router.HandleFunc("", controllers.VerificationHandler)
}
