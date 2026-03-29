# AI Testing Strategy

**Status**: Design Complete  
**Last Updated**: 2026-03-29  

## 1. Testing Philosophy

AI features require a multi-layered testing approach:
1. **Unit Tests**: Fast, deterministic, no LLM calls
2. **Integration Tests**: Mock LLM with predefined responses
3. **E2E Tests**: Real LLM calls (optional, requires API key)
4. **Evaluation Tests**: Quality assurance with golden files

## 2. Test Structure

```
internal/ai/
├── provider.go
├── provider_test.go           # Unit tests for provider logic
├── mocks/
│   └── mock_provider.go       # Mock provider implementation
└── testdata/
    ├── notes/                 # Test notes
    │   ├── simple.md
    │   ├── api-design.md
    │   └── meeting.md
    └── golden/                # Expected outputs
        ├── search_api_design.json
        ├── summarize_simple.json
        └── ...

cmd/
├── ai_search_test.go          # Search command tests
├── ai_summarize_test.go       # Summarize command tests
├── ai_related_test.go         # Related command tests
├── ai_tags_test.go            # Tags command tests
├── ai_ask_test.go             # Ask command tests
├── ai_extract_test.go         # Extract command tests
├── ai_dedup_test.go           # Dedup command tests
└── ai_setup_test.go           # Setup command tests
```

## 3. Layer 1: Unit Tests

### 3.1 Provider Interface Tests

Test individual components without actual LLM calls:

```go
// internal/ai/provider_test.go

func TestProviderConfigValidation(t *testing.T) {
    tests := []struct {
        name    string
        config  ProviderConfig
        wantErr bool
    }{
        {
            name: "valid openai config",
            config: ProviderConfig{
                Provider:  "openai",
                Model:     "gpt-4o-mini",
                APIKeyEnv: "OPENAI_API_KEY",
            },
            wantErr: false,
        },
        {
            name: "missing api key env",
            config: ProviderConfig{
                Provider: "openai",
                Model:    "gpt-4o-mini",
            },
            wantErr: true,
        },
        {
            name: "custom provider missing base_url",
            config: ProviderConfig{
                Provider: "custom",
                Model:    "model-name",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewProvider(tt.config)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 3.2 Context Building Tests

```go
// internal/ai/context_test.go

func TestContextBuilder(t *testing.T) {
    notes := []ScoredNote{
        {
            NoteIndex: NoteIndex{
                Slug:    "note1",
                Title:   "API Design",
                Excerpt: "REST vs GraphQL discussion...",
            },
            Score: 0.95,
        },
        {
            NoteIndex: NoteIndex{
                Slug:    "note2",
                Title:   "Meeting Notes",
                Excerpt: "Team discussed architecture...",
            },
            Score: 0.80,
        },
    }
    
    builder := ContextBuilder{
        MaxContextNotes: 5,
        MaxTokens:       1000,
    }
    
    context, included, err := builder.BuildContext(notes)
    assert.NoError(t, err)
    assert.Len(t, included, 2)
    assert.Contains(t, context, "note1")
    assert.Contains(t, context, "note2")
}

func TestNoteScoring(t *testing.T) {
    note := NoteIndex{
        Title:   "API Design Decisions",
        Tags:    []string{"api", "architecture"},
        Excerpt: "We decided to use REST over GraphQL...",
    }
    
    score := calculateRelevance(note, "API design")
    assert.Greater(t, score, 0.7) // Should be highly relevant
    
    score2 := calculateRelevance(note, "database")
    assert.Less(t, score2, 0.5) // Should not be relevant
}
```

### 3.3 Token Estimation Tests

```go
// internal/ai/context_test.go

func TestTokenEstimation(t *testing.T) {
    tests := []struct {
        text     string
        expected int
    }{
        {"Hello", 1},                          // ~1 token
        {"Hello world", 2},                    // ~2 tokens
        {strings.Repeat("a", 400), 100},       // ~100 tokens
        {strings.Repeat("hello ", 100), 150},    // ~150 tokens
    }
    
    for _, tt := range tests {
        tokens := EstimateTokens(tt.text)
        // Allow 20% margin for rough estimation
        assert.InDelta(t, tt.expected, tokens, float64(tt.expected)*0.2)
    }
}
```

### 3.4 Prompt Construction Tests

```go
// internal/ai/prompts_test.go

