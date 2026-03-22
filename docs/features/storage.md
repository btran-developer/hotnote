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

Create a new note file (fails if exists):

```go
func (s *Store) Ensure(id string) error {
    path, err := s.Path(id)
    if err != nil {
        return err
    }

    // O_CREATE: Create if not exists
    // O_EXCL: Fail if file exists
    // O_WRONLY: Write-only access
    f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("ensure: %w", err)
    }
    defer f.Close()

    return nil
}
```

### Write Frontmatter

Create note with initial metadata:

```go
func createNote(title string, tags []string) error {
    slug := slugify(title)
    id := uuid.New().String()

    path := filepath.Join(workspacePath, slug+".md")
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    // Write YAML frontmatter
    fmt.Fprintf(f, `---
id: %s
title: %s
created_at: %s
updated_at: %s
tags: %v
---
`,
        id,
        title,
        time.Now().Format(time.RFC3339),
        time.Now().Format(time.RFC3339),
        tags,
    )

    return f.Sync()
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
func TestStore_Ensure(t *testing.T) {
    tests := []struct {
        name    string
        slug    string
        wantErr bool
    }{
        {"new note", "test-note", false},
        {"duplicate", "test-note", true},  // Already exists
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := store.Ensure(tt.slug)
            if (err != nil) != tt.wantErr {
                t.Errorf("Ensure() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Best Practices

1. **Sync after writes**: Always call `file.Sync()` for durability
2. **Atomic writes**: Write to temp file, then rename for atomic updates
3. **Permissions**: Use `0644` for files, `0755` for directories
4. **Slug collision**: Ensure slugs are unique within a workspace
