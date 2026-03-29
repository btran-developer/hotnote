# AI Context Management

**Status**: Design Complete  
**Last Updated**: 2026-03-29  

## 1. Overview

Context management is the process of:
1. Selecting relevant notes for an AI operation
2. Building a context string that fits within token limits
3. Preserving the most important information

This document defines the context building strategy for all AI operations.

## 2. Context Building Pipeline

```
User Query → Fast Retrieval → Candidate Scoring → Context Assembly → Prompt Construction
```

### 2.1 Step 1: Fast Retrieval

Build a lightweight index for all notes:

```go
type NoteIndex struct {
    Slug        string
    Path        string
    Title       string
    Tags        []string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Excerpt     string    // First 500 chars of content
    WordCount   int
}
```

**Index Building**:
- Scan all `.md` files in workspace
- Extract frontmatter (title, tags, dates)
- Read first 500 characters of content
- Store in memory for fast access

**When to Rebuild**:
- On first AI command in session
- When file modification times differ from index
- Can be cached to disk for faster startup

### 2.2 Step 2: Candidate Scoring

Score each note's relevance to the query:

```go
type ScoredNote struct {
    NoteIndex
    Score float64 // 0.0 to 1.0
}

func ScoreNotes(index []NoteIndex, query string) []ScoredNote {
    var scored []ScoredNote
    
    for _, note := range index {
        score := calculateRelevance(note, query)
        scored = append(scored, ScoredNote{NoteIndex: note, Score: score})
    }
    
    // Sort by score descending
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].Score > scored[j].Score
    })
    
    return scored
}
```

**Scoring Factors**:

| Factor | Weight | Description |
|--------|--------|-------------|
| Title match | 0.30 | Query appears in title |
| Tag match | 0.25 | Query keywords match tags |
| Content match | 0.30 | Query appears in excerpt |
| Recency | 0.10 | Newer notes score higher |
| File path | 0.05 | Query matches folder name |

**Scoring Algorithm**:
```go
func calculateRelevance(note NoteIndex, query string) float64 {
    score := 0.0
    queryLower := strings.ToLower(query)
    
    // Title match
    if strings.Contains(strings.ToLower(note.Title), queryLower) {
        score += 0.30
    }
    
    // Tag match
    for _, tag := range note.Tags {
        if strings.Contains(strings.ToLower(tag), queryLower) {
            score += 0.25
            break
        }
    }
    
    // Content match
    if strings.Contains(strings.ToLower(note.Excerpt), queryLower) {
        score += 0.30
    }
    
    // Recency (exponential decay)
    daysOld := time.Since(note.UpdatedAt).Hours() / 24
    recencyScore := math.Exp(-daysOld / 30) * 0.10 // 30-day half-life
    score += recencyScore
    
    // Path match
    if strings.Contains(strings.ToLower(note.Path), queryLower) {
        score += 0.05
    }
    
    return score
}
```

### 2.3 Step 3: Context Assembly

Select top K candidates and build context string:

```go
type ContextBuilder struct {
    MaxContextNotes int    // Default: 20
    MaxTokens       int    // Default: 8000 (leaving room for system prompt + response)
    NoteCharLimit   int    // Default: 2000 chars per note (≈500 tokens)
}

func (b *ContextBuilder) BuildContext(notes []ScoredNote) (string, []string, error) {
    var included []string
    var context strings.Builder
    totalTokens := 0
    
    for i, note := range notes {
        if i >= b.MaxContextNotes {
            break
        }
        
        // Read full note content
        content, err := os.ReadFile(note.Path)
        if err != nil {
            continue // Skip unreadable notes
        }
        
        // Truncate if needed
        noteContent := string(content)
        if len(noteContent) > b.NoteCharLimit {
            noteContent = noteContent[:b.NoteCharLimit] + "... [truncated]"
        }
        
        // Calculate tokens (rough estimate)
        noteTokens := len(noteContent) / 4
        if totalTokens+noteTokens > b.MaxTokens {
            break
        }
        
        // Format note for context
        formatted := formatNote(note, noteContent)
        context.WriteString(formatted)
        context.WriteString("\n\n")
        
        included = append(included, note.Slug)
        totalTokens += noteTokens
    }
    
    return context.String(), included, nil
}
```

**Note Format**:
```xml
<note slug="projects/api-design" tags="api, architecture">
<created>2026-03-15</created>
<updated>2026-03-20</updated>
<relevance>0.92</relevance>
<content>
# API Design Decisions

Decided on REST over GraphQL for simplicity and team familiarity...
</content>
</note>
```

**Why XML format?**
- Unambiguous structure
- Easy to parse (for both AI and code)
- Handles multiline content well
- No escaping issues with markdown

### 2.4 Step 4: Prompt Construction

Combine system prompt, context, and user query:

