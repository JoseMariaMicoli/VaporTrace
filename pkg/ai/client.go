package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LLMProvider defines the contract for AI backends
type LLMProvider interface {
	Configure(apiKey, model, endpoint string)
	Analyze(input string) (string, error)
	GeneratePayloads(context string, count int) ([]string, error)
}

// --- OLLAMA PROVIDER ---

type OllamaClient struct {
	Endpoint string
	Model    string
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func (o *OllamaClient) Configure(apiKey, model, endpoint string) {
	// Robust Endpoint Construction
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}
	// Strip trailing slash
	endpoint = strings.TrimSuffix(endpoint, "/")

	// If the user didn't provide the full API path, append it
	if !strings.HasSuffix(endpoint, "/api/generate") {
		endpoint = endpoint + "/api/generate"
	}

	o.Endpoint = endpoint
	if model == "" {
		model = "mistral"
	}
	o.Model = model
}

func (o *OllamaClient) Analyze(input string) (string, error) {
	reqBody := OllamaRequest{
		Model:  o.Model,
		Prompt: input,
		System: SystemPersona,
		Stream: false,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(o.Endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error (%d): %s - (Hint: Try 'ollama pull %s')", resp.StatusCode, string(body), o.Model)
	}

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Response, nil
}

func (o *OllamaClient) GeneratePayloads(context string, count int) ([]string, error) {
	// Prompt engineering for specific attack payloads
	prompt := fmt.Sprintf("Generate %d unique, aggressive fuzzing payloads for: %s. Return only the payloads, one per line. No markdown, no explanations.", count, context)
	resp, err := o.Analyze(prompt)
	if err != nil {
		return nil, err
	}
	// Clean output
	lines := strings.Split(resp, "\n")
	var clean []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "`") {
			clean = append(clean, l)
		}
	}
	return clean, nil
}

// --- OPENAI PROVIDER ---

type OpenAIClient struct {
	APIKey   string
	Endpoint string
	Model    string
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (o *OpenAIClient) Configure(apiKey, model, endpoint string) {
	o.APIKey = apiKey
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}
	o.Endpoint = endpoint
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	o.Model = model
}

func (o *OpenAIClient) Analyze(input string) (string, error) {
	reqBody := OpenAIRequest{
		Model: o.Model,
		Messages: []Message{
			{Role: "system", Content: SystemPersona},
			{Role: "user", Content: input},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", o.Endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai error (%d): %s", resp.StatusCode, string(body))
	}

	var result OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("empty response")
}

func (o *OpenAIClient) GeneratePayloads(context string, count int) ([]string, error) {
	prompt := fmt.Sprintf("Generate %d attack payloads for %s. Raw text only.", count, context)
	resp, err := o.Analyze(prompt)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(resp), "\n"), nil
}
