package routers

import (
	"final/controllers"
	"final/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func ProfileRouter(profileRouter *mux.Router) {
	profileRouter.Use(middlewares.JWTAuthentication)
	profileRouter.HandleFunc("/", controllers.ProfilePageHandler)
	profileRouter.HandleFunc("/orders", controllers.OrdersPageHandler).Methods(http.MethodGet)
	profileRouter.HandleFunc("/delivery", controllers.DeliveryPageHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPut)
	CartRouter(profileRouter.PathPrefix("/cart").Subrouter())

}