func TestGetSystemPrompt(t *testing.T) {
    prompt := GetSystemPrompt("search")
    assert.NotEmpty(t, prompt)
    assert.Contains(t, prompt, "JSON")
    assert.Contains(t, prompt, "slug")
    assert.Contains(t, prompt, "score")
}

func TestBuildPrompt(t *testing.T) {
    context := "<note>test content</note>"
    query := "test query"
    
    systemPrompt, userPrompt := BuildPrompt("search", context, query)
    
    assert.NotEmpty(t, systemPrompt)
    assert.Contains(t, userPrompt, "<context>")
    assert.Contains(t, userPrompt, "<query>")
    assert.Contains(t, userPrompt, context)
    assert.Contains(t, userPrompt, query)
}
```

### 3.5 Configuration Tests

```go
// internal/ai/config_test.go

func TestLoadAIConfig(t *testing.T) {
    // Create temporary config file
    configContent := `
ai:
  provider: openai
  model: gpt-4o-mini
  api_key_env: OPENAI_API_KEY
  max_tokens: 2048
`
    tmpFile := createTempConfig(t, configContent)
    defer os.Remove(tmpFile)
    
    config, err := LoadAIConfig(tmpFile)
    assert.NoError(t, err)
    assert.Equal(t, "openai", config.Provider)
    assert.Equal(t, "gpt-4o-mini", config.Model)
    assert.Equal(t, "OPENAI_API_KEY", config.APIKeyEnv)
    assert.Equal(t, 2048, config.MaxTokens)
}

func TestConfigDefaults(t *testing.T) {
    config := &AIConfig{}
    config.ApplyDefaults()
    
    assert.Equal(t, 4096, config.MaxTokens)
    assert.Equal(t, 5, config.Batch.Size)
    assert.Equal(t, 20, config.Batch.MaxContextNotes)
    assert.Equal(t, 50, config.RateLimit.RequestsPerMinute)
}
```

## 4. Layer 2: Integration Tests (Mock Provider)

### 4.1 Mock Provider Implementation

```go
// internal/ai/mocks/mock_provider.go

package mocks

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "hotnotego/internal/ai"
)

type MockProvider struct {
    Responses map[string]string           // key -> response content
    Errors    map[string]error            // key -> error
    Delays    map[string]time.Duration    // key -> delay
    Calls     []ai.Request                // record all calls
}

func NewMockProvider() *MockProvider {
    return &MockProvider{
        Responses: make(map[string]string),
        Errors:    make(map[string]error),
        Delays:    make(map[string]time.Duration),
        Calls:     []ai.Request{},
    }
}

func (m *MockProvider) Complete(ctx context.Context, req ai.Request) (ai.Response, error) {
    m.Calls = append(m.Calls, req)
    
    key := hashPrompt(req.SystemPrompt + req.UserPrompt)
    
    // Check for predefined error
    if err, ok := m.Errors[key]; ok {
        return ai.Response{}, err
    }
    
    // Apply delay if specified
    if delay, ok := m.Delays[key]; ok {
        time.Sleep(delay)
    }
    
    // Return predefined response or default
    content := m.Responses[key]
    if content == "" {
        content = `{"results": []}` // Default empty response
    }
    
    return ai.Response{
        Content: content,
        Usage: ai.UsageStats{
            PromptTokens:     len(req.UserPrompt) / 4,
            CompletionTokens: len(content) / 4,
            TotalTokens:      (len(req.UserPrompt) + len(content)) / 4,
        },
    }, nil
}

func (m *MockProvider) Name() string {
    return "mock"
}

func (m *MockProvider) Model() string {
    return "mock-model"
}

func (m *MockProvider) SetResponse(promptPattern string, response string) {
    key := hashPrompt(promptPattern)
    m.Responses[key] = response
}

