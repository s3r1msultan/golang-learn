package models

type CartData struct {
	Dishes []DishData `bson:"dishes" json:"dishes"`
}
