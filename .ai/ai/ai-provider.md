# AI Provider System

**Status**: Design Complete  
**Last Updated**: 2026-03-29  

## 1. Provider Interface

### 1.1 Core Interface

```go
type Provider interface {
    // Complete sends a prompt and returns the response
    Complete(ctx context.Context, req Request) (Response, error)
    
    // Name returns the provider name
    Name() string
    
    // Model returns the configured model
    Model() string
}

type Request struct {
    SystemPrompt string
    UserPrompt   string
    MaxTokens    int
}

type Response struct {
    Content   string
    Usage     UsageStats
    Model     string
}

type UsageStats struct {
    PromptTokens     int
    CompletionTokens int
    TotalTokens      int
}
```

### 1.2 Provider Factory

```go
type ProviderConfig struct {
    Provider    string // openai | anthropic | ollama | custom
    Model       string
    APIKeyEnv   string
    BaseURL     string
    MaxTokens   int
    Timeout     time.Duration
}

func NewProvider(config ProviderConfig) (Provider, error) {
    switch config.Provider {
    case "openai":
        return NewOpenAIProvider(config)
    case "anthropic":
        return NewAnthropicProvider(config)
    case "ollama":
        return NewOllamaProvider(config)
    case "custom":
        return NewCustomProvider(config) // OpenAI-compatible
    default:
        return nil, fmt.Errorf("unknown provider: %s", config.Provider)
    }
}
```

## 2. Provider Implementations

### 2.1 OpenAI

**API**: https://api.openai.com/v1  
**Models**: gpt-4o, gpt-4o-mini, gpt-3.5-turbo  
**Format**: OpenAI Chat Completions API

```go
type OpenAIProvider struct {
    client  *http.Client
    apiKey  string
    model   string
    baseURL string
}

func (p *OpenAIProvider) Complete(ctx context.Context, req Request) (Response, error) {
    payload := map[string]interface{}{
        "model": p.model,
        "messages": []map[string]string{
            {"role": "system", "content": req.SystemPrompt},
            {"role": "user", "content": req.UserPrompt},
        },
        "max_tokens": req.MaxTokens,
    }
    
    // HTTP POST to /v1/chat/completions
    // Parse response, extract content and usage
}
```

**Configuration**:
```yaml
ai:
  provider: openai
  model: gpt-4o-mini
  api_key_env: OPENAI_API_KEY
  max_tokens: 4096
```

### 2.2 Anthropic

**API**: https://api.anthropic.com  
**Models**: claude-3-5-sonnet-20241022, claude-3-opus-20240229, claude-3-haiku-20240307  
**Format**: Anthropic Messages API

```go
type AnthropicProvider struct {
    client  *http.Client
    apiKey  string
    model   string
}

func (p *AnthropicProvider) Complete(ctx context.Context, req Request) (Response, error) {
    payload := map[string]interface{}{
        "model": p.model,
        "system": req.SystemPrompt,
        "messages": []map[string]string{
            {"role": "user", "content": req.UserPrompt},
        },
        "max_tokens": req.MaxTokens,
    }
    
    // HTTP POST to /v1/messages
    // Parse response, extract content and usage
}
```

**Configuration**:
```yaml
ai:
  provider: anthropic
  model: claude-3-5-sonnet-20241022
  api_key_env: ANTHROPIC_API_KEY
  max_tokens: 4096
```

### 2.3 Ollama (Local)

**API**: http://localhost:11434 (default)  
**Models**: llama3, mistral, codellama, etc.  
**Format**: Ollama API (OpenAI-compatible)

```go
type OllamaProvider struct {
    client  *http.Client
    baseURL string
    model   string
}

func (p *OllamaProvider) Complete(ctx context.Context, req Request) (Response, error) {
    payload := map[string]interface{}{
        "model": p.model,
        "system": req.SystemPrompt,
        "prompt": req.UserPrompt,
        "stream": false,
        "options": map[string]interface{}{
            "num_predict": req.MaxTokens,
        },
    }
    
    // HTTP POST to /api/generate
    // Parse response, extract content
    // Note: Ollama doesn't return token counts, estimate based on chars
}
```

**Configuration**:
```yaml
ai:
  provider: ollama
  base_url: http://localhost:11434
  model: llama3
  max_tokens: 4096
```

**Health Check**:
```go
func (p *OllamaProvider) CheckHealth() error {
    // GET /api/tags
    // Check if model is available
}
```

### 2.4 Custom (OpenAI-Compatible)

**Use Cases**:
- Together AI
- Groq
- Anyscale
- OpenRouter
- Self-hosted (vLLM, LocalAI)

**Implementation**: Same as OpenAI provider, but with custom base URL

