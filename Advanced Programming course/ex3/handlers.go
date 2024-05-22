package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OpenAICaller func(topic, userMessage string) string

func handleChat(openAICall OpenAICaller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		responseMessage := openAICall(req.Topic, req.Message)

		ChatHistory = append(ChatHistory, ChatSession{
			UserMessage:     req.Message,
			ChatGPTResponse: responseMessage,
		})

		resp := ChatResponse{Message: responseMessage}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Error processing response", http.StatusInternalServerError)
			return
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "public/index.html")
}

func historyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ChatHistory); err != nil {
		http.Error(w, "Failed to encode history", http.StatusInternalServerError)
	}
}
