package routers

import (
	"final/controllers"
	"github.com/gorilla/mux"
)

func VerificationRouter(router *mux.Router) {
	router.HandleFunc("", controllers.VerificationHandler)
}