```go
type CustomProvider struct {
    // Reuse OpenAI implementation
    openaiProvider *OpenAIProvider
}

func NewCustomProvider(config ProviderConfig) (*CustomProvider, error) {
    // Validate base URL
    if config.BaseURL == "" {
        return nil, errors.New("base_url required for custom provider")
    }
    
    // Create OpenAI provider with custom base URL
    config.Provider = "openai" // Reuse OpenAI client
    provider, err := NewOpenAIProvider(config)
    if err != nil {
        return nil, err
    }
    
    return &CustomProvider{openaiProvider: provider}, nil
}
```

**Configuration**:
```yaml
ai:
  provider: custom
  base_url: https://api.together.xyz/v1
  api_key_env: TOGETHER_API_KEY
  model: meta-llama/Llama-3-70b-chat-hf
  max_tokens: 4096
```

## 3. Configuration Management

### 3.1 Config Structure

```go
type AIConfig struct {
    Provider        string         `yaml:"provider"`
    Model           string         `yaml:"model"`
    APIKeyEnv       string         `yaml:"api_key_env"`
    BaseURL         string         `yaml:"base_url"`
    MaxTokens       int            `yaml:"max_tokens"`
    Timeout         time.Duration  `yaml:"timeout"`
    Batch           BatchConfig    `yaml:"batch"`
    Concurrency     ConcurrencyConfig `yaml:"concurrency"`
    RateLimit       RateLimitConfig   `yaml:"rate_limit"`
    Cache           CacheConfig       `yaml:"cache"`
}

type BatchConfig struct {
    Size            int `yaml:"size"`
    MaxContextNotes int `yaml:"max_context_notes"`
}

type ConcurrencyConfig struct {
    Enabled     bool `yaml:"enabled"`
    MaxRequests int  `yaml:"max_requests"`
}

type RateLimitConfig struct {
    RequestsPerMinute int    `yaml:"requests_per_minute"`
    RetryOnLimit      bool   `yaml:"retry_on_limit"`
    MaxRetries          int    `yaml:"max_retries"`
    Backoff           string `yaml:"backoff"` // exponential | linear
}

type CacheConfig struct {
    Enabled bool          `yaml:"enabled"`
    TTL     time.Duration `yaml:"ttl"`
    Path    string        `yaml:"path"`
}
```

### 3.2 Loading Configuration

```go
func LoadAIConfig(configPath string) (*AIConfig, error) {
    // Load from ~/.config/hotnote/config.yaml
    // Extract 'ai' section
    // Apply defaults for missing values
}

func (c *AIConfig) ApplyDefaults() {
    if c.MaxTokens == 0 {
        c.MaxTokens = 4096
    }
    if c.Timeout == 0 {
        c.Timeout = 60 * time.Second
    }
    if c.Batch.Size == 0 {
        c.Batch.Size = 5
    }
    if c.Batch.MaxContextNotes == 0 {
        c.Batch.MaxContextNotes = 20
    }
    if c.RateLimit.RequestsPerMinute == 0 {
        c.RateLimit.RequestsPerMinute = 50
    }
    if c.RateLimit.MaxRetries == 0 {
        c.RateLimit.MaxRetries = 3
    }
}
```

### 3.3 Saving Configuration

```go
func (c *AIConfig) Save(configPath string) error {
    // Serialize to YAML
    // Write to config file
    // Handle atomic writes
}
```

## 4. Setup Command Implementation

### 4.1 Interactive Flow

```go
type SetupStep struct {
    Prompt      string
    Options     []string
    Default     string
    Validate    func(string) error
    HelpText    string
}

func RunInteractiveSetup() error {
    steps := []SetupStep{
        {
            Prompt: "Choose a provider",
            Options: []string{"openai", "anthropic", "ollama", "custom"},
            Default: "openai",
        },
        // Provider-specific steps...
    }
    
    // Execute TUI prompts
    // Collect responses
    // Test configuration
    // Save if successful
}
```

### 4.2 Provider-Specific Steps

**OpenAI**:
1. Provider: openai
2. API key env var (detect OPENAI_API_KEY)
3. Model (gpt-4o-mini, gpt-4o)
4. Test request

**Anthropic**:
1. Provider: anthropic
2. API key env var (detect ANTHROPIC_API_KEY)
3. Model (claude-3-5-sonnet-20241022, claude-3-haiku-20240307)
4. Test request

