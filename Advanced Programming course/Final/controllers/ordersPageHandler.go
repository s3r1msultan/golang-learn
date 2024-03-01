package controllers

import (
	"final/middlewares"
	"final/models"
	"fmt"
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

	objectId, err := middlewares.ParseObjectIdFromJWT(r)
	if err == nil {
		headerData.ProfileID = objectId.Hex()
	}

	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
		User:       User,
		//	Dishes
		// TODO create dishes data
	}
	err = tmpl.ExecuteTemplate(w, "Orders.html", data)
	fmt.Println(User.Orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
