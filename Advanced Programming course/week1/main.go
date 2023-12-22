package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type SignUpRequest struct {
	FirstName string `json:"first_name"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type JsonRequest struct {
	Message string `json:"message"`
}

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const PORT string = ":8080"

func main() {
	http.HandleFunc("/", handlePostRequest)
	http.HandleFunc("/signup", handleSignup)
	fmt.Printf("Server listening on port %s...\n", PORT)
	http.ListenAndServe(PORT, nil)
}

func handlePostRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var requestData JsonRequest

	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if requestData.Message == "" {
		http.Error(res, "Invalid JSON message", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received message: %s\n", requestData.Message)

	response := JsonResponse{
		Status:  "success",
		Message: "Data successfully received",
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}

func handleSignup(res http.ResponseWriter, req *http.Request) {
	if http.MethodPost != req.Method {
		http.Error(res, "Method not allowed", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var requestData SignUpRequest
	err := decoder.Decode(&requestData)

	if err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	requestDataType := reflect.TypeOf(requestData)
	for i := 0; i < requestDataType.NumField(); i++ {
		if reflect.ValueOf(requestData).Field(i).Interface() == "" {
			http.Error(res, "Invalid JSON data, some fields are empty", http.StatusBadRequest)
			return
		}
	}

	for i := 0; i < requestDataType.NumField(); i++ {
		fieldName := requestDataType.Field(i).Name
		fieldValue := reflect.ValueOf(requestData).Field(i).Interface()
		fmt.Printf("%s : %v \n", fieldName, fieldValue)
	}

	response := JsonResponse{
		Status:  "OK",
		Message: "Data successfully received",
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}
