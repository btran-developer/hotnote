# AI Commands Specification

**Status**: Design Complete  
**Last Updated**: 2026-03-29  

## 1. Command Overview

```
hotnote ai <subcommand> [args] [flags]
```

## 2. Global Flags

All AI subcommands support these flags:

```
--json              Machine-readable JSON output
--provider          Override configured provider (openai | anthropic | ollama | custom)
--model             Override configured model
--max-tokens        Override max tokens (default: 4096)
--context-limit     Max notes to include in context (default: 20)
--concurrent        Enable concurrent processing (default: false)
--no-cache          Skip cache lookup (when caching enabled)
```

## 3. Subcommands

### 3.1 search

**Purpose**: Semantic search across notes using natural language queries.

**Usage**:
```bash
hotnote ai search <query> [flags]
```

**Flags**:
```
--workspace         Search specific workspace
--folder            Search within folder
--recent            Limit to notes from last N days
--limit             Max results (default: 10)
```

**Human Output**:
```
$ hotnote ai search "API design decisions"

Found 3 relevant notes:

1. projects/api-design.md (score: 0.92)
   "Decided on REST over GraphQL for simplicity..."

2. meetings/q4-planning.md (score: 0.78)
   "API redesign scheduled for Q4..."

3. notes/architecture.md (score: 0.71)
   "API follows resource-oriented design..."

Usage: 1,700 tokens (~$0.003)
```

**JSON Output**:
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

**Error Cases**:
- No results: Empty results array, success exit code
- Provider not configured: Exit code 5
- Request failed: Exit code 6

---

### 3.2 summarize

**Purpose**: Summarize a note or multiple notes.

**Usage**:
```bash
hotnote ai summarize <note> [flags]
hotnote ai summarize <folder> --all [flags]
```

**Flags**:
```
--all               Summarize all notes in folder/workspace
--recent            Summarize notes from last N days
--format            Output format: bullet | paragraph | outline (default: bullet)
--include-sources   Include source note paths in output
```

**Human Output (single note)**:
```
$ hotnote ai summarize projects/api-design.md

Summary:
API design decisions from March 2026. Chose REST over GraphQL for 
simplicity, decided on JSON:API spec, and established versioning strategy.

Key Points:
• REST selected for team familiarity and tooling
• GraphQL deferred to future if needed
• JSON:API spec for consistency
• URL versioning (/v1/, /v2/)
• Deprecation policy: 6 months notice

Topics: api, architecture, rest, graphql

Source: projects/api-design.md
Usage: 850 tokens (~$0.001)
```

**JSON Output (single note)**:
```json
{
  "slug": "projects/api-design",
  "path": "/Users/.../projects/api-design.md",
  "summary": "API design decisions from March 2026...",
  "key_points": [
    "REST selected for team familiarity and tooling",
    "GraphQL deferred to future if needed",
    "JSON:API spec for consistency"
  ],
  "topics": ["api", "architecture", "rest", "graphql"],
  "metadata": {
    "original_length": 2500,
    "summary_length": 180
  },
  "usage": {
    "prompt_tokens": 700,
    "completion_tokens": 150,
    "total_tokens": 850,
    "estimated_cost": 0.001
  }
}
```

**Human Output (batch)**:
```
$ hotnote ai summarize folder/meetings --all

Summarizing 23 notes (5 batches)...

Batch 1/5: ✓ 5 notes summarized
Batch 2/5: ✓ 5 notes summarized
Batch 3/5: ✓ 5 notes summarized
Batch 4/5: ✓ 5 notes summarized
Batch 5/5: ✓ 3 notes summarized

Done. 23 notes summarized.

Common Themes:
• Q4 planning and roadmap discussions
• API and architecture decisions
• Team process improvements

Individual Summaries:
1. 2026-03-15.md - Sprint planning, focused on performance...
2. 2026-03-22.md - API design decisions, chose REST...
...

Usage: 4,200 tokens (~$0.006)
```

**JSON Output (batch)**:
```json
{
  "mode": "batch",
  "total_notes": 23,
  "summaries": [
    {
      "slug": "meetings/2026-03-15",
      "summary": "Sprint planning...",
      "topics": ["sprint", "planning"]
    }
  ],
  "common_themes": ["Q4 planning", "API decisions", "process improvements"],
  "usage": {
    "total_tokens": 4200,
    "estimated_cost": 0.006
  }
}
```

---

### 3.3 related

**Purpose**: Find notes related to a given note.

**Usage**:
```bash
hotnote ai related <note> [flags]
```

**Flags**:
```
--limit             Max related notes (default: 5)
--explain           Include explanation of why related
--threshold         Minimum similarity score (default: 0.6)
```

**Human Output**:
```
$ hotnote ai related projects/api-design.md

Found 4 related notes:

1. meetings/2026-03-22.md (score: 0.89)
   Related: Documents the meeting where API design decisions were made

2. notes/architecture.md (score: 0.82)
   Related: Discusses broader architecture including API layer

3. projects/graphql-evaluation.md (score: 0.76)
   Related: Evaluates GraphQL vs REST, referenced in API decision

4. notes/rest-best-practices.md (score: 0.71)
   Related: REST guidelines that informed API design

Usage: 1,200 tokens (~$0.002)
```

