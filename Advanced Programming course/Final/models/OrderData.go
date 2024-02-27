package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderData struct {
	OrderId     primitive.ObjectID `bson:"_id" json:"_id"`
	TotalPrice  float64            `bson:"total_price" json:"total_price"`
	Dishes      []DishData         `bson:"dishes" json:"dishes"`
	OrderedDate time.Time          `bson:"ordered_date" json:"ordered_date"`
}
