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

// --- OLLAMA PROVIDER (Local Privacy) ---

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
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}
	endpoint = strings.TrimSuffix(endpoint, "/")

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
	prompt := fmt.Sprintf("Generate %d unique, aggressive fuzzing payloads for: %s. Return only the payloads, one per line. No markdown, no explanations.", count, context)
	resp, err := o.Analyze(prompt)
	if err != nil {
		return nil, err
	}
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

// --- OPENAI PROVIDER (High Performance Cloud) ---

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

// --- GOOGLE GEMINI PROVIDER (High Context / Free Tier) ---

type GeminiClient struct {
	APIKey string
	Model  string
	URL    string
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (g *GeminiClient) Configure(apiKey, model, endpoint string) {
	g.APIKey = apiKey

	// Auto-correct common model name issues for v1beta
	// The API often requires specific version aliases (e.g., -latest or -001)
	if model == "" || model == "gemini-1.5-flash" {
		model = "gemini-1.5-flash-latest"
	}
	g.Model = model

	// Gemini uses a different REST structure: url + ?key=API_KEY
	g.URL = fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)
}

func (g *GeminiClient) Analyze(input string) (string, error) {
	// Combine Persona and Input because Gemini v1beta doesn't strongly enforce system roles in basic rest
	fullPrompt := fmt.Sprintf("%s\n\nTASK:\n%s", SystemPersona, input)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: fullPrompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(g.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("gemini conn failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini api error (%d): %s", resp.StatusCode, string(body))
	}

	var result GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return result.Candidates[0].Content.Parts[0].Text, nil
	}
	return "", fmt.Errorf("empty gemini response")
}

func (g *GeminiClient) GeneratePayloads(context string, count int) ([]string, error) {
	prompt := fmt.Sprintf("Generate %d unique attack payloads for %s. Return ONLY raw text lines.", count, context)
	resp, err := g.Analyze(prompt)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(resp, "\n")
	var clean []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		// Clean typical LLM bullet points if present
		l = strings.TrimPrefix(l, "* ")
		l = strings.TrimPrefix(l, "- ")
		if l != "" {
			clean = append(clean, l)
		}
	}
	return clean, nil
}
