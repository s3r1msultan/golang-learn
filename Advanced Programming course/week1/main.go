package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
)


func sendRequest(url string, payload map[string]interface{}) (string, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseResponse(response string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
    url := "https://jsonplaceholder.typicode.com/posts"
	payload := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	// Send Request
	response, err := sendRequest(url, payload)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Println("Response:", response)

	// Parse Response
	parsedData, err := parseResponse(response)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	// Display Results
	fmt.Println("Parsed Data:")
	for key, value := range parsedData {
		fmt.Printf("%s: %v\n", key, value)
	}
}