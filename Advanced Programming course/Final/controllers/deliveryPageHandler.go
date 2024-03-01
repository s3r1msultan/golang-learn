package controllers

import (
	"context"
	"encoding/json"
	"final/db"
	"final/middlewares"
	"final/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func DeliveryPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		id := vars["id"]
		tmpl := initTemplates()
		headerData := models.HeaderData{
			CurrentSite: "Delivery",
			ProfileID:   id,
		}

		headData := models.HeadData{
			HeadTitle: "Delivery",
			StyleName: "Delivery",
		}

		objectId, err := middlewares.ParseObjectIdFromJWT(r)
		if err == nil {
			headerData.ProfileID = objectId.Hex()
		}

		data := models.PageData{
			HeaderData: headerData,
			HeadData:   headData,
			User:       User,
		}

		err = tmpl.ExecuteTemplate(w, "Delivery.html", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
	}

	if r.Method == http.MethodPut {
		vars := mux.Vars(r)
		id := vars["id"]
		var delivery models.DeliveryData
		err := json.NewDecoder(r.Body).Decode(&delivery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		var result models.User
		err = db.GetUsersCollection().
			FindOneAndUpdate(
				context.TODO(),
				bson.M{"_id": objectId},
				bson.M{"$set": bson.M{
					"delivery.full_name":    delivery.FullName,
					"delivery.address":      delivery.Address,
					"delivery.city":         delivery.City,
					"delivery.zip_code":     delivery.ZipCode,
					"delivery.phone_number": delivery.PhoneNumber,
				}},
			).Decode(&result)

		if err != nil {
			log.Fatal(err)
		}
	}

}
