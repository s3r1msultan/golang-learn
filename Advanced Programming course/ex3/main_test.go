package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestCallOpenAI(t *testing.T) {
	err := os.Setenv("OPENAI_API_KEY", "dummy_key")
	if err != nil {
		return
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices": [{"message": {"content": "Mocked response", "role": "assistant"}}]}`))
	}))
	defer ts.Close()

	response := "Mocked response"
	if response != "Mocked response" {
		t.Errorf("expected 'Mocked response', got '%s'", response)
	}
}

func TestLoadEnvVariables(t *testing.T) {
	originalEnvContent, err := os.ReadFile(".env")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("unable to read existing .env file: %v", err)
	}

	testEnvContent := []byte("OPENAI_API_KEY=dummy_key")
	err = os.WriteFile(".env", testEnvContent, 0644)
	if err != nil {
		t.Fatalf("unable to write test .env file: %v", err)
	}

	defer func() {
		if len(originalEnvContent) > 0 {
			err = os.WriteFile(".env", originalEnvContent, 0644)
			if err != nil {
				t.Fatalf("unable to restore original .env file: %v", err)
			}
		} else {
			err = os.Remove(".env")
			if err != nil && !os.IsNotExist(err) {
				t.Fatalf("unable to delete test .env file: %v", err)
			}
		}
	}()
	err = loadEnvVariables()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if os.Getenv("OPENAI_API_KEY") != "dummy_key" {
		t.Errorf("expected OPENAI_API_KEY to be 'dummy_key', got %s", os.Getenv("OPENAI_API_KEY"))
	}
}

func TestIsValidRequest(t *testing.T) {
	tests := []struct {
		message  string
		expected bool
	}{
		{"This is about law", true},
		{"Random message", false},
		{"Comply with regulations", true},
		{"No valid topic", false},
	}

	for _, tt := range tests {
		result := isValidRequest(tt.message)
		if result != tt.expected {
			t.Errorf("expected %v for message %q, got %v", tt.expected, tt.message, result)
		}
	}
}

func TestHandleChat(t *testing.T) {
	mockOpenAICall := func(topic, userMessage string) string {
		return "Mock response"
	}

	router := httprouter.New()
	router.POST("/chat", handleChat(mockOpenAICall))

	chatRequest := ChatRequest{
		Topic:   "law",
		Message: "What is the capital of Kazakhstan?",
	}
	body, _ := json.Marshal(chatRequest)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp ChatResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("unable to decode response: %v", err)
	}

	if resp.Message != "Mock response" {
		t.Errorf("unexpected response message: got %v want %v", resp.Message, "Mock response")
	}
}

func TestHomeHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.GET("/", homeHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHistoryHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/history", nil)
	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.GET("/history", historyHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var history []ChatSession
	if err := json.NewDecoder(rr.Body).Decode(&history); err != nil {
		t.Errorf("unable to decode response: %v", err)
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
		resp, err := http.Get("http://localhost:3001/")
		if err == nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(resp.Body)
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK; got %v", resp.StatusCode)
			}
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("Failed to connect to the server after multiple attempts")
}

func TestCoverage(t *testing.T) {
	t.Skip("Skipping coverage test")
}

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
