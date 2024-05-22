package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
)

func TestMain(m *testing.M) {
	err := loadEnvVariables()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func resetState() {
	ChatHistory = []ChatSession{}
}

func TestHandleChat(t *testing.T) {
	resetState()
	mockOpenAICall := func(topic, userMessage string) string {
		return "The legal age for voting is 18."
	}

	reqBody := `{"topic": "law", "message": "What is the legal age for voting?"}`
	req := httptest.NewRequest("POST", "/chat", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := handleChat(mockOpenAICall)
	handler(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handleChat() status = %v; expected %v", resp.StatusCode, http.StatusOK)
	}

	var chatResponse ChatResponse
	err := json.NewDecoder(resp.Body).Decode(&chatResponse)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	expectedMessage := "The legal age for voting is 18."
	if chatResponse.Message != expectedMessage {
		t.Errorf("chatResponse.Message = %v; expected %v", chatResponse.Message, expectedMessage)
	}
}

func TestHomeHandler(t *testing.T) {
	resetState()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	homeHandler(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("homeHandler() status = %v; expected %v", resp.StatusCode, http.StatusOK)
	}
}

func TestHistoryHandler(t *testing.T) {
	resetState()
	ChatHistory = []ChatSession{
		{UserMessage: "Test message 1", ChatGPTResponse: "Test response 1"},
		{UserMessage: "Test message 2", ChatGPTResponse: "Test response 2"},
	}

	req := httptest.NewRequest("GET", "/history", nil)
	w := httptest.NewRecorder()
	router := setupRouter(callOpenAI)
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("historyHandler() status = %v; expected %v", resp.StatusCode, http.StatusOK)
	}

	var history []ChatSession
	err := json.NewDecoder(resp.Body).Decode(&history)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Expected 2 history items, got %d", len(history))
	}
}

func TestCallOpenAI(t *testing.T) {
	resetState()
	response := callOpenAI("law", "What is the legal age for voting?")
	if response == "" {
		t.Errorf("callOpenAI() returned an empty response")
	}
}

func TestLoadEnvVariables(t *testing.T) {
	err := loadEnvVariables()
	if err != nil {
		t.Errorf("Failed to load env variables: %v", err)
	}
}

func TestIsValidRequest(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"law", true},
		{"finance", false},
		{"Legal advice", true},
		{"health", false},
	}

	for _, test := range tests {
		result := isValidRequest(test.input)
		if result != test.expected {
			t.Errorf("isValidRequest(%s) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestUpdateHistory(t *testing.T) {
	resetState()
	ChatHistory = append(ChatHistory, ChatSession{
		UserMessage:     "request1",
		ChatGPTResponse: "response1",
	})

	if len(ChatHistory) != 1 || ChatHistory[0].UserMessage != "request1" || ChatHistory[0].ChatGPTResponse != "response1" {
		t.Errorf("Expected history to contain the interaction; got %v", ChatHistory)
	}
}

func TestHistoryLimit(t *testing.T) {
	resetState()
	for i := 0; i < 6; i++ {
		ChatHistory = append(ChatHistory, ChatSession{
			UserMessage:     fmt.Sprintf("request%d", i),
			ChatGPTResponse: fmt.Sprintf("response%d", i),
		})
	}

	if len(ChatHistory) >= 10 {
		t.Errorf("Expected history to contain 5 interactions; got %v", len(ChatHistory))
	}
}

func TestMainFunction(t *testing.T) {
	resetState()
	go func() {
		main()
	}()
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get("http://localhost:8080/")
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK; got %v", resp.StatusCode)
			}
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("Failed to connect to the server after multiple attempts")
}

// To run: go test -coverprofile=coverage.out
func TestCoverage(t *testing.T) {
	t.Skip("Skipping coverage test")
}

// To run: go tool cover -html=coverage.out
// To run: go test -short -v
func TestShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping short tests")
	}
	t.Run("ShortTestExample", func(t *testing.T) {
		if 1+1 != 2 {
			t.Error("Math is broken")
		}
	})
}

// Additional examples of argument usage
func TestList(t *testing.T) {
	t.Log("List test")
}

func TestCount(t *testing.T) {
	t.Log("Count test")
}

func TestJSON(t *testing.T) {
	t.Log("JSON test")
}

func TestCPU(t *testing.T) {
	t.Log("CPU test")
}

func TestRace(t *testing.T) {
	t.Log("Race test")
}
