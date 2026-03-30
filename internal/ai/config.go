package ai

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// AIConfig stores the AI provider configuration
type AIConfig struct {
	Provider  string        `yaml:"provider"`
	Model     string        `yaml:"model"`
	APIKeyEnv string        `yaml:"api_key_env"`
	BaseURL   string        `yaml:"base_url"`
	MaxTokens int           `yaml:"max_tokens"`
	Timeout   int           `yaml:"timeout"`
	Batch     BatchConfig   `yaml:"batch"`
	Context   ContextConfig `yaml:"context"`
}

// BatchConfig holds batch processing settings
type BatchConfig struct {
	Size            int `yaml:"size"`
	MaxContextNotes int `yaml:"max_context_notes"`
}

// ContextConfig holds context building settings
type ContextConfig struct {
	MaxNotes         int `yaml:"max_notes"`
	MaxTokens        int `yaml:"max_tokens"`
	CharLimitPerNote int `yaml:"char_limit_per_note"`
}

// DefaultAIConfig returns a default configuration
func DefaultAIConfig() *AIConfig {
	return &AIConfig{
		MaxTokens: 4096,
		Timeout:   60,
		Batch: BatchConfig{
			Size:            5,
			MaxContextNotes: 20,
		},
		Context: ContextConfig{
			MaxNotes:         20,
			MaxTokens:        6000,
			CharLimitPerNote: 2000,
		},
	}
}

// LoadAIConfig loads AI configuration from the config file
func LoadAIConfig(configPath string) (*AIConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultAIConfig(), nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg struct {
		AI AIConfig `yaml:"ai"`
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Apply defaults for missing values
	defaults := DefaultAIConfig()

	if cfg.AI.MaxTokens == 0 {
		cfg.AI.MaxTokens = defaults.MaxTokens
	}
	if cfg.AI.Timeout == 0 {
		cfg.AI.Timeout = defaults.Timeout
	}
	if cfg.AI.Batch.Size == 0 {
		cfg.AI.Batch.Size = defaults.Batch.Size
	}
	if cfg.AI.Batch.MaxContextNotes == 0 {
		cfg.AI.Batch.MaxContextNotes = defaults.Batch.MaxContextNotes
	}
	if cfg.AI.Context.MaxNotes == 0 {
		cfg.AI.Context.MaxNotes = defaults.Context.MaxNotes
	}
	if cfg.AI.Context.MaxTokens == 0 {
		cfg.AI.Context.MaxTokens = defaults.Context.MaxTokens
	}
	if cfg.AI.Context.CharLimitPerNote == 0 {
		cfg.AI.Context.CharLimitPerNote = defaults.Context.CharLimitPerNote
	}

	return &cfg.AI, nil
}

// SaveAIConfig saves AI configuration, preserving existing non-AI sections
func SaveAIConfig(configPath string, cfg *AIConfig) error {
	var existing map[string]interface{}

	// Read existing config if it exists
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := yaml.Unmarshal(data, &existing); err != nil {
			return fmt.Errorf("parse existing config: %w", err)
		}
	}

	if existing == nil {
		existing = make(map[string]interface{})
	}

	// Update only the AI section
	existing["ai"] = cfg

	data, err = yaml.Marshal(existing)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// DetectAPIKeyEnv detects API key environment variable for a provider
func DetectAPIKeyEnv(provider string) string {
	switch provider {
	case "openai":
		if os.Getenv("OPENAI_API_KEY") != "" {
			return "OPENAI_API_KEY"
		}
	case "anthropic":
		if os.Getenv("ANTHROPIC_API_KEY") != "" {
			return "ANTHROPIC_API_KEY"
		}
	}
	return ""
}
