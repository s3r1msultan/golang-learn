package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ObjectId          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName         string             `bson:"first_name" json:"first_name"`
	LastName          string             `bson:"last_name" json:"last_name"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"password"`
	PhoneNumber       string             `bson:"phone_number" json:"phone_number"`
	VerificationToken string             `bson:"verification_token"`
	EmailVerified     bool               `bson:"email_verified"`
	Orders            []OrderData        `bson:"orders" json:"orders"`
	Cart              []DishData         `bson:"cart" json:"cart"`
	Delivery          DeliveryData       `bson:"delivery" json:"delivery"`
	IsAdmin           bool               `bson:"is_admin" json:"isAdmin"`
}