**JSON Output**:
```json
{
  "source_slug": "projects/api-design",
  "related": [
    {
      "slug": "meetings/2026-03-22",
      "path": "/Users/.../meetings/2026-03-22.md",
      "score": 0.89,
      "explanation": "Documents the meeting where API design decisions were made",
      "shared_topics": ["api", "rest", "decisions"]
    }
  ],
  "usage": {
    "total_tokens": 1200,
    "estimated_cost": 0.002
  }
}
```

---

### 3.4 tags

**Purpose**: Suggest or apply tags for a note.

**Usage**:
```bash
hotnote ai tags <note> [flags]
```

**Flags**:
```
--apply             Apply suggested tags to note
--dry-run           Preview without applying (default when no --apply)
--confidence        Minimum confidence threshold (0-1, default: 0.7)
--all               Tag all notes in folder/workspace
```

**Human Output (dry run)**:
```
$ hotnote ai tags projects/api-design.md

Suggested Tags:

• api          (confidence: 0.96)
• rest         (confidence: 0.94)
• architecture (confidence: 0.88)
• graphql      (confidence: 0.82)
• json-api     (confidence: 0.79)

Run with --apply to add these tags to the note.

Usage: 650 tokens (~$0.001)
```

**Human Output (with apply)**:
```
$ hotnote ai tags projects/api-design.md --apply

Applying 5 tags to projects/api-design.md:
✓ api
✓ rest
✓ architecture
✓ graphql
✓ json-api

Note updated successfully.
Usage: 650 tokens (~$0.001)
```

**JSON Output (dry run)**:
```json
{
  "slug": "projects/api-design",
  "suggested_tags": [
    {"tag": "api", "confidence": 0.96},
    {"tag": "rest", "confidence": 0.94},
    {"tag": "architecture", "confidence": 0.88}
  ],
  "applied": false,
  "usage": {
    "total_tokens": 650,
    "estimated_cost": 0.001
  }
}
```

**JSON Output (with apply)**:
```json
{
  "slug": "projects/api-design",
  "applied_tags": ["api", "rest", "architecture", "graphql", "json-api"],
  "applied": true,
  "usage": {
    "total_tokens": 650,
    "estimated_cost": 0.001
  }
}
```

---

### 3.5 ask

**Purpose**: Ask questions about notes and get answers with citations.

**Usage**:
```bash
hotnote ai ask <question> [flags]
```

**Flags**:
```
--sources           Include source citations (default: true)
--recent            Limit context to notes from last N days
--notes             Specific notes to search (comma-separated slugs)
--workspace         Search specific workspace
```

**Human Output**:
```
$ hotnote ai ask "What did I decide about the API design?"

Answer:
You decided to use REST over GraphQL for the API design, primarily 
for simplicity and team familiarity. The JSON:API specification was 
chosen for consistency, and URL versioning (/v1/, /v2/) was established.

Sources:
1. projects/api-design.md
   Excerpt: "Decided on REST over GraphQL for simplicity and team familiarity"
   Relevance: Direct decision statement

2. meetings/2026-03-22.md
   Excerpt: "Team voted 4-1 in favor of REST"
   Relevance: Meeting decision record

Usage: 2,100 tokens (~$0.003)
```

**JSON Output**:
```json
{
  "question": "What did I decide about the API design?",
  "answer": "You decided to use REST over GraphQL...",
  "citations": [
    {
      "slug": "projects/api-design",
      "path": "/Users/.../projects/api-design.md",
      "excerpt": "Decided on REST over GraphQL for simplicity and team familiarity",
      "relevance": "Direct decision statement"
    }
  ],
  "usage": {
    "prompt_tokens": 1800,
    "completion_tokens": 300,
    "total_tokens": 2100,
    "estimated_cost": 0.003
  }
}
```

**Error Cases**:
- No relevant notes found: "I don't have information about that in your notes."
- Provider not configured: Exit code 5

---

### 3.6 extract

**Purpose**: Extract specific passages relevant to a query.

**Usage**:
```bash
hotnote ai extract <query> [flags]
```

**Flags**:
```
--format            Output format: quotes | bullet | json (default: bullet)
--notes             Specific notes to search (comma-separated)
--workspace         Search specific workspace
--folder            Search within folder
--recent            Limit to notes from last N days
```

**Human Output**:
```
$ hotnote ai extract "API design decisions" --format bullet

Found 4 relevant passages:

• "Decided on REST over GraphQL for simplicity and team familiarity" 
  — projects/api-design.md

• "Team voted 4-1 in favor of REST. Main concerns: learning curve, 
   tooling complexity" 
  — meetings/2026-03-22.md

• "GraphQL evaluation postponed to Q4 if REST proves insufficient"
  — projects/api-design.md

• "JSON:API spec selected for consistency across endpoints"
  — notes/rest-best-practices.md

Usage: 1,400 tokens (~$0.002)
```

