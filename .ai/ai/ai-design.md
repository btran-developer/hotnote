# AI Features Design Document

**Status**: Design Phase Complete  
**Last Updated**: 2026-03-29  

## 1. Overview

This document provides the comprehensive design for AI features in Hotnote. It covers architecture decisions, provider strategy, CLI surface, and implementation guidance.

## 2. Core Philosophy

### 2.1 Design Principles

1. **Explicit Configuration**: No implicit AI providers. Users must explicitly configure their provider via `hotnote ai setup`.
2. **Provider Agnostic**: Support multiple LLM providers through a common interface.
3. **Privacy-First Options**: Local models (Ollama) alongside cloud providers.
4. **Cost Transparency**: Always show token usage and estimated costs.
5. **Agent-Friendly**: JSON output with citations for AI agents.

### 2.2 Audience

- **Human users**: Interactive, human-readable output with progress indicators
- **AI agents**: Machine-readable JSON with deterministic exit codes and citations

## 3. Architecture

### 3.1 Package Structure

```
internal/ai/
├── provider.go        # Provider interface
├── openai.go          # OpenAI implementation
├── anthropic.go       # Anthropic implementation
├── ollama.go          # Local Ollama implementation
├── config.go          # AI configuration
├── context.go         # Context building logic
├── cache.go           # Response caching (future)
├── errors.go          # AI-specific errors
└── prompts/           # System prompts (future)
    ├── search.go
    ├── summarize.go
    ├── related.go
    ├── tags.go
    ├── ask.go
    ├── extract.go
    └── dedup.go

cmd/
├── ai.go              # Parent 'ai' command
├── ai_search.go       # 'ai search' subcommand
├── ai_summarize.go    # 'ai summarize' subcommand
├── ai_related.go      # 'ai related' subcommand
├── ai_tags.go         # 'ai tags' subcommand
├── ai_ask.go          # 'ai ask' subcommand
├── ai_extract.go      # 'ai extract' subcommand
├── ai_dedup.go        # 'ai dedup' subcommand
└── ai_setup.go        # 'ai setup' command
```

### 3.2 Provider Interface

```go
type Provider interface {
    // Complete sends a prompt and returns the response
    Complete(ctx context.Context, req Request) (Response, error)
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

## 4. Configuration

### 4.1 Config Schema

Configuration stored in `~/.config/hotnote/config.yaml`:

```yaml
current_workspace: default
workspaces:
  default: /Users/.../workspaces/default

ai:
  provider: openai           # openai | anthropic | ollama | custom
  model: gpt-4o-mini
  api_key_env: OPENAI_API_KEY  # Environment variable name
  base_url: ""               # For Ollama or custom endpoints
  max_tokens: 4096
  timeout: 60s
  
  batch:
    size: 5                  # Notes per batch for batch operations
    max_context_notes: 20    # Max notes to include in context
  
  concurrency:
    enabled: false           # Sequential by default
    max_requests: 3
  
  rate_limit:
    requests_per_minute: 50
    retry_on_limit: true
    max_retries: 3
    backoff: exponential     # exponential | linear
  
  cache:
    enabled: false           # Future feature
    ttl: 24h
    path: ~/.cache/hotnote/ai/
```

### 4.2 Setup Command

```bash
$ hotnote ai setup

? Choose a provider:
  > OpenAI
    Anthropic
    Ollama (local)
    Custom (OpenAI-compatible)

? API key environment variable: OPENAI_API_KEY
  [current value: sk-...abc123]

? Model (default: gpt-4o-mini): gpt-4o-mini

✓ Configuration saved to ~/.config/hotnote/config.yaml
✓ Test request successful (45ms, 150 tokens)

