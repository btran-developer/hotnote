# Hotnote Documentation

Welcome to the Hotnote documentation. This is a file-based note-taking CLI application written in Go.

## Quick Links

### Getting Started
- [Quick Start](getting-started/quick-start.md) - Build, run, and common commands
- [Go Concepts](getting-started/go-concepts.md) - Go patterns explained for newcomers

### Architecture
- [Overview](architecture/overview.md) - Project structure and component relationships
- [Data Flow](architecture/data-flow.md) - How data moves through the system

### Features
- [Commands](features/commands.md) - CLI commands reference
- [Workspace](features/workspace.md) - Multi-workspace management
- [Storage](features/storage.md) - File-based storage and YAML frontmatter

## Overview

Hotnote is a terminal-based note-taking application with these key features:

- **File-based storage** - Notes stored as markdown files with YAML frontmatter
- **Multiple workspaces** - Organize notes into separate workspace directories
- **Editor integration** - Open notes in your preferred `$EDITOR`
- **Markdown rendering** - Convert notes to HTML

## Tech Stack

- **Language**: Go 1.24+
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Markdown**: [Goldmark](https://github.com/yuin/goldmark)
- **UUID**: [google/uuid](https://github.com/google/uuid)
