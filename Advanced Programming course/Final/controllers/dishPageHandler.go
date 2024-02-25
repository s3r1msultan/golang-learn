package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func DishPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars["id"])
}
