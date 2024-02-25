package controllers

import (
	"final/models"

	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := initTemplates()
	headData := models.HeadData{HeadTitle: "Home", StyleName: "Home"}
	headerData := models.HeaderData{CurrentSite: "Home", ProfileID: User.ObjectId.Hex()}
	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
	}
	err := tmpl.ExecuteTemplate(w, "Home.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
