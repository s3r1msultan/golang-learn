package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DishData struct {
	ObjectID      primitive.ObjectID `bson:"_id" json:"_id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	Price         float32            `bson:"price" json:"price"`
	Weight        float32            `bson:"weight" json:"weight"`
	Proteins      float32            `bson:"protein" json:"protein"`
	Fats          float32            `bson:"fats" json:"fats"`
	Carbohydrates float32            `bson:"carbohydrates" json:"carbohydrates"`
	ImgURL        string             `bson:"img_URL" json:"img_URL"`
}
