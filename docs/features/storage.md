# Storage

Hotnote uses a file-based storage system with markdown files and YAML frontmatter.

## Overview

```
workspace/
├── note1.md
├── note2.md
└── my-note.md
```

Each `.md` file contains optional YAML frontmatter followed by markdown content.

## File Format

```
---
id: <uuid>
title: <title>
created_at: <RFC3339 timestamp>
updated_at: <RFC3339 timestamp>
tags: [tag1, tag2]
---
# Note Title

Markdown content starts here...
```

### Frontmatter

The YAML block between `---` markers contains metadata:

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Unique identifier |
| `title` | string | Human-readable title |
| `created_at` | timestamp | RFC3339 format |
| `updated_at` | timestamp | RFC3339 format |
| `tags` | array | Optional categorization |

### Why Frontmatter?

- **Separation**: Metadata distinct from content
- **Portability**: Standard format (YAML)
- **Editor support**: Most markdown editors recognize frontmatter
- **Flexibility**: Easy to extend with new fields

## Storage Architecture

### Interface

```go
type WorkspaceManager interface {
    Current() (name string, path string, err error)
}
```

### Store Implementation

```go
type Store struct {
    wm WorkspaceManager  // Current workspace
}

func NewStore(wm WorkspaceManager) *Store {
    return &Store{wm: wm}
}
```

The Store depends on WorkspaceManager to determine where to store files.

## Core Operations

### Path Resolution

Convert a note slug to full file path:

```go
func (s *Store) Path(id string) (string, error) {
    _, rootPath, err := s.wm.Current()
    if err != nil {
        return "", err
    }
    return filepath.Join(rootPath, id+".md"), nil
}
```

### Ensure Note Exists

Create a new note file atomically (fails if exists):

```go
func (s *Store) Ensure(id string, content []byte) error {
    path, err := s.Path(id)
    if err != nil {
        return err
    }
    if err := fsutil.AtomicWriteExclusive(path, content, 0644); err != nil {
        return fmt.Errorf("ensure: %w", err)
    }
    return nil
}
```

Uses `fsutil.AtomicWriteExclusive` for atomic file creation with `os.O_EXCL`.

### Write Frontmatter

Create note with initial metadata:

```go
func createNote(title string, tags []string) error {
    slug := slugify(title)
    id := uuid.New().String()

    frontmatter := fmt.Sprintf(`---
id: %s
title: %s
created_at: %s
updated_at: %s
tags: %v
---

# %s
`,
        id,
        title,
        time.Now().Format(time.RFC3339),
        time.Now().Format(time.RFC3339),
        tags,
        title,
    )

    return store.Ensure(slug, []byte(frontmatter))
}
```

## Note Lookup

### By Slug

The slug is derived from the filename:

```
~/.local/share/hotnote/workspaces/default/my-note.md
                                         ^^^^^^^
                                         slug: my-note
```

### By UUID

UUID is stored in frontmatter `id` field.

## File Operations

### Read Note

```go
func readNote(path string) (string, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("read: %w", err)
    }
    return string(content), nil
}
```

### Update Note

After editing, frontmatter `updated_at` should be updated:

```go
func updateTimestamp(path string) error {
    content, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    // Parse frontmatter, update updated_at, rewrite file
    updated := updateFrontmatterTime(content, time.Now())

    return os.WriteFile(path, []byte(updated), 0644)
}
```

### Delete Note

```go
func deleteNote(slug string) error {
    path, err := store.Path(slug)
    if err != nil {
        return err
    }
    return os.Remove(path)
}
```

## Error Handling

| Error | Cause | Handling |
|-------|-------|----------|
| `os.ErrNotExist` | Note doesn't exist | Create or show error |
| `*os.PathError` | File system error | Log and retry |
| `*os.LinkError` | Symbolic link issue | Resolve path |

## Testing

The storage layer is tested with table-driven tests:

```go
func TestEnsure_NewFile(t *testing.T) {
    store := NewStore(mockWM)

    err := store.Ensure("new-note", []byte("# Test\n"))
    require.NoError(t, err)

    // Verify file was created
    content, _ := os.ReadFile(filepath.Join(workspaceDir, "new-note.md"))
    assert.Equal(t, "# Test\n", string(content))
}

func TestEnsure_AlreadyExists(t *testing.T) {
    store := NewStore(mockWM)

    err := store.Ensure("existing-note", []byte("# First\n"))
    require.NoError(t, err)

    // Second create should fail
    err = store.Ensure("existing-note", []byte("# Second\n"))
    assert.Error(t, err)
}
```

## Best Practices

1. **Sync after writes**: Always call `file.Sync()` for durability
2. **Atomic writes**: Write to temp file, then rename for atomic updates
3. **Permissions**: Use `0644` for files, `0755` for directories
4. **Slug collision**: Ensure slugs are unique within a workspace
