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

const anthropicBaseURL = "https://api.anthropic.com/v1"

// AnthropicProvider implements the Provider interface for Anthropic's Messages API.
type AnthropicProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

// NewAnthropicProvider creates an Anthropic provider from the given configuration.
func NewAnthropicProvider(config ProviderConfig) (*AnthropicProvider, error) {
	apiKey := os.Getenv(config.APIKeyEnv) // Resolve env var name to value
	if apiKey == "" {
		return nil, ErrAIAPIKeyMissing
	}

	model := config.Model
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = anthropicBaseURL
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 60
	}

	return &AnthropicProvider{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
	}, nil
}

// Complete sends a completion request to the Anthropic API.
func (p *AnthropicProvider) Complete(ctx context.Context, req Request) (Response, error) {
	payload := map[string]interface{}{
		"model":  p.model,
		"system": req.SystemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": req.UserPrompt},
		},
		"max_tokens": req.MaxTokens,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return Response{}, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

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
			Provider: "anthropic",
			Type:     getErrorType(resp.StatusCode),
			Message:  string(respBody),
		}
	}

	var completion struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(respBody, &completion); err != nil {
		return Response{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(completion.Content) == 0 || completion.Content[0].Type != "text" {
		return Response{}, ErrAIResponseInvalid
	}

	return Response{
		Content: completion.Content[0].Text,
		Usage: UsageStats{
			PromptTokens:     completion.Usage.InputTokens,
			CompletionTokens: completion.Usage.OutputTokens,
			TotalTokens:      completion.Usage.InputTokens + completion.Usage.OutputTokens,
		},
		Model: completion.Model,
	}, nil
}

func (p *AnthropicProvider) Name() string  { return "anthropic" }
func (p *AnthropicProvider) Model() string { return p.model }