**JSON Output**:
```json
{
  "query": "API design decisions",
  "extractions": [
    {
      "passage": "Decided on REST over GraphQL for simplicity and team familiarity",
      "slug": "projects/api-design",
      "path": "/Users/.../projects/api-design.md",
      "relevance_score": 0.94,
      "context": "..."
    }
  ],
  "usage": {
    "total_tokens": 1400,
    "estimated_cost": 0.002
  }
}
```

---

### 3.7 dedup

**Purpose**: Find duplicate or highly similar notes.

**Usage**:
```bash
hotnote ai dedup [flags]
```

**Flags**:
```
--scan              Scan entire workspace for duplicates
--notes             Specific notes to compare (comma-separated)
--threshold         Similarity threshold (0-1, default: 0.8)
--merge-suggest     Suggest which note to keep when duplicates found
--folder            Limit scan to folder
--recent            Only check notes from last N days
```

**Human Output**:
```
$ hotnote ai dedup --scan

Scanning 156 notes for duplicates...

Found 2 potential duplicates:

1. Similarity: 0.87
   • notes/ideas/api-v2.md
   • projects/api-design.md
   
   These notes cover similar API design topics. The newer note 
   (api-v2.md) appears to be an update/refinement.
   
   Suggested action: Review and consider merging or linking.

2. Similarity: 0.82
   • meetings/2026-03-15.md
   • meetings/2026-03-15-backup.md
   
   These appear to be duplicate meeting notes (possibly a backup).
   
   Suggested action: Delete the backup copy.

Usage: 3,800 tokens (~$0.005)
```

**JSON Output**:
```json
{
  "scanned_notes": 156,
  "duplicates_found": 2,
  "duplicates": [
    {
      "notes": [
        {
          "slug": "notes/ideas/api-v2",
          "path": "/Users/.../notes/ideas/api-v2.md"
        },
        {
          "slug": "projects/api-design",
          "path": "/Users/.../projects/api-design.md"
        }
      ],
      "similarity": 0.87,
      "explanation": "These notes cover similar API design topics...",
      "suggested_action": "Review and consider merging or linking"
    }
  ],
  "usage": {
    "total_tokens": 3800,
    "estimated_cost": 0.005
  }
}
```

---

### 3.8 setup

**Purpose**: Configure AI provider interactively.

**Usage**:
```bash
hotnote ai setup [flags]
```

**Flags**:
```
--provider          Provider to configure (skip interactive selection)
--model             Model to use
--api-key-env       Environment variable name for API key
--base-url          Custom endpoint URL
--force             Overwrite existing configuration
```

**Interactive Flow**:
```
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

**Ollama Specific**:
```
? Choose a provider: Ollama (local)

? Ollama endpoint (default: http://localhost:11434): http://localhost:11434
? Model (available: llama3, mistral, codellama): llama3

✓ Configuration saved
✓ Test request successful (120ms, 150 tokens)
⚠ Using local model. Responses may be slower but data stays on your machine.
```

**Custom Provider**:
```
? Choose a provider: Custom (OpenAI-compatible)

? Endpoint URL: https://api.together.xyz/v1
? API key environment variable: TOGETHER_API_KEY
? Model: meta-llama/Llama-3-70b-chat-hf

✓ Configuration saved
✓ Test request successful (200ms, 150 tokens)
```

**JSON Output**:
```json
{
  "status": "configured",
  "provider": "openai",
  "model": "gpt-4o-mini",
  "api_key_env": "OPENAI_API_KEY",
  "test_result": {
    "success": true,
    "latency_ms": 45,
    "tokens": 150
  }
}
```

**Error Cases**:
- API key not set: Show which env var is expected
- Test request fails: Show error details, suggest checking configuration
- No network: Suggest checking connection or using Ollama (local)

---

## 4. Output Schema Reference

### 4.1 Usage Object

All commands include usage information:

```json
{
  "usage": {
    "prompt_tokens": 1500,
    "completion_tokens": 200,
    "total_tokens": 1700,
    "estimated_cost": 0.003
  }
}
```

### 4.2 Metadata Object

Notes include metadata where relevant:

```json
{
  "metadata": {
    "tags": ["api", "architecture"],
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-20T14:30:00Z"
  }
}
```

### 4.3 Citation Object

Used in Q&A and extract commands:

```json
{
  "citations": [
    {
      "slug": "projects/api-design",
      "path": "/Users/.../projects/api-design.md",
      "excerpt": "Decided on REST over GraphQL...",
      "relevance": "Direct decision statement"
    }
  ]
}
```

## 5. Progress Indicators

For long operations (batch operations, dedup), show progress:

```
Summarizing 23 notes (5 batches)...

Batch 1/5: ✓ 5 notes summarized
Batch 2/5: ✓ 5 notes summarized
...
```

Or for non-batch long operations:
```
Scanning 156 notes for duplicates...
[████████░░░░░░░░░░░░] 45%
```

**JSON Mode**: No progress indicators, complete response at end.
