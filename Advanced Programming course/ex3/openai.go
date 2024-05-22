package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func callOpenAI(topic, userMessage string) string {
	if !isValidRequest(topic) {
		return "Your request was declined because your question is not related to the vision of the jurisdiction company."
	}

	systemMessage := fmt.Sprintf("You are a helpful assistant. You need to provide information on the topic: %s. If questions are not related to this topic, then say that you can't answer on this question", topic)
	data := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "system", Content: systemMessage},
			{Role: "user", Content: userMessage},
		},
	}

	dataJSON, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(dataJSON))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to OpenAI API: %s", err)
		return "Error communicating with OpenAI API"
	}
	defer resp.Body.Close()

	var apiResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("Error decoding response from OpenAI API: %s", err)
		return "Failed to decode response from OpenAI API"
	}

	if len(apiResponse.Choices) > 0 {
		return apiResponse.Choices[0].Message.Content
	}

	return "No response received"
}
