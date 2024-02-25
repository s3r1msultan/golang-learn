package models

type DeliveryData struct {
	FullName    string `bson:"full_name" json:"full_name"`
	Address     string `bson:"address" json:"address"`
	City        string `bson:"city" json:"city"`
	ZipCode     string `bson:"zip_code" json:"zip_code"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
}
