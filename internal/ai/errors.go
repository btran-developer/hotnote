package ai

import "errors"

var (
	// ErrAIProviderNotConfigured is returned when AI provider is not configured.
	ErrAIProviderNotConfigured = errors.New("AI provider not configured")

	// ErrAIAPIKeyMissing is returned when the API key environment variable is not set.
	ErrAIAPIKeyMissing = errors.New("AI API key not set")

	// ErrAIRequestFailed is returned when an AI API request fails.
	ErrAIRequestFailed = errors.New("AI request failed")

	// ErrAIResponseInvalid is returned when an AI response cannot be parsed.
	ErrAIResponseInvalid = errors.New("AI response invalid")

	// ErrAIContextTooLarge is returned when context exceeds the token limit.
	ErrAIContextTooLarge = errors.New("context exceeds token limit")

	// ErrAITimeout is returned when an AI request times out.
	ErrAITimeout = errors.New("AI request timed out")

	// ErrAIProviderNotFound is returned when an unknown provider name is used.
	ErrAIProviderNotFound = errors.New("unknown provider")
)

// ProviderError represents an error from an AI provider.
type ProviderError struct {
	Provider string
	Type     string
	Message  string
	RawError error
}

func (e *ProviderError) Error() string {
	return e.Provider + " provider error: " + e.Type + " - " + e.Message
}

func (e *ProviderError) Unwrap() error {
	return e.RawError
}

const (
	// ErrTypeAuth indicates an authentication failure.
	ErrTypeAuth = "authentication"
	// ErrTypeRateLimit indicates a rate limit was exceeded.
	ErrTypeRateLimit = "rate_limit"
	// ErrTypeTimeout indicates a request timeout.
	ErrTypeTimeout = "timeout"
	// ErrTypeInvalidReq indicates an invalid request.
	ErrTypeInvalidReq = "invalid_request"
	// ErrTypeServerError indicates a server error from the provider.
	ErrTypeServerError = "server_error"
)
