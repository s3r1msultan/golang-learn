package controllers

import (
	"context"
	"encoding/json"
	"final/db"
	"final/middlewares"
	"final/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type dishId struct {
	DishId string `json:"dishId"`
}

func GetDishes(w http.ResponseWriter, r *http.Request) {
	tmpl := initTemplates()
	headerData := models.HeaderData{
		CurrentSite: "Cart",
	}

	headData := models.HeadData{
		HeadTitle: "Cart",
		StyleName: "Cart",
	}

	objectId, err := middlewares.ParseObjectIdFromJWT(r)
	if err == nil {
		headerData.ProfileID = objectId.Hex()
	}

	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
		User:       User,
		Dishes:     User.Cart,
	}

	err = tmpl.ExecuteTemplate(w, "Cart.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func AddDish(w http.ResponseWriter, r *http.Request) {
	var req dishId
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	dishObjectId, err := primitive.ObjectIDFromHex(req.DishId)
	if err != nil {
		http.Error(w, "Invalid dish ID format", http.StatusBadRequest)
		return
	}

	objectId, err := middlewares.ParseObjectIdFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: unable to parse user ID from JWT", http.StatusUnauthorized)
		return
	}
	var dish models.DishData
	err = db.GetDishesCollection().FindOne(context.TODO(), bson.M{"_id": dishObjectId}).Decode(&dish)

	_, err = db.GetUsersCollection().
		UpdateOne(
			context.TODO(),
			bson.M{"_id": objectId},
			bson.M{"$addToSet": bson.M{"cart": dish}}, // Note: Adjust this if your cart structure is different
		)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dish added to cart successfully"})
}
func DeleteDish(w http.ResponseWriter, r *http.Request) {
	var req dishId
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dishObjectId, err := primitive.ObjectIDFromHex(req.DishId)
	if err != nil {
		http.Error(w, "Invalid dish ID", http.StatusBadRequest)
		return
	}

	userId, err := middlewares.ParseObjectIdFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: unable to parse user ID from JWT", http.StatusUnauthorized)
		return
	}

	_, err = db.GetUsersCollection().UpdateOne(
		context.TODO(),
		bson.M{"_id": userId},
		bson.M{"$pull": bson.M{"cart": bson.M{"_id": dishObjectId}}},
	)
	if err != nil {
		http.Error(w, "Failed to remove dish from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dish removed from cart successfully"})
}
