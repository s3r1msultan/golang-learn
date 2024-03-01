package routers

import (
	"final/controllers"
	"final/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func MenuRouter(menuRouter *mux.Router) {
	menuRouter.StrictSlash(true)
	menuRouter.Use(middlewares.JWTAuthentication)
	menuRouter.HandleFunc("/", controllers.MenuPageHandler).Methods(http.MethodGet)
	menuRouter.HandleFunc("/{id}", controllers.DishPageHandler).Methods(http.MethodGet)
}
