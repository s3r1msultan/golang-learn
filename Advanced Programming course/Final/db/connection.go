package db

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func Connect() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(getURI())
	client, err := mongo.Connect(context.TODO(), clientOptions)
	return client, err
}

func getURI() string {
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		log.Fatal("MongoDB URI is not appropriate. Please provide correct one. ")
	}
	return URI
}