```go
func BuildPrompt(operation string, context string, query string) (string, string) {
    systemPrompt := GetSystemPrompt(operation)
    
    userPrompt := fmt.Sprintf(`
<context>
%s
</context>

<query>
%s
</query>
`, context, query)
    
    return systemPrompt, userPrompt
}
```

## 3. Operation-Specific Context Strategies

### 3.1 search

**Strategy**: Context is the index itself (no full content needed)

```go
func BuildSearchContext(index []NoteIndex, query string) string {
    // Score and sort notes
    scored := ScoreNotes(index, query)
    
    // Include top N in context
    var context strings.Builder
    for i, note := range scored[:min(20, len(scored))] {
        context.WriteString(fmt.Sprintf("%d. %s (score: %.2f)\n", i+1, note.Slug, note.Score))
        context.WriteString(fmt.Sprintf("   Title: %s\n", note.Title))
        context.WriteString(fmt.Sprintf("   Tags: %v\n", note.Tags))
        context.WriteString(fmt.Sprintf("   Excerpt: %s\n\n", note.Excerpt))
    }
    
    return context.String()
}
```

**Output**: AI returns ranked list with scores and excerpts.

### 3.2 summarize

**Strategy**: Full content for single note, batched for multiple

```go
func BuildSummarizeContext(notes []Note, config SummarizeConfig) string {
    if len(notes) == 1 {
        // Single note: include full content
        return formatNoteFull(notes[0])
    }
    
    // Multiple notes: batch processing
    var context strings.Builder
    for i, note := range notes {
        context.WriteString(fmt.Sprintf("\n--- Note %d ---\n", i+1))
        context.WriteString(formatNoteFull(note))
    }
    
    return context.String()
}
```

**Batch Processing**:
- Group notes into batches of 5
- Process each batch separately
- Combine results

### 3.3 related

**Strategy**: Source note + candidate notes

```go
func BuildRelatedContext(source Note, candidates []NoteIndex) string {
    var context strings.Builder
    
    // Source note (full content)
    context.WriteString("<source_note>\n")
    context.WriteString(formatNoteFull(source))
    context.WriteString("</source_note>\n\n")
    
    // Candidate notes (excerpts only)
    context.WriteString("<candidate_notes>\n")
    for i, note := range candidates {
        context.WriteString(fmt.Sprintf("%d. %s\n", i+1, note.Slug))
        context.WriteString(fmt.Sprintf("   Title: %s\n", note.Title))
        context.WriteString(fmt.Sprintf("   Tags: %v\n", note.Tags))
        context.WriteString(fmt.Sprintf("   Excerpt: %s\n\n", note.Excerpt))
    }
    context.WriteString("</candidate_notes>")
    
    return context.String()
}
```

### 3.4 tags

**Strategy**: Full content of single note

```go
func BuildTagsContext(note Note) string {
    return formatNoteFull(note)
}
```

### 3.5 ask (Q&A)

**Strategy**: Most context-heavy operation

```go
func BuildAskContext(query string, notes []ScoredNote, config AskConfig) string {
    // For Q&A, we want high-quality context
    // Include full content of top notes, excerpts of others
    
    var context strings.Builder
    fullContentLimit := 5 // Full content for top 5 notes
    
    for i, note := range notes {
        if i >= config.MaxContextNotes {
            break
        }
        
        if i < fullContentLimit {
            // Full content for top notes
            content, _ := os.ReadFile(note.Path)
            context.WriteString(formatNoteWithFullContent(note, string(content)))
        } else {
            // Excerpt for remaining notes
            context.WriteString(formatNoteWithExcerpt(note))
        }
        context.WriteString("\n\n")
    }
    
    return context.String()
}
```

### 3.6 extract

**Strategy**: Similar to ask, but focus on specific passages

```go
func BuildExtractContext(query string, notes []ScoredNote) string {
    // Similar to ask context
    // AI will identify and extract relevant passages
    return BuildAskContext(query, notes, AskConfig{})
}
```

### 3.7 dedup

**Strategy**: Compare notes pairwise

```go
func BuildDedupContext(note1, note2 Note) string {
    content1, _ := os.ReadFile(note1.Path)
    content2, _ := os.ReadFile(note2.Path)
    
    return fmt.Sprintf(`
<note1 slug="%s">
%s
</note1>

<note2 slug="%s">
%s
</note2>
`, note1.Slug, content1, note2.Slug, content2)
}
```

**Optimization**: Don't compare all pairs (O(n²))
- Use scoring to find likely duplicates first
- Only compare top candidates
- Hash-based pre-filtering (same title, similar length)

## 4. Token Management

### 4.1 Token Estimation

```go
func EstimateTokens(text string) int {
    // Rough estimate: 1 token ≈ 4 characters
    return len(text) / 4
}

func EstimateNoteTokens(note Note) int {
    content, _ := os.ReadFile(note.Path)
    return EstimateTokens(string(content))
}
```

