package ai

import "context"

// Provider defines the interface for LLM providers.
type Provider interface {
	Complete(ctx context.Context, req Request) (Response, error)
	Name() string
	Model() string
}

// Request encapsulates a prompt to send to the LLM.
type Request struct {
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
}

// Response contains the LLM's completion and usage statistics.
type Response struct {
	Content string
	Usage   UsageStats
	Model   string
}

// UsageStats tracks token usage for a request.
type UsageStats struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// ProviderConfig configures a provider instance.
type ProviderConfig struct {
	Provider  string
	Model     string
	APIKeyEnv string
	BaseURL   string
	MaxTokens int
	Timeout   int
}

// NewProvider creates a provider based on the given configuration.
func NewProvider(config ProviderConfig) (Provider, error) {
	switch config.Provider {
	case "openai":
		return NewOpenAIProvider(config)
	case "anthropic":
		return NewAnthropicProvider(config)
	case "ollama":
		return NewOllamaProvider(config)
	case "custom":
		return NewCustomProvider(config)
	default:
		return nil, ErrAIProviderNotFound
	}
}
