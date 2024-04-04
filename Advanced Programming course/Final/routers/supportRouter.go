package routers

import (
	"final/controllers"
	"github.com/gorilla/mux"
)

func SupportRouter(r *mux.Router) {
	r.HandleFunc("/", controllers.SupportPageHandler)
}