func (m *MockProvider) SetError(promptPattern string, err error) {
    key := hashPrompt(promptPattern)
    m.Errors[key] = err
}

func hashPrompt(prompt string) string {
    // Simple hash for matching prompts
    return fmt.Sprintf("%d", len(prompt))
}
```

### 4.2 Command Integration Tests

```go
// cmd/ai_search_test.go

func TestSearchCommand(t *testing.T) {
    mock := mocks.NewMockProvider()
    
    // Setup mock response
    mockResponse := `{
        "results": [
            {
                "slug": "projects/api-design",
                "score": 0.92,
                "excerpt": "Decided on REST over GraphQL...",
                "reason": "Direct match"
            }
        ]
    }`
    mock.SetResponse("search", mockResponse)
    
    // Create command with mock provider
    cmd := NewSearchCommandWithProvider(mock)
    
    // Execute
    output, err := cmd.Run("API design")
    
    // Assert
    assert.NoError(t, err)
    assert.Contains(t, output, "projects/api-design")
    assert.Contains(t, output, "0.92")
    
    // Verify provider was called
    assert.Len(t, mock.Calls, 1)
    assert.Contains(t, mock.Calls[0].UserPrompt, "API design")
}

func TestSearchCommandJSON(t *testing.T) {
    mock := mocks.NewMockProvider()
    mock.SetResponse("search", `{"results": [{"slug": "test", "score": 0.9}]}`)
    
    cmd := NewSearchCommandWithProvider(mock)
    output, err := cmd.RunWithFlags("query", map[string]interface{}{"json": true})
    
    assert.NoError(t, err)
    
    // Verify JSON structure
    var result map[string]interface{}
    err = json.Unmarshal([]byte(output), &result)
    assert.NoError(t, err)
    assert.Contains(t, result, "results")
}

func TestSearchCommandProviderError(t *testing.T) {
    mock := mocks.NewMockProvider()
    mock.SetError("search", errors.New("API rate limit exceeded"))
    
    cmd := NewSearchCommandWithProvider(mock)
    _, err := cmd.Run("query")
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "rate limit")
}
```

### 4.3 Batch Processing Tests

```go
// cmd/ai_summarize_test.go

func TestSummarizeBatchCommand(t *testing.T) {
    mock := mocks.NewMockProvider()
    
    // Setup responses for multiple batches
    mock.SetResponse("batch1", `{
        "summaries": [
            {"slug": "note1", "summary": "Summary 1", "topics": ["topic1"]},
            {"slug": "note2", "summary": "Summary 2", "topics": ["topic2"]}
        ]
    }`)
    mock.SetResponse("batch2", `{
        "summaries": [
            {"slug": "note3", "summary": "Summary 3", "topics": ["topic3"]}
        ]
    }`)
    
    cmd := NewSummarizeCommandWithProvider(mock, BatchConfig{Size: 2})
    output, err := cmd.RunBatch([]string{"note1", "note2", "note3"})
    
    assert.NoError(t, err)
    
    // Should have called provider twice (2 batches of 2, 1 note left)
    assert.Len(t, mock.Calls, 2)
    
    // Verify all notes summarized
    assert.Contains(t, output, "Summary 1")
    assert.Contains(t, output, "Summary 2")
    assert.Contains(t, output, "Summary 3")
}
```

## 5. Layer 3: E2E Tests

### 5.1 E2E Test Setup

```go
// cmd/ai_e2e_test.go

//go:build e2e
// +build e2e

package cmd

