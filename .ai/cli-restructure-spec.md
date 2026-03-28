# CLI Command Restructure Specification

## Overview

Restructure folder commands under a unified `folder` subcommand and add short aliases for better consistency, usability, and AI agent compatibility.

**Status:** Planned (not yet implemented)
**Issue:** Issue O (CLI Command Restructure), Issue P (Add folder list)

## Breaking Changes

This is a breaking change since the app is still in early development (Phase 1).

| Old Command | New Command | Alias |
|--------------|-------------|-------|
| `hotnote mkdir <path>` | `hotnote folder create <path>` | `hotnote folder new <path>`, `hotnote folder cr <path>` |
| `hotnote rmdir <path>` | `hotnote folder delete <path>` | `hotnote folder del <path>` |

## New Features

1. **Folder list command** - List contents of a folder (Issue P)
2. **Short aliases** for all commands (Issue O)

---

## Command Mapping

### Folder Commands (Unified under `folder`)

| Primary | Alias | Description |
|---------|-------|-------------|
| `folder make <path>` | `mk` | Create a folder |
| `folder remove <path>` | `rmv` | Delete a folder |
| `folder list [path]` | `ls` | List folder contents |
| `folder rename <old> <new>` | `rn` | Rename a folder |

### Note Commands (Aliases added)

| Primary | Alias | Description |
|---------|-------|-------------|
| `list [path]` | `ls` | List all notes |
| `delete <slug>` | `del` | Delete a note |
| `open <slug>` | `opn` | Open a note in editor |
| `render <slug>` | `rdr` | Render note as markdown |
| `rename <old> <new>` | `rn` | Rename a note |

### Workspace Commands (Aliases to add)

| Primary | Alias | Description |
|---------|-------|-------------|
| `workspace list` | `workspace ls` | List workspaces |
| `workspace delete <name>` | `workspace del` | Delete workspace |

---

## folder list Specification

### Usage
```bash
hotnote folder list [path]
hotnote folder ls [path]
```

### Arguments
| Argument | Description |
|----------|-------------|
| `path` | Optional folder path (defaults to workspace root) |

### Behavior
- Lists all files and folders in the specified path
- Supports nested paths (e.g., `folder list projects/2024`)
- Shows direct children only (not recursive)

### JSON Output Format
Follows `workspace list` pattern - flat array:

```json
[
  {"name": "notes.md", "path": "projects/notes.md", "type": "file"},
  {"name": "subfolder", "path": "projects/subfolder", "type": "folder"}
]
```

**Fields:**
- `name` - File or folder name
- `path` - Relative path from workspace root
- `type` - Either `"file"` or `"folder"`

### Human Output Format

```
projects/
  notes.md
  subfolder/
```

### Error Handling

| Error | Exit Code | Message |
|-------|-----------|---------|
| Workspace not initialized | 4 | `workspace not initialized` |
| Path does not exist | 2 | `folder not found: <path>` |
| Path outside workspace | 3 | `invalid folder path: must be inside workspace` |

---

## Alias Strategy

Use Cobra's native `Aliases` field:

```go
var listCmd = &cobra.Command{
    Use:     "list",
    Aliases: []string{"ls"},
    Short:   "List all notes",
    // ...
}
```

**Benefits:**
- Built-in help text automatically shows aliases
- No code duplication
- Works at all levels (top-level and subcommands)

---

## Exit Codes

Consistent across all commands:

| Code | Meaning | Example |
|------|---------|---------|
| 0 | Success | Command executed successfully |
| 1 | General error | Unexpected error, empty folder |
| 2 | Not found | Note, workspace, or folder not found |
| 3 | Invalid input | Missing arguments, path outside workspace |
| 4 | Config error | Workspace not initialized |

---

## Implementation Details

### New Files Required

1. `cmd/folder_list.go` - New folder list command
2. `cmd/folder_make.go` - Refactored from mkdir.go
3. `cmd/folder_remove.go` - Refactored from rmdir.go
4. Test files for above

### Files to Update

| File | Changes |
|------|---------|
| `cmd/folder.go` | Add list subcommand, add aliases to rename |
| `cmd/list.go` | Add `ls` alias |
| `cmd/delete.go` | Add `del` alias |
| `cmd/open.go` | Add `opn` alias |
| `cmd/render.go` | Add `rdr` alias |
| `cmd/rename.go` | Add `rn` alias |
| `cmd/workspace.go` | Add `ls`, `del` aliases |

### Files to Delete

- `cmd/mkdir.go` - Replaced by folder_make.go
- `cmd/rmdir.go` - Replaced by folder_remove.go

### Test Updates Required

- `cmd/mkdir_test.go` → Refactor to `cmd/folder_make_test.go`
- `cmd/rmdir_test.go` → Refactor to `cmd/folder_remove_test.go`
- Update all test calls to use new command names
- Add tests for alias functionality

---

## Documentation Updates

### Required Updates

1. **docs/features/commands.md**
   - Rewrite folder command section with new structure
   - Add alias information to all command documentation
   - Update examples to use new syntax

2. **docs/architecture/overview.md**
   - Update command reference table

3. **docs/index.md**
   - Update quick reference guide with new commands and aliases

4. **.ai/phase1/hotnote_cli_spec.md**
   - Update CLI specification with new command structure

### Documentation Example

```markdown
## hotnote folder

Folder management commands.

### Subcommands

| Command | Alias | Description |
|---------|-------|-------------|
| `folder make <path>` | `mk` | Create a folder |
| `folder remove <path>` | `rmv` | Delete a folder |
| `folder list [path]` | `ls` | List folder contents |
| `folder rename <old> <new>` | `rn` | Rename a folder |

### Examples

```bash
# Create a folder
hotnote folder make projects
hotnote folder mk projects

# List folder contents
hotnote folder list projects
hotnote folder ls projects

# Remove a folder
hotnote folder remove projects
hotnote folder rmv projects
```
```

---

## Implementation Order

1. **Issue P first** - Add folder list command
   - Create cmd/folder_list.go
   - Add tests
   - Document in commands.md

2. **Issue O** - CLI Command Restructure
   - Refactor mkdir → folder make
   - Refactor rmdir → folder remove
   - Add aliases to all commands
   - Update tests
   - Update all documentation

---

## Backlog Reference

See `.ai/backlog.md`:
- **Issue P:** Add folder list Command
- **Issue O:** CLI Command Restructure - Breaking Changes