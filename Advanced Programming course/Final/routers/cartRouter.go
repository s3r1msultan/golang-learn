package routers

import (
	"final/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func CartRouter(cartRouter *mux.Router) {
	cartRouter.StrictSlash(true)
	cartRouter.HandleFunc("/", controllers.GetDishes).Methods(http.MethodGet)
	cartRouter.HandleFunc("/add", controllers.AddDish).Methods(http.MethodPut)
	cartRouter.HandleFunc("/delete", controllers.DeleteDish).Methods(http.MethodDelete)
}