import (
    "os"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestSearchE2E(t *testing.T) {
    // Skip if no API key
    if os.Getenv("OPENAI_API_KEY") == "" {
        t.Skip("Skipping E2E test: OPENAI_API_KEY not set")
    }
    
    // Setup test workspace
    workspace := setupTestWorkspace(t)
    defer cleanupTestWorkspace(workspace)
    
    // Create test notes
    createTestNote(t, workspace, "api-design.md", `# API Design

We decided to use REST over GraphQL.
`)
    
    createTestNote(t, workspace, "database.md", `# Database Design

We chose PostgreSQL over MySQL.
`)
    
    // Run search command
    cmd := NewSearchCommand()
    output, err := cmd.Run("API design decisions")
    
    assert.NoError(t, err)
    assert.Contains(t, output, "api-design")
    assert.NotContains(t, output, "database") // Should not match
}

func TestAskE2E(t *testing.T) {
    if os.Getenv("OPENAI_API_KEY") == "" {
        t.Skip("Skipping E2E test: OPENAI_API_KEY not set")
    }
    
    workspace := setupTestWorkspace(t)
    defer cleanupTestWorkspace(workspace)
    
    createTestNote(t, workspace, "decisions.md", `# Decisions

We chose REST over GraphQL on March 15, 2026.
Team voted 4-1 in favor of REST.
`)
    
    cmd := NewAskCommand()
    output, err := cmd.Run("What did we decide about the API?")
    
    assert.NoError(t, err)
    assert.Contains(t, output, "REST")
    assert.Contains(t, output, "citations") // Should include citations
}
```

### 5.2 E2E Test Helpers

```go
// cmd/test_helpers_test.go

func setupTestWorkspace(t *testing.T) string {
    tmpDir := t.TempDir()
    
    // Initialize workspace
    cmd := exec.Command("hotnote", "workspace", "init")
    cmd.Env = append(os.Environ(), "HOME="+tmpDir)
    err := cmd.Run()
    require.NoError(t, err)
    
    return filepath.Join(tmpDir, ".local", "share", "hotnote", "workspaces", "default")
}

func cleanupTestWorkspace(path string) {
    os.RemoveAll(path)
}

func createTestNote(t *testing.T, workspace string, filename string, content string) {
    path := filepath.Join(workspace, filename)
    err := os.WriteFile(path, []byte(content), 0644)
    require.NoError(t, err)
}
```

## 6. Layer 4: Evaluation Tests

### 6.1 Golden Files

Store expected outputs for quality validation:

```json
// testdata/golden/search_api_design.json
{
  "query": "API design decisions",
  "expected_results": [
    {
      "slug": "projects/api-design",
      "min_score": 0.8,
      "required_keywords": ["REST", "GraphQL"]
    }
  ],
  "notes": [
    "projects/api-design.md",
    "meetings/2026-03-22.md",
    "notes/architecture.md"
  ]
}
```

### 6.2 Evaluation Test Implementation

```go
// cmd/ai_evaluation_test.go

//go:build evaluation
// +build evaluation

package cmd

func TestSearchQuality(t *testing.T) {
    // Load golden file
    golden := loadGoldenFile(t, "search_api_design.json")
    
    // Create test notes from golden file
    workspace := setupTestWorkspace(t)
    for _, notePath := range golden.Notes {
        createTestNoteFromFile(t, workspace, notePath)
    }
    
    // Run search
    cmd := NewSearchCommand()
    output, err := cmd.Run(golden.Query)
    require.NoError(t, err)
    
    // Parse response
    var result SearchResult
    err = json.Unmarshal([]byte(output), &result)
    require.NoError(t, err)
    
    // Validate against expected results
    for _, expected := range golden.ExpectedResults {
        found := findResult(result.Results, expected.Slug)
        assert.NotNil(t, found, "Expected to find %s", expected.Slug)
        assert.GreaterOrEqual(t, found.Score, expected.MinScore,
            "Score for %s should be >= %f", expected.Slug, expected.MinScore)
        
        // Check required keywords in excerpt
        for _, keyword := range expected.RequiredKeywords {
            assert.Contains(t, found.Excerpt, keyword,
                "Excerpt for %s should contain '%s'", expected.Slug, keyword)
        }
    }
}

func loadGoldenFile(t *testing.T, filename string) GoldenFile {
    path := filepath.Join("testdata", "golden", filename)
    data, err := os.ReadFile(path)
    require.NoError(t, err)
    
    var golden GoldenFile
    err = json.Unmarshal(data, &golden)
    require.NoError(t, err)
    
    return golden
}
```

### 6.3 Quality Metrics

Track quality metrics over time:

```go
type QualityMetrics struct {
    TestName         string
    Timestamp        time.Time
    TotalTests       int
    PassedTests      int
    AverageScore     float64
    FalsePositives   int
    FalseNegatives   int
}

func RecordMetrics(metrics QualityMetrics) {
    // Write to metrics file for tracking
    // Can be used to detect quality regression
}
```

## 7. Test Data Management

### 7.1 Test Notes

Create diverse test notes in `testdata/notes/`:

```markdown
<!-- testdata/notes/simple.md -->
# Simple Note

This is a simple test note.
```

```markdown
<!-- testdata/notes/api-design.md -->
---
title: API Design Decisions
created_at: 2026-03-15T10:00:00Z
tags: [api, architecture, rest]
---

# API Design Decisions

Date: March 15, 2026

## Decisions

1. **Protocol**: REST over GraphQL
   - Reason: Team familiarity and simplicity
   - Vote: 4-1 in favor of REST

2. **Specification**: JSON:API
   - Reason: Consistency across endpoints

3. **Versioning**: URL versioning (/v1/, /v2/)
   - Deprecation policy: 6 months notice

## GraphQL Evaluation

Deferred to Q4 2026 if REST proves insufficient.
```

### 7.2 Note Generator

Generate test notes programmatically:

```go
func generateTestNotes(count int) []TestNote {
    var notes []TestNote
    
    templates := []string{
        "Meeting notes from %s",
        "API design discussion - %s",
        "Architecture decision record - %s",
        "Research notes: %s",
    }
    
    for i := 0; i < count; i++ {
        template := templates[i%len(templates)]
        date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
        
        notes = append(notes, TestNote{
            Title:   fmt.Sprintf(template, date),
            Content: generateContent(i),
            Tags:    generateTags(i),
        })
    }
    
    return notes
}
```

## 8. CI/CD Integration

### 8.1 Test Stages

```yaml
# .github/workflows/test.yml

name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      
      - name: Unit Tests
        run: go test ./... -v -short
  
  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      
      - name: Integration Tests
        run: go test ./... -v -tags=integration
  
  # E2E tests only on main branch or manual trigger
  e2e-tests:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main' || github.event_name == 'workflow_dispatch'
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      
      - name: E2E Tests
        env:
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
        run: go test ./cmd -v -tags=e2e -timeout 30m
```

### 8.2 Test Coverage

Track test coverage for AI package:

```bash
go test -coverprofile=coverage.out ./internal/ai/...
go tool cover -html=coverage.out
```

Target: 80%+ coverage for AI package.

## 9. Performance Tests

### 9.1 Benchmark Tests

```go
// internal/ai/benchmark_test.go

func BenchmarkContextBuilding(b *testing.B) {
    // Generate 1000 test notes
    notes := generateTestNotes(1000)
    builder := ContextBuilder{}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _, _ = builder.BuildContext(notes)
    }
}

