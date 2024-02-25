package controllers

import (
	"final/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func MenuPageHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := initTemplates()
	headData := models.HeadData{
		HeadTitle: "Menu",
		StyleName: "Menu",
	}

	headerData := models.HeaderData{CurrentSite: "Menu"}

	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
		//	Dishes
		// TODO create dishes data
	}

	err := tmpl.ExecuteTemplate(w, "Menu.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
