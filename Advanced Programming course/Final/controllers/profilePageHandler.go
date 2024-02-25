package controllers

import (
	"final/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tmpl := initTemplates()
	headerData := models.HeaderData{
		CurrentSite: "Profile",
		ProfileID:   id,
	}

	headData := models.HeadData{
		HeadTitle: "Profile",
		StyleName: "Profile",
	}

	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
		User:       User,
	}

	err := tmpl.ExecuteTemplate(w, "Profile.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
