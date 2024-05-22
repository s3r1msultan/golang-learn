package main

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
