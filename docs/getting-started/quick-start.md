# Quick Start

## Prerequisites

- Go 1.24 or later
- `$EDITOR` environment variable set (defaults to `vim`)

## Build

```bash
go build -o hotnote ./cmd/hotnote
```

This creates the `hotnote` binary in the current directory.

## Run

```bash
./hotnote [command] [flags]
```

Or install globally:

```bash
go install ./cmd/hotnote
```

## Common Commands

### Create a new note

```bash
hotnote new "My First Note"
```

### List all notes

```bash
hotnote list
```

### Open a note for editing

```bash
hotnote open "my-first-note"
```

### Render note as HTML

```bash
hotnote render "my-first-note"
```

### Manage workspaces

```bash
hotnote workspace init          # Initialize default workspace
hotnote workspace list         # List all workspaces
hotnote workspace new notes    # Create new workspace
hotnote workspace use notes    # Switch to workspace
```

## Data Locations

| Purpose | Location |
|---------|----------|
| Workspace config | `~/.config/hotnote/config.yaml` |
| Notes storage | `~/.local/share/hotnote/workspaces/<name>/` |

## Debugging

### Exit codes

Check the exit code to understand what went wrong:

```bash
./hotnote open missing-note; echo $?
# Output: 2 (not found)
```

| Code | Meaning |
|------|---------|
| 1 | General error |
| 2 | Not found |
| 3 | Invalid input |
| 4 | Config error |

### Verbose output

Set the `$DEBUG` environment variable:

```bash
DEBUG=1 hotnote list
```

### Check workspace config

```bash
cat ~/.config/hotnote/config.yaml
```

### List workspace directories

```bash
ls -la ~/.local/share/hotnote/workspaces/
```
