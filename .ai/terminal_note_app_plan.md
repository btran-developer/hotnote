# Terminal-Based Markdown Note App (Go) — Development Plan

## Vision
A terminal-first, markdown-native knowledge system designed for both humans and AI agents.

- Human UX: TUI editor + explorer
- Machine UX: CLI with structured I/O
- Content: Markdown as the source of truth

---

## Product Direction

**Approach: Note-system first (recommended)**

Focus on:
- Strong storage model
- Clean CLI interface
- AI-friendly design

Editing experience can evolve later.

---

## Tech Stack

### Core
- Language: Go
- Storage: Filesystem (.md files)
- Config: YAML or TOML

### TUI
- Bubble Tea
- Lip Gloss

### CLI
- Cobra

### Markdown
- goldmark

---

## Architecture

```
/cmd
  note-cli
  note-tui

/internal
  /core
  /storage
  /markdown
  /search
  /ai
  /tui
  /cli
```

### Core Domain Models

```go
type Note struct {
    ID        string
    Title     string
    Path      string
    Tags      []string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Workspace struct {
    RootPath string
}
```

---

## Development Phases

### Phase 1 — MVP (2–3 weeks)

Goal: Usable CLI-based note system

Features:
- Create note
- Edit note (via $EDITOR)
- List notes
- Workspace/folder support
- Markdown rendering

Example CLI:

```
note new "My idea"
note list
note open my-idea
note render my-idea
```

---

### Phase 2 — TUI

Goal: Terminal UI for navigation

Features:
- Folder tree (left)
- Note list (center)
- Preview/editor (right)
- Toggle raw ↔ rendered markdown
- Keyboard navigation

---

### Phase 3 — Editor Improvements

Features:
- Syntax highlighting
- Basic editing capabilities
- Markdown-aware behavior

Note: Full editor is complex — keep it simple initially.

---

### Phase 4 — AI Interface

Goal: Enable AI interaction via CLI

Examples:

```
note ai write --title "Research Plan" --input file.txt
note ai summarize my-note
note ai query "distributed systems"
```

Design principle:
- Output should be machine-readable (JSON)

---

### Phase 5 — Search & Knowledge System

Features:
- Full-text search
- Tagging
- Backlinks ([[note]])
- Graph (optional)

---

## Key Design Decisions

### File Structure

```
notes/
  project-a/
    idea.md
```

Optional frontmatter:

```yaml
---
title: Idea
tags: [ai, system]
---
```

---

### ID Strategy

- Filename: slug (human-readable)
- Internal ID: UUID

---

### AI Integration

```go
type AIProvider interface {
    Generate(prompt string) (string, error)
}
```

---

## Execution Plan

### Week 1
- CLI setup (Cobra)
- Note CRUD
- File storage

### Week 2
- Markdown rendering
- Workspace support
- CLI improvements

### Week 3
- Basic TUI (list + preview)

### Week 4+
- Editor improvements
- AI features

---

## Guiding Principles

- Keep everything markdown-first
- Make it useful early
- Avoid over-engineering
- Build for both humans and AI

---

## Long-Term Potential

- AI-native knowledge base
- Developer research notebook
- CLI-based Notion alternative
- Agent memory system
