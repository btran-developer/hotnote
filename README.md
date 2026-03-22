# Hotnote - Terminal-Based Markdown Note App

A terminal-first, markdown-native knowledge system designed for both humans and AI agents.

## Prerequisites
- Go 1.22+
- Basic understanding of command-line interfaces

## Installation
1. Clone this repository
2. Build the application:
```bash   
go build -o hotnote ./cmd/hotnote
```

## Usage
### Create a new note
```bash
./hotnote new "My idea"
```

### List all notes
```bash
./hotnote list
```

### Open a note for editing
```bash
./hotnote open "My idea"
```

### Render markdown to HTML
```bash
./hotnote render "My idea"
```

### AI-powered operations
```bash
./hotnote ai write --title "Research Plan" --input file.txt
./hotnote ai summarize my-note./hotnote ai query "distributed systems"
```

## Exit Codes

Hotnote uses exit codes to indicate the result of command execution:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Not found |
| 3 | Invalid input |
| 4 | Configuration error |

Example:
```bash
./hotnote open nonexistent; echo $?  # Outputs: 2
```

## Data Storage
- Notes are stored as `.md` files in the `notes/` directory by default
- Each note uses a filename-based slug for its ID
- Optional frontmatter can be added to notes for metadata

## Project Structure
```
hotnote/
├── cmd/
│   └── hotnote/          # Main command implementation
├── internal/
│   ├── core/             # Domain models
│   ├── storage/          # Data storage implementation
│   ├── markdown/         # Markdown processing
│   ├── search/           # Search functionality
│   ├── ai/               # AI integration
│   ├── tui/              # Terminal UI (planned)
│   └── cli/              # CLI utilities
├── go.mod                # Go module definition
└── README.md             # This file
```

## Development Roadmap- Phase 1: MVP CLI (implemented)
- Phase 2: TUI interface- Phase 3: Editor improvements
- Phase 4: AI interface
- Phase 5: Search & knowledge system

## License
MIT License