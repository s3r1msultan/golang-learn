package controllers

import (
	"final/middlewares"
	"final/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SupportPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := initTemplates()
	headerData := models.HeaderData{CurrentSite: "Support", ProfileID: User.ObjectId.Hex()}
	headData := models.HeadData{
		HeadTitle: "Support",
		StyleName: "Support",
	}
	objectId, err := middlewares.ParseObjectIdFromJWT(r)
	if err == nil {
		headerData.ProfileID = objectId.Hex()
	}
	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
		//	Dishes
		// TODO create dishes data
	}

	err = tmpl.ExecuteTemplate(w, "Support.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
