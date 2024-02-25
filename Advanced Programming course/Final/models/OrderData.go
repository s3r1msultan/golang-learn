package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderData struct {
	OrderId     primitive.ObjectID `bson:"_id" json:"_id"`
	TotalPrice  int32              `bson:"total_price" json:"total_price"`
	Dishes      []DishData         `bson:"dishes" json:"dishes"`
	OrderedDate primitive.DateTime `bson:"ordered_date" json:"ordered_date"`
}
