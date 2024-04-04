package routers

import (
	"final/controllers"
	"github.com/gorilla/mux"
)

func AdminRouter(adminRouter *mux.Router) {
	adminRouter.HandleFunc("/", controllers.AdminPageHandler)
	adminRouter.HandleFunc("/api/users/{id}", controllers.GetUser).Methods("GET")
	adminRouter.HandleFunc("/api/users", controllers.ListUsers).Methods("GET")
	adminRouter.HandleFunc("/api/users", controllers.AddUser).Methods("POST")
	adminRouter.HandleFunc("/api/users/{id}", controllers.UpdateUser).Methods("PUT")
	adminRouter.HandleFunc("/api/users/{id}", controllers.DeleteUser).Methods("DELETE")
}
