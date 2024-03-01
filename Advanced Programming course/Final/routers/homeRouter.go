package routers

import (
	"final/controllers"
	"final/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func HomeRouter(authRouter *mux.Router) {
	authRouter.Use(middlewares.FormToJSONMiddleware)
	authRouter.StrictSlash(true)
	authRouter.HandleFunc("/", controllers.HomePageHandler).Methods(http.MethodGet)
}
