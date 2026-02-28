package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// STTClient implements bot.STTService using the OpenAI-compatible Whisper API.
type STTClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewSTTClient creates a new speech-to-text client using the OpenAI Whisper API.
func NewSTTClient(apiKey, baseURL, model string) *STTClient {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if model == "" {
		model = "whisper-1"
	}
	return &STTClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

type transcriptionResponse struct {
	Text  string `json:"text"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *STTClient) Transcribe(ctx context.Context, audioData []byte) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("model", c.model); err != nil {
		return "", fmt.Errorf("write model field: %w", err)
	}

	part, err := writer.CreateFormFile("file", "audio.wav")
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(audioData); err != nil {
		return "", fmt.Errorf("write audio data: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	endpoint := c.baseURL + "/audio/transcriptions"
	log.Printf("STT request: audio_size=%d bytes, model=%s, endpoint=%s", len(audioData), c.model, endpoint)
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &buf)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
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
		return "", fmt.Errorf("STT API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var transResp transcriptionResponse
	if err := json.Unmarshal(respBody, &transResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if transResp.Error != nil {
		return "", fmt.Errorf("STT API error: %s", transResp.Error.Message)
	}

	log.Printf("STT response: duration=%s, text_length=%d", time.Since(start), len(transResp.Text))
	return transResp.Text, nil
}