Run 'hotnote ai --help' to see available commands.
```

**Non-interactive mode:**
```bash
hotnote ai setup --provider openai --model gpt-4o-mini --api-key-env OPENAI_API_KEY
```

## 5. Key Design Decisions

### 5.1 Provider Strategy

| Provider | Type | Use Case |
|----------|------|----------|
| OpenAI | Cloud | Balanced quality/cost (GPT-4o-mini) |
| Anthropic | Cloud | Best quality (Claude 3.5 Sonnet) |
| Ollama | Local | Privacy-conscious users |
| Custom | Varies | OpenAI-compatible APIs (Together, Groq, etc.) |

**Decision**: Support explicit configuration only. No default provider.

### 5.2 Streaming

**Decision**: Defer streaming support.

**Reasoning**:
- Adds complexity (different output handling for human vs JSON)
- JSON mode requires complete response anyway
- Most operations are fast enough (1-5s)
- Can add later for `ai ask` if needed

### 5.3 Batch Operations

**Decision**: Hybrid approach with batching.

**Strategy**:
- Group notes into batches (default: 5 notes)
- Process each batch in single LLM call
- Aggregate results
- Show progress after each batch

**Benefits**:
- Balances speed and reliability
- Manages context window limits
- Partial success possible
- Progress visible to user

### 5.4 Context Window Management

**Decision**: Tiered retrieval with user control.

**Strategy**:
1. **Fast Retrieval**: Index filename + frontmatter + first 500 chars
2. **Candidate Scoring**: Score all notes by relevance
3. **Context Assembly**: Include top K candidates (default: 20)
4. **Truncation**: Truncate each note if total exceeds threshold

**User Control**:
```bash
--context-limit N     # Max notes to include
--context-tokens N    # Max tokens (advanced)
```

### 5.5 Privacy

**Decision**: No separate privacy flag.

**Reasoning**:
- Redundant with explicit provider setup
- Users know their chosen provider's privacy implications
- Setup flow clarifies which providers are cloud vs local

### 5.6 Rate Limiting

**Decision**: Client-side rate limiting with backoff.

**Implementation**:
- Configurable requests per minute (default: 50)
- Exponential backoff on rate limit errors
- Respect `Retry-After` header from APIs
- Show "Rate limited, waiting..." in progress output

### 5.7 Concurrency

**Decision**: Sequential default, opt-in concurrent.

**Strategy**:
- Sequential processing by default (safer for rate limits)
- Optional concurrency with `--concurrent` flag or config
- Max concurrent requests: 3 (configurable)

### 5.8 Caching

**Decision**: Optional caching, disabled by default.

**Future Implementation**:
- Cache key: hash(provider + model + system_prompt + context_hash + user_query)
- TTL: 24 hours (configurable)
- Cache invalidation: Note content changes, TTL expires, `--no-cache` flag

## 6. CLI Surface

See `ai-commands.md` for detailed command specifications.

### 6.1 Command Structure

```
hotnote ai <subcommand> [args] [flags]
```

### 6.2 Subcommands

| Command | Purpose |
|---------|---------|
| `search <query>` | Semantic search across notes |
| `summarize <note>` | Summarize a note or notes |
| `related <note>` | Find related notes |
| `tags <note>` | Suggest tags for a note |
| `ask <question>` | Q&A over notes |
| `extract <query>` | Extract relevant passages |
| `dedup [--scan]` | Find duplicate/similar notes |
| `setup` | Configure AI provider |

### 6.3 Common Flags

```
--json              Machine-readable output
--provider          Override configured provider
--model             Override configured model
--max-tokens        Override max tokens
--context-limit     Max notes to include in context
--concurrent        Enable concurrent processing
--no-cache          Skip cache (when caching enabled)
```

## 7. Output Formats

### 7.1 Human Output

```
$ hotnote ai search "API design decisions"

Found 3 relevant notes:

1. projects/api-design.md (score: 0.92)
   "Decided on REST over GraphQL for simplicity..."

2. meetings/q4-planning.md (score: 0.78)
   "API redesign scheduled for Q4..."

Usage: 1,700 tokens (~$0.003)
```

### 7.2 JSON Output

```json
{
  "query": "API design decisions",
  "results": [
    {
      "slug": "projects/api-design",
      "path": "/Users/.../projects/api-design.md",
      "score": 0.92,
      "excerpt": "Decided on REST over GraphQL for simplicity...",
      "metadata": {
        "tags": ["api", "architecture"],
        "created_at": "2026-03-15T10:00:00Z",
        "updated_at": "2026-03-20T14:30:00Z"
      }
    }
  ],
  "usage": {
    "prompt_tokens": 1500,
    "completion_tokens": 200,
    "total_tokens": 1700,
    "estimated_cost": 0.003
  }
}
```

## 8. Error Handling

### 8.1 New Error Types

```go
var (
    ErrAIProviderNotConfigured = errors.New("AI provider not configured")
    ErrAIAPIKeyMissing         = errors.New("AI API key not set")
    ErrAIRequestFailed         = errors.New("AI request failed")
    ErrAIResponseInvalid       = errors.New("AI response invalid")
    ErrAIContextTooLarge       = errors.New("context exceeds token limit")
    ErrAITimeout               = errors.New("AI request timed out")
)
```

### 8.2 Exit Codes

| Code | Meaning |
|------|---------|
| 5 | AI provider not configured |
| 6 | AI request failed |

### 8.3 Error Output

**Human:**
```
Error: AI provider not configured
Run 'hotnote ai setup' to configure your provider.
```

**JSON:**
```json
{"error": "AI provider not configured", "code": 5}
```

## 9. Testing Strategy

See `ai-testing.md` for detailed testing specifications.

### 9.1 Layers

1. **Unit Tests**: Test components without LLM calls
2. **Integration Tests**: Use mock provider with predefined responses
3. **E2E Tests**: Run against real LLM (optional, requires API key)
4. **Evaluation Tests**: Quality assurance with golden files

### 9.2 Test Data

- Store test notes in `testdata/notes/`
- Store golden files in `testdata/golden/`
- Use `testdata` to avoid polluting user directories

## 10. Future Considerations

### 10.1 Streaming

- Add streaming support for `ai ask` command
- Different output handling for human vs JSON mode

### 10.2 Embeddings

- Pre-compute embeddings for faster semantic search
- Store embeddings in local database
- Hybrid search: embeddings + LLM re-ranking

### 10.3 Multi-modal

- Support for images in notes
- Vision-capable models for diagram analysis

### 10.4 Advanced Context

- Long context window models
- Recursive summarization for very large note collections
- RAG with vector database

## 11. References

- [ai-commands.md](ai-commands.md) - CLI specifications
- [ai-provider.md](ai-provider.md) - Provider system
- [ai-context.md](ai-context.md) - Context management
- [ai-prompts.md](ai-prompts.md) - System prompts
- [ai-testing.md](ai-testing.md) - Testing strategy
- [ai-backlog.md](ai-backlog.md) - Implementation backlog
