package ai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultAIConfig(t *testing.T) {
	cfg := DefaultAIConfig()

	if cfg.MaxTokens != 4096 {
		t.Errorf("Expected MaxTokens 4096, got %d", cfg.MaxTokens)
	}

	if cfg.Timeout != 60 {
		t.Errorf("Expected Timeout 60, got %d", cfg.Timeout)
	}

	if cfg.Batch.Size != 5 {
		t.Errorf("Expected Batch.Size 5, got %d", cfg.Batch.Size)
	}

	if cfg.Context.MaxNotes != 20 {
		t.Errorf("Expected Context.MaxNotes 20, got %d", cfg.Context.MaxNotes)
	}
}

func TestLoadAIConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Test loading non-existent file (should return defaults)
	cfg, err := LoadAIConfig(configPath)
	if err != nil {
		t.Fatalf("LoadAIConfig failed: %v", err)
	}

	if cfg.Provider != "" {
		t.Errorf("Expected empty provider, got %s", cfg.Provider)
	}

	if cfg.MaxTokens != 4096 {
		t.Errorf("Expected default MaxTokens 4096, got %d", cfg.MaxTokens)
	}
}

func TestSaveAIConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &AIConfig{
		Provider:  "openai",
		Model:     "gpt-4o-mini",
		APIKeyEnv: "OPENAI_API_KEY",
		MaxTokens: 2048,
	}

	err := SaveAIConfig(configPath, cfg)
	if err != nil {
		t.Fatalf("SaveAIConfig failed: %v", err)
	}

	// Verify file was created
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestDetectAPIKeyEnv(t *testing.T) {
	// Set up environment variables
	originalOpenAI := os.Getenv("OPENAI_API_KEY")
	originalAnthropic := os.Getenv("ANTHROPIC_API_KEY")

	defer func() {
		if originalOpenAI != "" {
			os.Setenv("OPENAI_API_KEY", originalOpenAI)
		} else {
			os.Unsetenv("OPENAI_API_KEY")
		}
		if originalAnthropic != "" {
			os.Setenv("ANTHROPIC_API_KEY", originalAnthropic)
		} else {
			os.Unsetenv("ANTHROPIC_API_KEY")
		}
	}()

	// Test OpenAI detection
	os.Setenv("OPENAI_API_KEY", "test-key")
	result := DetectAPIKeyEnv("openai")
	if result != "OPENAI_API_KEY" {
		t.Errorf("Expected OPENAI_API_KEY, got %s", result)
	}

	// Test Anthropic detection
	os.Unsetenv("OPENAI_API_KEY")
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	result = DetectAPIKeyEnv("anthropic")
	if result != "ANTHROPIC_API_KEY" {
		t.Errorf("Expected ANTHROPIC_API_KEY, got %s", result)
	}

	// Test unknown provider
	result = DetectAPIKeyEnv("unknown")
	if result != "" {
		t.Errorf("Expected empty string for unknown provider, got %s", result)
	}
}

func TestProviderError(t *testing.T) {
	err := &ProviderError{
		Provider: "openai",
		Type:     ErrTypeAuth,
		Message:  "Invalid API key",
	}

	expected := "openai provider error: authentication - Invalid API key"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// TestSaveAndLoadAIConfig tests the round-trip save and load of config
func TestSaveAndLoadAIConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	original := &AIConfig{
		Provider:  "openai",
		Model:     "gpt-4o-mini",
		APIKeyEnv: "OPENAI_API_KEY",
		MaxTokens: 2048,
		Timeout:   30,
		Batch: BatchConfig{
			Size:            10,
			MaxContextNotes: 30,
		},
		Context: ContextConfig{
			MaxNotes:         25,
			MaxTokens:        5000,
			CharLimitPerNote: 2500,
		},
	}

	// Save the config
	err := SaveAIConfig(configPath, original)
	if err != nil {
		t.Fatalf("SaveAIConfig failed: %v", err)
	}

	// Load the config
	loaded, err := LoadAIConfig(configPath)
	if err != nil {
		t.Fatalf("LoadAIConfig failed: %v", err)
	}

	// Verify all fields match
	if loaded.Provider != original.Provider {
		t.Errorf("Provider mismatch: %s != %s", loaded.Provider, original.Provider)
	}
	if loaded.Model != original.Model {
		t.Errorf("Model mismatch: %s != %s", loaded.Model, original.Model)
	}
	if loaded.APIKeyEnv != original.APIKeyEnv {
		t.Errorf("APIKeyEnv mismatch: %s != %s", loaded.APIKeyEnv, original.APIKeyEnv)
	}
	if loaded.MaxTokens != original.MaxTokens {
		t.Errorf("MaxTokens mismatch: %d != %d", loaded.MaxTokens, original.MaxTokens)
	}
	if loaded.Timeout != original.Timeout {
		t.Errorf("Timeout mismatch: %d != %d", loaded.Timeout, original.Timeout)
	}
	if loaded.Batch.Size != original.Batch.Size {
		t.Errorf("Batch.Size mismatch: %d != %d", loaded.Batch.Size, original.Batch.Size)
	}
	if loaded.Batch.MaxContextNotes != original.Batch.MaxContextNotes {
		t.Errorf("Batch.MaxContextNotes mismatch: %d != %d", loaded.Batch.MaxContextNotes, original.Batch.MaxContextNotes)
	}
	if loaded.Context.MaxNotes != original.Context.MaxNotes {
		t.Errorf("Context.MaxNotes mismatch: %d != %d", loaded.Context.MaxNotes, original.Context.MaxNotes)
	}
	if loaded.Context.MaxTokens != original.Context.MaxTokens {
		t.Errorf("Context.MaxTokens mismatch: %d != %d", loaded.Context.MaxTokens, original.Context.MaxTokens)
	}
	if loaded.Context.CharLimitPerNote != original.Context.CharLimitPerNote {
		t.Errorf("Context.CharLimitPerNote mismatch: %d != %d", loaded.Context.CharLimitPerNote, original.Context.CharLimitPerNote)
	}
}
