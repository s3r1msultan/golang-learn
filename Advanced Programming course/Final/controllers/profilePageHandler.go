package controllers

import (
	"context"
	"final/db"
	"final/initializers"
	"final/middlewares"
	"final/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := initTemplates()
	headerData := models.HeaderData{
		CurrentSite: "Profile",
	}

	headData := models.HeadData{
		HeadTitle: "Profile",
		StyleName: "Profile",
	}

	objectId, err := middlewares.ParseObjectIdFromJWT(r)
	if err == nil {
		headerData.ProfileID = objectId.Hex()
	}

	err = db.GetUsersCollection().FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&User)
	if err != nil {
		initializers.LogError("finding the user", err, nil)
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
		User:       User,
	}

	err = tmpl.ExecuteTemplate(w, "Profile.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
