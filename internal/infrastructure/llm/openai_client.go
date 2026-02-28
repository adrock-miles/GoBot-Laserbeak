package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/adrock-miles/go-laserbeak/internal/domain/bot"
)

// OpenAIClient implements bot.LLMService using the OpenAI-compatible chat completions API.
type OpenAIClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewOpenAIClient creates a new OpenAI-compatible LLM client.
func NewOpenAIClient(apiKey, baseURL, model string) *OpenAIClient {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

type chatRequest struct {
	Model    string       `json:"model"`
	Messages []chatMsg    `json:"messages"`
}

type chatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *OpenAIClient) ChatCompletion(ctx context.Context, messages []bot.LLMMessage) (string, error) {
	msgs := make([]chatMsg, len(messages))
	for i, m := range messages {
		msgs[i] = chatMsg{Role: m.Role, Content: m.Content}
	}

	reqBody := chatRequest{
		Model:    c.model,
		Messages: msgs,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint := c.baseURL + "/chat/completions"
	log.Printf("LLM request: messages=%d, model=%s, endpoint=%s", len(msgs), c.model, endpoint)
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	result := chatResp.Choices[0].Message.Content
	log.Printf("LLM response: duration=%s, length=%d", time.Since(start), len(result))
	return result, nil
}