### 4.2 Token Budget Allocation

For a context window of 8000 tokens:

| Component | Allocation | Description |
|-----------|------------|-------------|
| System prompt | 500 tokens | Fixed cost |
| User query | 200 tokens | Variable |
| Context | 6000 tokens | Note content |
| Response | 1000 tokens | AI output buffer |
| Safety margin | 300 tokens | Buffer |

### 4.3 Truncation Strategies

When content exceeds budget:

1. **Head-only**: First N characters
2. **Head-tail**: First N/2 + Last N/2 (for code with imports/exports)
3. **Smart truncation**: Preserve headings, truncate body paragraphs

```go
func TruncateNote(content string, maxChars int) string {
    if len(content) <= maxChars {
        return content
    }
    
    // Try smart truncation first
    if strings.Contains(content, "#") {
        return truncateWithHeadings(content, maxChars)
    }
    
    // Default: head-only
    return content[:maxChars] + "\n\n... [truncated]"
}

func truncateWithHeadings(content string, maxChars int) string {
    // Preserve all headings
    // Truncate content between headings
    // Return: headings + truncated content
}
```

## 5. Caching Strategies

### 5.1 Context Cache

Cache the built context to avoid re-reading files:

```go
type ContextCache struct {
    cache map[string]CachedContext
    ttl   time.Duration
}

type CachedContext struct {
    Context      string
    IncludedSlugs []string
    Timestamp    time.Time
}

func (c *ContextCache) Get(key string) (CachedContext, bool) {
    cached, ok := c.cache[key]
    if !ok {
        return CachedContext{}, false
    }
    
    if time.Since(cached.Timestamp) > c.ttl {
        delete(c.cache, key)
        return CachedContext{}, false
    }
    
    return cached, true
}
```

**Cache Key**: hash(query + operation + workspace_hash)

### 5.2 Note Index Cache

Persist note index to disk for faster startup:

```go
func SaveIndex(index []NoteIndex, path string) error {
    data, _ := json.Marshal(index)
    return os.WriteFile(path, data, 0644)
}

func LoadIndex(path string) ([]NoteIndex, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var index []NoteIndex
    err = json.Unmarshal(data, &index)
    return index, err
}
```

**Invalidation**: Rebuild when file modification times differ from cached times.

## 6. User Control

### 6.1 Flags

```bash
--context-limit N      # Max notes to include (default: 20)
--context-tokens N     # Max tokens for context (default: 6000)
--recent N             # Only include notes from last N days
```

### 6.2 Configuration

```yaml
ai:
  context:
    max_notes: 20
    max_tokens: 6000
    char_limit_per_note: 2000
```

## 7. Performance Considerations

### 7.1 Index Building

- Build index asynchronously on startup
- Use goroutines to read files in parallel
- Limit concurrent file reads (e.g., max 10)

### 7.2 Lazy Loading

- Don't build index until first AI command
- Build index in background while user interacts

### 7.3 Incremental Updates

- Only re-read changed files
- Track file modification times
- Update index incrementally

## 8. Error Handling

### 8.1 File Reading Errors

```go
func safeReadFile(path string) ([]byte, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        log.Printf("Warning: Could not read %s: %v", path, err)
        return nil, err
    }
    return content, nil
}
```

**Strategy**: Skip unreadable files, continue with available notes.

### 8.2 Context Too Large

```go
if tokenCount > maxTokens {
    return "", ErrAIContextTooLarge
}
```

**Recovery**: Reduce `--context-limit` or `--context-tokens`.

## 9. Testing

### 9.1 Context Building Tests

```go
func TestContextBuilder(t *testing.T) {
    notes := []ScoredNote{
        {Slug: "note1", Score: 0.9, Path: "note1.md"},
        {Slug: "note2", Score: 0.8, Path: "note2.md"},
    }
    
    builder := ContextBuilder{
        MaxContextNotes: 2,
        MaxTokens: 1000,
    }
    
    context, included, err := builder.BuildContext(notes)
    assert.NoError(t, err)
    assert.Len(t, included, 2)
    assert.Contains(t, context, "note1")
}
```

### 9.2 Token Estimation Tests

```go
func TestTokenEstimation(t *testing.T) {
    text := "Hello world, this is a test."
    tokens := EstimateTokens(text)
    // Rough check: should be around 6-8 tokens
    assert.True(t, tokens >= 5 && tokens <= 10)
}
```

### 9.3 Truncation Tests

```go
func TestTruncateNote(t *testing.T) {
    longContent := strings.Repeat("a", 10000)
    truncated := TruncateNote(longContent, 1000)
    assert.LessOrEqual(t, len(truncated), 1100)
    assert.Contains(t, truncated, "[truncated]")
}
```
