package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type ChatSession struct {
	UserMessage     string `json:"user_message"`
	ChatGPTResponse string `json:"chatgpt_response"`
}

var ChatHistory []ChatSession

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
	} `json:"choices"`
}

type ChatRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type ChatResponse struct {
	Message string `json:"message"`
}

func loadEnvVariables() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("public"))
	router.POST("/chat", handleChat)
	router.GET("/", homeHandler)
	router.GET("/history", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ChatHistory); err != nil {
			http.Error(w, "Failed to encode history", http.StatusInternalServerError)
		}
	})
	err := loadEnvVariables()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}

func homeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "public/index.html")
}

func handleChat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req ChatRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(body, &req)

	responseMessage := callOpenAI(req.Topic, req.Message)

	ChatHistory = append(ChatHistory, ChatSession{
		UserMessage:     req.Message,
		ChatGPTResponse: responseMessage,
	})

	resp := ChatResponse{Message: responseMessage}
	responseData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

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

func isValidRequest(message string) bool {
	validTopics := []string{"law", "regulations", "statutes", "legal", "compliance"}
	for _, topic := range validTopics {
		if strings.Contains(strings.ToLower(message), topic) {
			return true
		}
	}
	return false
}
