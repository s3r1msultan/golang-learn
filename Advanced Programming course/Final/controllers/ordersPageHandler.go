package controllers

import (
	"final/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func OrdersPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tmpl := initTemplates()
	headData := models.HeadData{
		HeadTitle: "Orders",
		StyleName: "Orders",
	}

	headerData := models.HeaderData{
		CurrentSite: "Orders",
		ProfileID:   id,
	}

	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
		User:       User,
		//	Dishes
		// TODO create dishes data
	}
	err := tmpl.ExecuteTemplate(w, "Orders.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
