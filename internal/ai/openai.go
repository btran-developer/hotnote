package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const openAIBaseURL = "https://api.openai.com/v1"

// OpenAIProvider implements the Provider interface for OpenAI's Chat Completions API.
type OpenAIProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

// NewOpenAIProvider creates an OpenAI provider from the given configuration.
func NewOpenAIProvider(config ProviderConfig) (*OpenAIProvider, error) {
	apiKey := os.Getenv(config.APIKeyEnv) // Resolve env var name to value
	if apiKey == "" {
		return nil, ErrAIAPIKeyMissing
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = openAIBaseURL
	}

	model := config.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 60
	}

	return &OpenAIProvider{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
	}, nil
}

// Complete sends a completion request to the OpenAI API.
func (p *OpenAIProvider) Complete(ctx context.Context, req Request) (Response, error) {
	payload := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": req.SystemPrompt},
			{"role": "user", "content": req.UserPrompt},
		},
		"max_tokens": req.MaxTokens,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Response{}, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Response{}, &ProviderError{
			Provider: "openai",
			Type:     getErrorType(resp.StatusCode),
			Message:  string(respBody),
		}
	}

	var completion struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(respBody, &completion); err != nil {
		return Response{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(completion.Choices) == 0 {
		return Response{}, ErrAIResponseInvalid
	}

	return Response{
		Content: completion.Choices[0].Message.Content,
		Usage: UsageStats{
			PromptTokens:     completion.Usage.PromptTokens,
			CompletionTokens: completion.Usage.CompletionTokens,
			TotalTokens:      completion.Usage.TotalTokens,
		},
		Model: completion.Model,
	}, nil
}

func (p *OpenAIProvider) Name() string  { return "openai" }
func (p *OpenAIProvider) Model() string { return p.model }

func getErrorType(statusCode int) string {
	switch statusCode {
	case 401:
		return ErrTypeAuth
	case 429:
		return ErrTypeRateLimit
	case 400:
		return ErrTypeInvalidReq
	default:
		if statusCode >= 500 {
			return ErrTypeServerError
		}
		return "unknown"
	}
}
