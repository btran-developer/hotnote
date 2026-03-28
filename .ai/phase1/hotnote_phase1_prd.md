# HotNote — Phase 1 PRD (MVP)

## 1. Objective

Build a minimal, reliable CLI-based note system that:

- Stores notes as Markdown files
- Supports basic CRUD operations
- Organizes notes via folders/workspaces
- Is usable daily by developers
- Forms the foundation for future TUI and AI features

---

## 2. Non-Goals (Phase 1 Scope Control)

Phase 1 will NOT include:

- TUI interface
- Built-in text editor (use $EDITOR)
- AI features
- Full-text search
- Tag-based querying
- Backlinks / graph visualization

---

## 3. Target Users

Primary:
- Developers using terminal workflows

Secondary (future):
- AI agents interacting via CLI

---

## 4. Core Use Cases

Create a note:
    hotnote new "My Idea"
    hotnote new "My Idea" --path projects    # In subfolder

List notes (including subfolders):
    hotnote list
    hotnote list --json

Open/edit a note:
    hotnote open my-idea
    hotnote open projects/my-idea           # Direct path

Render note:
    hotnote render my-idea
    hotnote render projects/my-idea          # Direct path

Create folder:
    hotnote folder create projects
    hotnote folder create projects/2024
    hotnote folder cr projects                    # alias

Delete note:
    hotnote delete my-idea
    hotnote delete my-idea --force
    hotnote del my-idea --force                  # alias

Delete folder:
    hotnote folder delete projects
    hotnote folder delete projects --force
    hotnote folder del projects --force            # alias

Workspace management:
    hotnote workspace init
    hotnote workspace list
    hotnote workspace use <name>
    hotnote workspace new <name>
    hotnote workspace delete <name>
    hotnote workspace delete <name> --force

---

## 5. Functional Requirements

### Note Creation

- Generate slug (lowercase, hyphen-separated, no special chars)
- Generate UUID
- Open in $EDITOR

Example:
"My Research Plan!" → my-research-plan

File format:

---
id: <uuid>
title: My Idea
created_at: <timestamp>
updated_at: <timestamp>
tags: []
---

# My Idea

---

### Storage

Default:
~/.local/share/hotnote/workspaces/default/

Structure:

workspace/
  note-1.md
  project/
    note-2.md

---

### Listing Notes

Default output:

project-a/my-idea        2026-03-21
project-a/design         2026-03-20
journal                  2026-03-18

Flags:
    --flat
    --sort updated
    --json

JSON output:

[
  {
    "slug": "my-idea",
    "path": "project-a/my-idea.md",
    "updated_at": "2026-03-21T10:00:00Z"
  }
]

---

### Open Note

Uses $EDITOR, fallback to vim or nano

---

### Workspace

Config:
~/.config/hotnote/config.yaml

Data:
~/.local/share/hotnote/workspaces/

Example config:

current_workspace: default

workspaces:
  default: ~/.local/share/hotnote/workspaces/default
  work: ~/notes/work

---

## 6. CLI Structure

hotnote
  new <title>
  list
  open <slug>
  render <slug>
  workspace
    init
    list
    use <name>

---

## 7. Non-Functional Requirements

- Fast (<100ms)
- Reliable (no data loss)
- Simple (no DB)
- Portable (macOS + Linux)
- Deterministic behavior

---

## 8. Edge Cases

- Slug collision:
  my-idea.md
  my-idea-1.md

- Missing $EDITOR
- Non-existent note
- Corrupted frontmatter

---

## 9. Success Criteria

- Replaces ad-hoc markdown notes
- Used daily
- Fast and natural workflow

---

## 10. Implementation Plan

Week 1:
- CLI (Cobra)
- Storage
- new + list

Week 2:
- open
- render
- workspace

---

## 11. Principles

- Markdown-first
- Simple
- Human + AI friendly
- Ship early

---

## 12. Future

- TUI
- AI
- Search
- Tags
- Graph
