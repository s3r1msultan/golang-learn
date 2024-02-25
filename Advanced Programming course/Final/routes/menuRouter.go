package routes

import (
	"final/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func MenuRouter(r *mux.Router) {
	menuRouter := r.PathPrefix("/menu").Subrouter()
	menuRouter.HandleFunc("", controllers.MenuPageHandler).Methods(http.MethodGet)
	menuRouter.HandleFunc("/{id}", controllers.DishPageHandler).Methods(http.MethodGet)
}
