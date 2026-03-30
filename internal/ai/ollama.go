package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const ollamaBaseURL = "http://localhost:11434"

// OllamaProvider implements the Provider interface for local Ollama models.
type OllamaProvider struct {
	client  *http.Client
	baseURL string
	model   string
}

// NewOllamaProvider creates an Ollama provider from the given configuration.
func NewOllamaProvider(config ProviderConfig) (*OllamaProvider, error) {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = ollamaBaseURL
	}

	model := config.Model
	if model == "" {
		model = "llama3"
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 60
	}

	return &OllamaProvider{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		baseURL: baseURL,
		model:   model,
	}, nil
}

// Complete sends a completion request to the local Ollama API.
func (p *OllamaProvider) Complete(ctx context.Context, req Request) (Response, error) {
	payload := map[string]interface{}{
		"model":  p.model,
		"system": req.SystemPrompt,
		"prompt": req.UserPrompt,
		"stream": false,
		"options": map[string]interface{}{
			"num_predict": req.MaxTokens,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return Response{}, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

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
			Provider: "ollama",
			Type:     getErrorType(resp.StatusCode),
			Message:  string(respBody),
		}
	}

	var completion struct {
		Response   string `json:"response"`
		Model      string `json:"model"`
		TotalDur   int64  `json:"total_duration"`
		EvalCount  int    `json:"eval_count"`
		PromptEval int    `json:"prompt_eval_count"`
	}

	if err := json.Unmarshal(respBody, &completion); err != nil {
		return Response{}, fmt.Errorf("unmarshal response: %w", err)
	}

	promptTokens := completion.PromptEval
	if promptTokens == 0 {
		promptTokens = EstimateTokens(req.UserPrompt)
	}

	return Response{
		Content: completion.Response,
		Usage: UsageStats{
			PromptTokens:     promptTokens,
			CompletionTokens: completion.EvalCount,
			TotalTokens:      promptTokens + completion.EvalCount,
		},
		Model: completion.Model,
	}, nil
}

func (p *OllamaProvider) Name() string  { return "ollama" }
func (p *OllamaProvider) Model() string { return p.model }

// CheckHealth verifies the Ollama server is running and responsive.
func (p *OllamaProvider) CheckHealth() error {
	resp, err := p.client.Get(p.baseURL + "/api/tags")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned %d", resp.StatusCode)
	}

	return nil
}
