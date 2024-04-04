package controllers

import (
	"context"
	"final/db"
	"final/middlewares"
	"final/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
)

func MenuPageHandler(w http.ResponseWriter, r *http.Request) {
	minPriceParam := r.URL.Query().Get("minPrice")
	maxPriceParam := r.URL.Query().Get("maxPrice")
	sortParam := r.URL.Query().Get("sort")
	pageParam := r.URL.Query().Get("page")
	pageSizeParam := r.URL.Query().Get("pageSize")
	minPrice, _ := strconv.ParseFloat(minPriceParam, 64)
	maxPrice, _ := strconv.ParseFloat(maxPriceParam, 64)
	page, _ := strconv.Atoi(pageParam)
	pageSize, _ := strconv.Atoi(pageSizeParam)

	if pageSize <= 0 {
		pageSize = 15
	}
	if page <= 0 {
		page = 1
	}

	var filter bson.M = bson.M{}
	if minPriceParam != "" {
		filter["price"] = bson.M{"$gte": minPrice}
	}
	if maxPriceParam != "" {
		filter["price"] = bson.M{"$lte": maxPrice}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))

	if sortParam == "asc" {
		findOptions.SetSort(bson.D{{"price", 1}})
	} else if sortParam == "desc" {
		findOptions.SetSort(bson.D{{"price", -1}})
	}

	collection := db.GetDishesCollection()
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var dishes []models.DishData
	if err = cursor.All(ctx, &dishes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := initTemplates()
	headData := models.HeadData{
		HeadTitle: "Menu",
		StyleName: "Menu",
	}

	headerData := models.HeaderData{CurrentSite: "Menu"}
	objectId, err := middlewares.ParseObjectIdFromJWT(r)
	if err == nil {
		headerData.ProfileID = objectId.Hex()
	}

	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
		User:       User,
		Dishes:     dishes,
	}

	err = tmpl.ExecuteTemplate(w, "Menu.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}