func BenchmarkNoteScoring(b *testing.B) {
    notes := generateTestNotes(100)
    query := "API design"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ScoreNotes(notes, query)
    }
}
```

### 9.2 Load Tests

```bash
# Test with many notes
mkdir -p testdata/many_notes
for i in {1..1000}; do
    echo "# Note $i" > testdata/many_notes/note_$i.md
done

time hotnote ai search "test query"
```

## 10. Testing Best Practices

### 10.1 Test Independence

Each test should be independent:
- Use `t.TempDir()` for temporary files
- Clean up after each test
- Don't depend on test execution order

### 10.2 Deterministic Tests

Use mock providers for deterministic behavior:
- Same input → same output
- No network calls in unit/integration tests
- Fixed random seeds if applicable

### 10.3 Comprehensive Coverage

Test these scenarios:
- Happy path
- Empty results
- API errors (rate limit, auth, timeout)
- Malformed responses
- Large notes (truncation)
- Many notes (batching)
- Special characters in content
- Unicode content

### 10.4 Documentation

Document test purpose:

```go
// TestSearchCommand validates the search command with a mock provider.
// It verifies:
// - Correct prompt construction
// - Proper JSON response parsing
// - Human-readable output formatting
// - Error handling for provider failures
func TestSearchCommand(t *testing.T) {
    // ...
}
```