**Ollama**:
1. Provider: ollama
2. Endpoint URL (default: http://localhost:11434)
3. Available models (fetch from /api/tags)
4. Select model
5. Test request
6. Privacy notice

**Custom**:
1. Provider: custom
2. Endpoint URL
3. API key env var
4. Model
5. Test request

### 4.3 API Key Detection

```go
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
```

### 4.4 Test Request

```go
func TestProvider(config ProviderConfig) (*TestResult, error) {
    provider, err := NewProvider(config)
    if err != nil {
        return nil, err
    }
    
    start := time.Now()
    resp, err := provider.Complete(context.Background(), Request{
        SystemPrompt: "You are a helpful assistant.",
        UserPrompt:   "Say 'Hotnote AI is ready' and nothing else.",
        MaxTokens:    50,
    })
    
    if err != nil {
        return nil, err
    }
    
    return &TestResult{
        Success:    true,
        Latency:    time.Since(start),
        Tokens:     resp.Usage.TotalTokens,
        Response:   resp.Content,
    }, nil
}
```

## 5. Error Handling

### 5.1 Provider Errors

```go
type ProviderError struct {
    Provider string
    Type     string // authentication | rate_limit | timeout | invalid_request | server_error
    Message  string
    RawError error
}

func (e *ProviderError) Error() string {
    return fmt.Sprintf("%s provider error: %s - %s", e.Provider, e.Type, e.Message)
}
```

### 5.2 Error Mapping

| Provider Error | Hotnote Error | Exit Code |
|----------------|---------------|-----------|
| 401 Unauthorized | ErrAIAPIKeyMissing | 5 |
| 429 Rate Limit | ErrAIRequestFailed (retry) | 6 |
| 408/Timeout | ErrAITimeout | 6 |
| 400 Invalid Request | ErrAIRequestFailed | 6 |
| 5xx Server Error | ErrAIRequestFailed | 6 |
| Network Error | ErrAIRequestFailed | 6 |

## 6. Provider Capabilities

### 6.1 Feature Matrix

| Feature | OpenAI | Anthropic | Ollama | Custom |
|---------|--------|-----------|--------|--------|
| Streaming | Yes | Yes | No | Yes* |
| Token Count | Yes | Yes | No** | Yes* |
| Tool Use | Yes | Yes | No | Yes* |
| JSON Mode | Yes | Yes | No | Yes* |

*Depends on provider implementation  
**Ollama doesn't return token counts, must estimate

### 6.2 Fallback Strategy

If a provider doesn't support a feature:
- Streaming: Return complete response
- Token counts: Estimate based on characters (1 token ≈ 4 chars)
- JSON mode: Parse text response, return error if invalid JSON

## 7. Security Considerations

### 7.1 API Key Storage

- **Never** store API keys in config files
- Use environment variable names only
- Validate key exists before operations
- Mask keys in error messages (show first/last 4 chars only)

### 7.2 Request Logging

```go
// Log requests at debug level (not in production)
if logLevel == Debug {
    log.Printf("AI Request: %s %s (tokens: %d)", 
        provider.Name(), 
        truncate(req.UserPrompt, 100),
        len(req.UserPrompt)/4) // rough estimate
}
```

### 7.3 Data Privacy

- Default providers (OpenAI, Anthropic) send data to cloud
- Ollama keeps data local
- Custom providers depend on endpoint
- Document this clearly in setup flow

## 8. Testing

### 8.1 Mock Provider

```go
type MockProvider struct {
    Responses map[string]string
    Delays    map[string]time.Duration
    Errors    map[string]error
}

func (m *MockProvider) Complete(ctx context.Context, req Request) (Response, error) {
    key := hash(req.SystemPrompt + req.UserPrompt)
    
    if err, ok := m.Errors[key]; ok {
        return Response{}, err
    }
    
    if delay, ok := m.Delays[key]; ok {
        time.Sleep(delay)
    }
    
    content := m.Responses[key]
    if content == "" {
        content = "Default mock response"
    }
    
    return Response{
        Content: content,
        Usage: UsageStats{
            PromptTokens:     len(req.UserPrompt) / 4,
            CompletionTokens: len(content) / 4,
            TotalTokens:      (len(req.UserPrompt) + len(content)) / 4,
        },
    }, nil
}
```

### 8.2 Provider Tests

Test each provider implementation with:
1. Successful completion
2. Authentication error
3. Rate limit error
4. Timeout error
5. Invalid response
6. Network error

## 9. Future Enhancements

### 9.1 Streaming Support

```go
type StreamingProvider interface {
    Provider
    CompleteStream(ctx context.Context, req Request) (<-chan StreamChunk, error)
}

type StreamChunk struct {
    Content string
    Done    bool
    Error   error
}
```

### 9.2 Caching

```go
type CachedProvider struct {
    Provider
    cache Cache
}

func (c *CachedProvider) Complete(ctx context.Context, req Request) (Response, error) {
    key := hash(req.SystemPrompt + req.UserPrompt)
    
    if cached, ok := c.cache.Get(key); ok {
        return cached, nil
    }
    
    resp, err := c.Provider.Complete(ctx, req)
    if err != nil {
        return Response{}, err
    }
    
    c.cache.Set(key, resp)
    return resp, nil
}
```

### 9.3 Multi-Provider Fallback

```go
type MultiProvider struct {
    Providers []Provider
}

func (m *MultiProvider) Complete(ctx context.Context, req Request) (Response, error) {
    for _, p := range m.Providers {
        resp, err := p.Complete(ctx, req)
        if err == nil {
            return resp, nil
        }
        // Log failure, try next
    }
    return Response{}, errors.New("all providers failed")
}
```
