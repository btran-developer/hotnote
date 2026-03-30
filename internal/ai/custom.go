package ai

import (
	"context"
	"fmt"
)

// CustomProvider wraps an OpenAI-compatible provider with a custom base URL.
type CustomProvider struct {
	openai *OpenAIProvider
}

// NewCustomProvider creates a custom provider from the given configuration.
// The configuration must include a BaseURL for the custom endpoint.
func NewCustomProvider(config ProviderConfig) (*CustomProvider, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base_url required for custom provider")
	}

	openai, err := NewOpenAIProvider(config)
	if err != nil {
		return nil, err
	}

	return &CustomProvider{openai: openai}, nil
}

// Complete delegates to the underlying OpenAI provider.
func (p *CustomProvider) Complete(ctx context.Context, req Request) (Response, error) {
	return p.openai.Complete(ctx, req)
}

// Name returns the provider identifier.
func (p *CustomProvider) Name() string { return "custom" }

// Model returns the configured model name.
func (p *CustomProvider) Model() string { return p.openai.Model() }
