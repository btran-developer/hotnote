package ai

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"
)

// MockProvider is a mock implementation of the Provider interface for testing
type MockProvider struct {
	NameVal   string
	ModelVal  string
	Responses map[string]string
	Delays    map[string]int64 // in milliseconds
	Errors    map[string]error
}

// NewMockProvider creates a new mock provider
func NewMockProvider(name, model string) *MockProvider {
	return &MockProvider{
		NameVal:   name,
		ModelVal:  model,
		Responses: make(map[string]string),
		Delays:    make(map[string]int64),
		Errors:    make(map[string]error),
	}
}

// AddResponse adds a mock response for a given prompt
func (m *MockProvider) AddResponse(prompt, response string) {
	m.Responses[m.hashKey("", prompt)] = response
}

// AddError adds a mock error for a given prompt
func (m *MockProvider) AddError(prompt string, err error) {
	m.Errors[m.hashKey("", prompt)] = err
}

// Complete implements the Provider interface
func (m *MockProvider) Complete(ctx context.Context, req Request) (Response, error) {
	key := m.hashKey(req.SystemPrompt, req.UserPrompt)

	// Fall back to lookup without system prompt for simpler test registration
	if _, ok := m.Responses[key]; !ok {
		if _, ok := m.Errors[key]; !ok {
			key = m.hashKey("", req.UserPrompt)
		}
	}

	if err, ok := m.Errors[key]; ok {
		return Response{}, err
	}

	if delay, ok := m.Delays[key]; ok {
		timer := time.NewTimer(time.Duration(delay) * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return Response{}, ctx.Err()
		case <-timer.C:
		}
	}

	content := m.Responses[key]
	if content == "" {
		content = "Mock response"
	}

	return Response{
		Content: content,
		Usage: UsageStats{
			PromptTokens:     len(req.UserPrompt) / 4,
			CompletionTokens: len(content) / 4,
			TotalTokens:      (len(req.UserPrompt) + len(content)) / 4,
		},
		Model: m.ModelVal,
	}, nil
}

func (m *MockProvider) Name() string  { return m.NameVal }
func (m *MockProvider) Model() string { return m.ModelVal }

func (m *MockProvider) hashKey(system, user string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(system)) // hash/fnv Write never fails
	_, _ = h.Write([]byte(user))   // hash/fnv Write never fails
	return fmt.Sprintf("%x", h.Sum32())
}
