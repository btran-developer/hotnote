package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"hotnotego/internal/ai"
)

const (
	defaultOpenAIModel    = "gpt-4o-mini"
	defaultAnthropicModel = "claude-3-5-sonnet-20241022"
	defaultOllamaModel    = "llama3"
)

var aiSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure AI provider",
	Long:  `Interactive setup for AI provider configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
		if nonInteractive {
			// Use the persistent flags from parent aiCmd
			provider := aiProvider
			model := aiModel
			apiKeyEnv, _ := cmd.Flags().GetString("api-key-env")
			baseURL, _ := cmd.Flags().GetString("base-url")
			return runNonInteractiveSetup(provider, model, apiKeyEnv, baseURL)
		}
		return runInteractiveSetup()
	},
}

func init() {
	aiSetupCmd.Flags().Bool("non-interactive", false, "Run without prompts")
	aiSetupCmd.Flags().String("api-key-env", "", "Environment variable name for API key")
	aiSetupCmd.Flags().String("base-url", "", "Custom API endpoint (for ollama/custom)")
}

func getConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("HOME")
		if configDir == "" {
			configDir = "."
		}
		configDir = filepath.Join(configDir, ".config")
	}
	return filepath.Join(configDir, "hotnote", "config.yaml")
}

func runInteractiveSetup() error {
	fmt.Println("AI Provider Setup")
	fmt.Println("=================")
	fmt.Println()
	fmt.Println("Available providers:")
	fmt.Println("  1. OpenAI (cloud, gpt-4o-mini)")
	fmt.Println("  2. Anthropic (cloud, claude-3-5-sonnet)")
	fmt.Println("  3. Ollama (local, llama3)")
	fmt.Println("  4. Custom (OpenAI-compatible API)")
	fmt.Println()

	var choice string
	fmt.Print("Choose provider (1-4): ")
	if _, err := fmt.Scanln(&choice); err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	providerMap := map[string]string{
		"1": "openai",
		"2": "anthropic",
		"3": "ollama",
		"4": "custom",
	}

	selectedProvider, ok := providerMap[strings.TrimSpace(choice)]
	if !ok {
		return fmt.Errorf("invalid choice")
	}

	config := &ai.AIConfig{
		Provider: selectedProvider,
	}

	switch selectedProvider {
	case "openai":
		config.Model = defaultOpenAIModel
		detected := ai.DetectAPIKeyEnv("openai")
		if detected != "" {
			config.APIKeyEnv = detected
		}
	case "anthropic":
		config.Model = defaultAnthropicModel
		detected := ai.DetectAPIKeyEnv("anthropic")
		if detected != "" {
			config.APIKeyEnv = detected
		}
	case "ollama":
		config.Model = defaultOllamaModel
		config.BaseURL = "http://localhost:11434"
	}

	if config.APIKeyEnv == "" && selectedProvider != "ollama" {
		fmt.Print("API key environment variable: ")
		if _, err := fmt.Scanln(&config.APIKeyEnv); err != nil {
			return fmt.Errorf("read API key env var: %w", err)
		}
	}

	config.MaxTokens = 4096
	config.Timeout = 60

	// Load existing config to preserve Batch and Context defaults
	configPath := getConfigPath()
	existing, err := ai.LoadAIConfig(configPath)
	if err != nil {
		existing = ai.DefaultAIConfig()
	}

	// Only update the fields the user selected
	existing.Provider = config.Provider
	existing.Model = config.Model
	existing.APIKeyEnv = config.APIKeyEnv
	existing.BaseURL = config.BaseURL
	existing.MaxTokens = config.MaxTokens
	existing.Timeout = config.Timeout

	if err := ai.SaveAIConfig(configPath, existing); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Println()
	fmt.Println("Configuration saved to", configPath)
	fmt.Println()

	return testProvider(existing)
}

func runNonInteractiveSetup(provider, model, apiKeyEnv, baseURL string) error {
	if provider == "" {
		return fmt.Errorf("--provider is required in non-interactive mode")
	}

	validProviders := map[string]bool{
		"openai":    true,
		"anthropic": true,
		"ollama":    true,
		"custom":    true,
	}
	if !validProviders[provider] {
		return fmt.Errorf("unknown provider %q: must be openai, anthropic, ollama, or custom", provider)
	}

	config := ai.DefaultAIConfig()
	config.Provider = provider

	// Apply provider-specific defaults
	switch provider {
	case "openai":
		if model == "" {
			config.Model = defaultOpenAIModel
		}
		if apiKeyEnv == "" {
			detected := ai.DetectAPIKeyEnv("openai")
			if detected != "" {
				config.APIKeyEnv = detected
			}
		}
	case "anthropic":
		if model == "" {
			config.Model = defaultAnthropicModel
		}
		if apiKeyEnv == "" {
			detected := ai.DetectAPIKeyEnv("anthropic")
			if detected != "" {
				config.APIKeyEnv = detected
			}
		}
	case "ollama":
		if model == "" {
			config.Model = defaultOllamaModel
		}
		if baseURL == "" {
			config.BaseURL = "http://localhost:11434"
		}
	}

	if model != "" {
		config.Model = model
	}
	if apiKeyEnv != "" {
		config.APIKeyEnv = apiKeyEnv
	}
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	configPath := getConfigPath()
	if err := ai.SaveAIConfig(configPath, config); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Println("Configuration saved to", configPath)

	return testProvider(config)
}

func testProvider(config *ai.AIConfig) error {
	fmt.Println("Testing provider...")

	providerConfig := ai.ProviderConfig{
		Provider:  config.Provider,
		Model:     config.Model,
		APIKeyEnv: config.APIKeyEnv, // Pass env var name, not value
		BaseURL:   config.BaseURL,
		MaxTokens: config.MaxTokens,
		Timeout:   config.Timeout,
	}

	provider, err := ai.NewProvider(providerConfig)
	if err != nil {
		return fmt.Errorf("create provider: %w", err)
	}

	ctx := context.Background()
	resp, err := provider.Complete(ctx, ai.Request{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "Say 'Hotnote AI is ready' and nothing else.",
		MaxTokens:    50,
	})

	if err != nil {
		return fmt.Errorf("test request failed: %w", err)
	}

	fmt.Printf("Success! Response: %s\n", resp.Content)
	fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)

	return nil
}
