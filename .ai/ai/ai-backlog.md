# AI Features Backlog

**Status**: Planning Phase  
**Last Updated**: 2026-03-29  

## 1. Implementation Phases

### Phase 3A: AI Infrastructure
**Priority**: High  
**Status**: Not Started

#### Tasks

- [ ] **Provider Interface** (`internal/ai/provider.go`)
  - Define `Provider` interface
  - Define `Request` and `Response` types
  - Implement `ProviderFactory`
  - Document interface

- [ ] **OpenAI Provider** (`internal/ai/openai.go`)
  - Implement `Complete()` method
  - HTTP client configuration
  - Error handling and mapping
  - Request/response logging
  - Unit tests

- [ ] **Anthropic Provider** (`internal/ai/anthropic.go`)
  - Implement `Complete()` method
  - HTTP client configuration
  - Error handling and mapping
  - Unit tests

- [ ] **Ollama Provider** (`internal/ai/ollama.go`)
  - Implement `Complete()` method
  - Health check endpoint
  - Model availability check
  - Token estimation (Ollama doesn't return counts)
  - Unit tests

- [ ] **Custom Provider** (`internal/ai/custom.go`)
  - OpenAI-compatible wrapper
  - Base URL validation
  - Unit tests

- [ ] **Configuration System** (`internal/ai/config.go`)
  - `AIConfig` struct definition
  - Load/save configuration
  - Apply defaults
  - Validation
  - Unit tests

- [ ] **Context Management** (`internal/ai/context.go`)
  - `NoteIndex` struct
  - Index building from workspace
  - Note scoring algorithm
  - Context assembly
  - Token estimation
  - Truncation strategies
  - Unit tests

- [ ] **Error Handling** (`internal/ai/errors.go`)
  - Define AI-specific errors
  - Error mapping from providers
  - Exit codes
  - Unit tests

- [ ] **Setup Command** (`cmd/ai_setup.go`)
  - Interactive TUI setup
  - Provider selection
  - API key detection
  - Model selection
  - Test request
  - Non-interactive mode (`--provider`, `--model`, etc.)
  - Unit tests

- [ ] **Parent AI Command** (`cmd/ai.go`)
  - Define parent `ai` command
  - Register subcommands
  - Global flags (`--json`, `--provider`, `--model`, etc.)
  - Help text

**Dependencies**: None  
**Estimated Effort**: 3-4 weeks  

---

### Phase 3B: Core AI Commands
**Priority**: High  
**Status**: Not Started

#### Tasks

- [ ] **Search Command** (`cmd/ai_search.go`)
  - Implement `search` subcommand
  - Context building for search
  - Prompt construction
  - Response parsing
  - Human-readable output
  - JSON output
  - Flags: `--workspace`, `--folder`, `--recent`, `--limit`
  - Unit tests
  - Integration tests

- [ ] **Summarize Command** (`cmd/ai_summarize.go`)
  - Implement `summarize` subcommand
  - Single note summarization
  - Batch summarization (hybrid approach)
  - Progress indicators
  - Output formats: `bullet`, `paragraph`, `outline`
  - Flags: `--all`, `--recent`, `--format`, `--include-sources`
  - Unit tests
  - Integration tests

- [ ] **Related Command** (`cmd/ai_related.go`)
  - Implement `related` subcommand
  - Source note context
  - Candidate note scoring
  - Relationship explanation
  - Flags: `--limit`, `--explain`, `--threshold`
  - Unit tests
  - Integration tests

- [ ] **Tags Command** (`cmd/ai_tags.go`)
  - Implement `tags` subcommand
  - Tag suggestion logic
  - Dry-run mode (default)
  - Apply mode (`--apply`)
  - Confidence threshold
  - Frontmatter update logic
  - Flags: `--apply`, `--dry-run`, `--confidence`, `--all`
  - Unit tests
  - Integration tests

**Dependencies**: Phase 3A  
**Estimated Effort**: 2-3 weeks  

---

### Phase 3C: Advanced AI Commands
**Priority**: Medium  
**Status**: Not Started

#### Tasks

- [ ] **Ask (Q&A) Command** (`cmd/ai_ask.go`)
  - Implement `ask` subcommand
  - High-context retrieval strategy
  - Citation generation
  - Confidence levels
  - Flags: `--sources`, `--recent`, `--notes`, `--workspace`
  - Unit tests
  - Integration tests
  - E2E tests

- [ ] **Extract Command** (`cmd/ai_extract.go`)
  - Implement `extract` subcommand
  - Passage extraction logic
  - Output formats: `quotes`, `bullet`, `json`
  - Flags: `--format`, `--notes`, `--workspace`, `--folder`, `--recent`
  - Unit tests
  - Integration tests

- [ ] **Dedup Command** (`cmd/ai_dedup.go`)
  - Implement `dedup` subcommand
  - Note comparison algorithm
  - Pairwise similarity scoring
  - Optimization (don't compare all pairs)
  - Merge suggestions
  - Flags: `--scan`, `--notes`, `--threshold`, `--merge-suggest`, `--folder`, `--recent`
  - Unit tests
  - Integration tests

**Dependencies**: Phase 3B  
**Estimated Effort**: 2-3 weeks  

---

### Phase 3D: Polish & Optimization
**Priority**: Medium  
**Status**: Not Started

#### Tasks

- [ ] **Rate Limiting** (`internal/ai/ratelimit.go`)
  - Client-side rate limiter
  - Configurable RPM
  - Exponential backoff
  - Retry logic
  - Respect `Retry-After` header
  - Configuration
  - Unit tests

- [ ] **Concurrent Processing** (`internal/ai/concurrent.go`)
  - Optional concurrent batch processing
  - Semaphore-based limiting
  - Progress tracking for concurrent ops
  - Configuration
  - Unit tests

- [ ] **Response Caching** (`internal/ai/cache.go`)
  - Cache key generation
  - Cache storage (disk)
  - TTL management
  - Cache invalidation
  - Configuration
  - Unit tests

- [ ] **Cost Tracking** (`internal/ai/cost.go`)
  - Token counting
  - Cost estimation per model
  - Usage tracking
  - Cost warnings
  - Unit tests

- [ ] **Progress Indicators** (`internal/ai/progress.go`)
  - Progress bar for batch operations
  - Progress messages
  - JSON mode (no progress)
  - Unit tests

- [ ] **Performance Optimization**
  - Index caching to disk
  - Incremental index updates
  - Parallel file reading
  - Benchmark tests

- [ ] **Comprehensive Testing**
  - E2E tests for all commands
  - Evaluation tests with golden files
  - Quality metrics tracking
  - CI/CD integration

**Dependencies**: Phase 3C  
**Estimated Effort**: 2-3 weeks  

---

### Phase 3E: Future Enhancements
**Priority**: Low  
**Status**: Not Started

#### Tasks

- [ ] **Streaming Support**
  - Streaming provider interface
  - Streaming for `ask` command
  - Output handling for streaming
  - JSON streaming mode

- [ ] **Embeddings Integration**
  - Pre-compute embeddings
  - Vector similarity search
  - Local vector database
  - Hybrid search (embeddings + LLM)

- [ ] **Multi-Modal Support**
  - Image analysis
  - Vision-capable models
  - Diagram extraction

- [ ] **Advanced Context Strategies**
  - Recursive summarization
  - Hierarchical note organization
  - Long-context models optimization

- [ ] **Knowledge Graph**
  - Automatic link detection
  - Backlink tracking
  - Graph visualization

**Dependencies**: Phase 3D  
**Estimated Effort**: 4-6 weeks  

---

## 2. Dependencies Graph

```
Phase 3A: Infrastructure
    │
    ├── Provider Interface
    │   ├── OpenAI Provider
    │   ├── Anthropic Provider
    │   ├── Ollama Provider
    │   └── Custom Provider
    │
    ├── Configuration System
    │
    ├── Context Management
    │
    ├── Error Handling
    │
    └── Setup Command
            │
            └── Phase 3B: Core Commands
                    │
                    ├── Search Command
                    ├── Summarize Command
                    ├── Related Command
                    └── Tags Command
                            │
                            └── Phase 3C: Advanced Commands
                                    │
                                    ├── Ask Command
                                    ├── Extract Command
                                    └── Dedup Command
                                            │
                                            └── Phase 3D: Polish & Optimization
                                                    │
                                                    ├── Rate Limiting
                                                    ├── Concurrent Processing
                                                    ├── Response Caching
                                                    ├── Cost Tracking
                                                    ├── Progress Indicators
                                                    └── Performance Optimization
                                                            │
                                                            └── Phase 3E: Future Enhancements
```

---

## 3. Timeline Estimate

| Phase | Duration | Cumulative |
|-------|----------|------------|
| Phase 3A | 3-4 weeks | 3-4 weeks |
| Phase 3B | 2-3 weeks | 5-7 weeks |
| Phase 3C | 2-3 weeks | 7-10 weeks |
| Phase 3D | 2-3 weeks | 9-13 weeks |
| Phase 3E | 4-6 weeks | 13-19 weeks |

**Total Estimated Duration**: 3-4.5 months for complete AI feature set (Phases 3A-3D)

---

## 4. Risk Assessment

### High Risk

1. **Context Window Limits**
   - Risk: Users with many notes may hit token limits
   - Mitigation: Tiered retrieval, truncation strategies, user-configurable limits

2. **API Costs**
   - Risk: Users concerned about API costs
   - Mitigation: Cost tracking, warnings, local model option (Ollama), caching

3. **Provider Reliability**
   - Risk: API outages or rate limits
   - Mitigation: Retry logic, exponential backoff, clear error messages

### Medium Risk

1. **Response Quality**
   - Risk: AI responses may not meet user expectations
   - Mitigation: Extensive prompt engineering, evaluation tests, user feedback loop

2. **Performance**
   - Risk: Slow response times for large workspaces
   - Mitigation: Index caching, concurrent processing, progress indicators

### Low Risk

1. **Security**
   - Risk: API key exposure
   - Mitigation: Env var only, never store in config

---

## 5. Success Metrics

### Phase 3A
- [ ] All provider implementations complete with tests
- [ ] Setup command working end-to-end
- [ ] Configuration system robust

### Phase 3B
- [ ] 4 core commands working (search, summarize, related, tags)
- [ ] >80% test coverage for AI package
- [ ] Human and JSON output working

### Phase 3C
- [ ] 3 advanced commands working (ask, extract, dedup)
- [ ] E2E tests passing
- [ ] Quality evaluation tests passing

### Phase 3D
- [ ] Rate limiting working
- [ ] Caching working
- [ ] Performance benchmarks meeting targets:
  - Index building: <5s for 1000 notes
  - Search: <3s end-to-end
  - Summarize: <5s for single note

---

## 6. Reference Documents

- [ai-design.md](ai-design.md) - Overall design
- [ai-commands.md](ai-commands.md) - CLI specifications
- [ai-provider.md](ai-provider.md) - Provider system
- [ai-context.md](ai-context.md) - Context management
- [ai-prompts.md](ai-prompts.md) - System prompts
- [ai-testing.md](ai-testing.md) - Testing strategy

---

## 7. Notes

- Prioritize Phase 3A-3B for initial release
- Phase 3C can be released incrementally (one command at a time)
- Phase 3D can be ongoing improvements
- Phase 3E is future roadmap
- Consider beta/alpha releases for early feedback
- Documentation should be updated as features are implemented
