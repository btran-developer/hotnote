# HotNote CLI — AI Agent Guide

## Overview

The HotNote CLI is designed to work with both human users and AI agents.
This guide covers machine-readable interfaces, error handling, and best practices
for building AI agents that interact with HotNote.

---

## Quick Start

```bash
# Initialize workspace
hotnote workspace init

# Create a note
hotnote new "My Note" --json

# List notes
hotnote list --json

# Delete a note
hotnote delete my-note --force --json
```

---

## JSON Output Mode

All commands support `--json` flag for machine-readable output.

### Success Responses

```bash
# Create note
hotnote new "My Idea" --json
# {"status":"created","slug":"my-idea","path":"/Users/.../workspaces/default/my-idea.md"}

# List notes
hotnote list --json
# [{"slug":"my-idea","path":"/Users/.../my-idea.md","updated_at":"2026-03-25T10:00:00Z"}]

# Open note
hotnote open my-idea --json
# {"status":"opened","path":"/Users/.../my-idea.md"}

# Delete note
hotnote delete my-idea --force --json
# {"status":"deleted","slug":"my-idea"}
```

### Error Responses

```json
{"error": "note not found"}
{"error": "workspace not initialized"}
{"error": "title is required"}
```

---

## Exit Codes

AI agents should check exit codes to determine success or failure type.

| Code | Name | Meaning | Example |
|------|------|---------|---------|
| 0 | Success | Command completed successfully | Normal operation |
| 1 | General error | Unexpected system error | File system error, permission denied |
| 2 | Not found | Resource not found | Note, folder, or workspace doesn't exist |
| 3 | Invalid input | Bad user input | Missing required arguments, duplicate names |
| 4 | Config error | Configuration issue | Workspace not initialized, corrupted config |

### Exit Code Usage

```bash
hotnote new "My Idea"
if [ $? -eq 0 ]; then
    echo "Success"
fi

hotnote open my-idea
if [ $? -eq 2 ]; then
    echo "Note not found"
fi
```

---

## Non-Interactive Mode

When `--json` is used, interactive prompts are suppressed.
For destructive operations, use `--force` to skip confirmation prompts.

### Examples

```bash
# Create without opening editor (recommended for AI agents)
# (--no-open flag to be added)

# Delete without prompt
hotnote delete my-idea --force --json

# Delete workspace without prompt
hotnote workspace delete work --force --json
```

---

## Path Resolution

The CLI supports two path resolution modes:

### Direct Path
If the argument contains `/`, it's treated as a direct relative path:

```bash
hotnote open projects/my-idea
hotnote render notes/tasks.md
hotnote delete drafts/idea
```

### Recursive Search
If the argument doesn't contain `/`, it searches recursively by slug:

```bash
hotnote open my-idea           # Searches all subfolders
hotnote render notes           # Searches for "notes" slug
```

### Conflict Resolution
If multiple notes match the slug, an error is returned:

```json
{"error": "multiple matches: my-idea found in ['my-idea.md', 'projects/my-idea.md']. Specify full path."}
```

---

## Deterministic Output

AI agents can rely on the following guarantees:

1. **JSON output is deterministic** - Same input produces same output structure
2. **Error messages are stable strings** - No system-specific details in errors
3. **File paths are absolute** - Full paths returned for clarity
4. **Timestamps in RFC3339** - Standard format for parsing

---

## AI Agent Patterns

### Pattern 1: Create and Track

```bash
# Create note and capture path
NOTE_PATH=$(hotnote new "Meeting Notes" --json | jq -r '.path')
# Use NOTE_PATH for subsequent operations
```

### Pattern 2: Find and Process

```bash
# Find all notes in JSON
NOTES=$(hotnote list --json)

# Find specific note
NOTE=$(echo "$NOTES" | jq -r '.[] | select(.slug == "my-idea")')
```

### Pattern 3: Safe Delete

```bash
# Check if note exists, then delete
hotnote open my-idea 2>/dev/null
if [ $? -eq 0 ]; then
    hotnote delete my-idea --force --json
fi
```

### Pattern 4: Workspace Management

```bash
# List workspaces
hotnote workspace list --json

# Switch workspace
hotnote workspace use work --json

# Create workspace
hotnote workspace new my-workspace --json
```

---

## Best Practices

1. **Always use `--json`** for parsing output programmatically
2. **Check exit codes** before processing output
3. **Use `--force`** for batch operations to avoid prompts
4. **Handle "multiple matches"** errors by specifying full path
5. **Use absolute paths** from JSON output for subsequent operations
6. **Initialize workspace** before running other commands
7. **Handle errors gracefully** - don't assume operations succeed

---

## Command Reference

| Command | JSON Support | --force Support | Notes |
|---------|--------------|-----------------|-------|
| create | ✓ | N/A | Creates note |
| list | ✓ | N/A | Lists all notes |
| open | ✓ | N/A | Opens in $EDITOR |
| render | ✓ | N/A | Renders to HTML |
| delete | ✓ | ✓ | Deletes note |
| folder create | ✓ | N/A | Creates folder |
| folder delete | ✓ | ✓ | Deletes folder |
| folder list | ✓ | N/A | Lists folder contents |
| workspace init | ✓ | N/A | Initializes workspace |
| workspace create | ✓ | N/A | Creates workspace |
| workspace list | ✓ | N/A | Lists workspaces |
| workspace use | ✓ | N/A | Switches workspace |
| workspace delete | ✓ | ✓ | Deletes workspace |

---

## Error Handling

### Common Errors

| Error | Exit Code | Cause | Resolution |
|-------|-----------|-------|------------|
| workspace not initialized | 4 | No config | Run `workspace init` |
| note not found | 2 | File doesn't exist | Check slug/path |
| multiple matches | 2 | Slug exists in multiple folders | Use full path |
| title is required | 3 | Missing argument | Provide title |
| workspace not found | 2 | Workspace name invalid | Check `workspace list` |
| cannot delete current workspace | 3 | Trying to delete active | Switch workspace first |

### Handling Errors in Scripts

```bash
#!/bin/bash

# Function to handle errors
handle_error() {
    local exit_code=$1
    case $exit_code in
        1) echo "System error" ;;
        2) echo "Resource not found" ;;
        3) echo "Invalid input" ;;
        4) echo "Workspace not initialized" ;;
    esac
}

hotnote new "My Note" --json
exit_code=$?
if [ $exit_code -ne 0 ]; then
    handle_error $exit_code
fi
```

---

## Security Notes

1. **No path traversal protection** - Be careful with `../` in paths
2. **Workspace deletion is recursive** - `workspace delete` removes all files
3. **No undo** - Deleted notes and workspaces cannot be recovered
4. **File permissions** - Uses 0644 for notes, 0755 for directories

---

## Future Considerations

- **AI-specific flags** - Planned for batch operations
- **Structured metadata** - For better note indexing
- **Search API** - Full-text search support
- **Tag filtering** - Query by tags