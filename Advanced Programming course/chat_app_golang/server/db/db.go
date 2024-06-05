package db

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func getURI() string {
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		log.Fatal("MongoDB URI is not appropriate. Please provide correct one. ")
	}
	return URI
}

func Connect() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(getURI()).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	MongoClient = client

	return err
}

func GetDB() string {
	DB := os.Getenv("DATABASE")
	return DB
}

func Disconnect() {
	MongoClient.Disconnect(context.TODO())
}
