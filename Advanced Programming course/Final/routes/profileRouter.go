package routes

import (
	"final/controllers"
	"final/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func ProfileRouter(r *mux.Router) {
	router := r.PathPrefix("/profile").Subrouter()
	router.StrictSlash(true)
	router.Use(middlewares.JWTAuthentication)
	router.HandleFunc("/{id}", controllers.ProfilePageHandler)
	router.HandleFunc("/{id}/orders", controllers.OrdersPageHandler).Methods(http.MethodGet)
	router.HandleFunc("/{id}/delivery", controllers.DeliveryPageHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPut)
}
