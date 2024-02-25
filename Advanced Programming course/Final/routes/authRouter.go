package routes

import (
	"final/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthRouter(r *mux.Router) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("", controllers.AuthPageHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/sign_up", controllers.SignupHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/sign_in", controllers.SigninHandler).Methods(http.MethodPost)
}
