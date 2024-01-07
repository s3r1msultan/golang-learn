package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"reflect"
)

type DeliveryRequest struct {
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	ZipCode     uint32 `json:"zip_code"`
	PhoneNumber string `json:"phone_number"`
}

type JSONResponse struct {
	Status  uint16 `json:"status"`
	Message string `json:"message"`
}

const PORT string = ":8080"
const DATABASE string = "go_restaurants"
const COLLECTION string = "restaurants"

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/delivery", handleDelivery)
	fmt.Printf("Server listening on port %s...\n", PORT)
	http.ListenAndServe(PORT, nil)
}

func handleHome(res http.ResponseWriter, req *http.Request) {
	if http.MethodGet != req.Method {
		http.Error(res, "Method not allowed", http.StatusBadRequest)
		return
	}

	http.ServeFile(res, req, "../client/FinalProject/DeliveryAddressPage/Delivery.html")
}

func handleDelivery(res http.ResponseWriter, req *http.Request) {
	if http.MethodPost != req.Method {
		http.Error(res, "Method not allowed", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var requestData DeliveryRequest
	var responseData JSONResponse
	var isSuccessful bool = true
	res.Header().Set("Content-Type", "application/json")

	err := decoder.Decode(&requestData)
	if err != nil {
		responseData = JSONResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON format",
		}
		res.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(res).Encode(responseData)
		if err != nil {
			return
		}
		return
	}

	requestDataType := reflect.TypeOf(requestData)
	var emptyFields string
	for i := 0; i < requestDataType.NumField(); i++ {
		if reflect.ValueOf(requestData).Field(i).Interface() == "" {
			emptyFields += requestDataType.Field(i).Name + " "
		}
	}

	if emptyFields != "" {
		responseData = JSONResponse{
			Status:  http.StatusBadRequest,
			Message: emptyFields + "fields are empty",
		}

		isSuccessful = false
	}

	if isSuccessful {
		err := insertDataToMongoDB(requestData)
		if err != nil {
			responseData = JSONResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error inserting data into MongoDB",
			}
			res.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(res).Encode(responseData)
			if err != nil {
				return
			}
			return
		}

		responseData = JSONResponse{
			Status:  http.StatusOK,
			Message: "Data successfully received and inserted into MongoDB",
		}
		res.WriteHeader(http.StatusOK)
	}

	err = json.NewEncoder(res).Encode(responseData)
	if err != nil {
		return
	}
}

func insertDataToMongoDB(data DeliveryRequest) error {
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:root@nosqlcourse.eyuhz6s.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(DATABASE).Collection(COLLECTION)
	_, err = collection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

func createDeliveryAddress(data DeliveryRequest) (*mongo.InsertOneResult, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:root@nosqlcourse.eyuhz6s.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(DATABASE).Collection(COLLECTION)
	result, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getDeliveryAddressByID(id string) (*DeliveryRequest, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:root@nosqlcourse.eyuhz6s.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(DATABASE).Collection(COLLECTION)

	var deliveryAddress DeliveryRequest
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&deliveryAddress)
	if err != nil {
		return nil, err
	}

	return &deliveryAddress, nil
}

func updateDeliveryAddressNameByID(id string, newName string) (*mongo.UpdateResult, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:root@nosqlcourse.eyuhz6s.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(DATABASE).Collection(COLLECTION)

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"full_name", newName}}},
		},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deleteDeliveryAddressByID(id string) (*mongo.DeleteResult, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:root@nosqlcourse.eyuhz6s.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(DATABASE).Collection(COLLECTION)

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getAllDeliveryAddresses() ([]DeliveryRequest, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:root@nosqlcourse.eyuhz6s.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(DATABASE).Collection(COLLECTION)

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var deliveryAddresses []DeliveryRequest
	err = cursor.All(context.Background(), &deliveryAddresses)
	if err != nil {
		return nil, err
	}

	return deliveryAddresses, nil
}
