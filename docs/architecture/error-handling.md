# Error Handling

Hotnote uses Go's error handling patterns with consistent exit codes and structured error messages.

## Exit Codes

All CLI commands return one of these exit codes:

| Code | Constant | Meaning |
|------|----------|---------|
| 0 | `ExitSuccess` | Command executed successfully |
| 1 | `ExitGeneral` | Unexpected error occurred |
| 2 | `ExitNotFound` | Note or workspace not found |
| 3 | `ExitInvalidInput` | Missing or invalid arguments |
| 4 | `ExitConfigError` | Configuration or workspace error |

### Implementation

Defined in `internal/errors/errors.go`:

```go
package exitorrors

const (
    ExitSuccess      = 0
    ExitGeneral      = 1
    ExitNotFound     = 2
    ExitInvalidInput = 3
    ExitConfigError  = 4
)
```

### Usage in Commands

```go
import exitorrors "hotnotego/internal/errors"

// In command Run function
if err != nil {
    fmt.Printf("error: %v\n", err)
    os.Exit(exitorrors.ExitGeneral)
}

if notFound {
    fmt.Println("note not found")
    os.Exit(exitorrors.ExitNotFound)
}
```

## Error Message Format

Error messages follow a consistent format:
- Lowercase
- No "Error:" prefix
- Include context for debugging

### Good Examples

```
note not found
workspace not initialized
create note: file exists
invalid slug: empty title
```

### Bad Examples

```
Error: Note not found
ERROR: The note could not be found
Failed to create note
```

## Error Wrapping

Use `fmt.Errorf` with `%w` to wrap errors while preserving the error chain:

```go
func (s *Store) Ensure(id string, content []byte) error {
    path, err := s.Path(id)
    if err != nil {
        return fmt.Errorf("ensure: %w", err)  // Wraps the error
    }
    if err := fsutil.AtomicWriteExclusive(path, content, 0644); err != nil {
        return fmt.Errorf("ensure: %w", err)
    }
    return nil
}
```

### Checking Wrapped Errors

Use `errors.Is` to check for specific error types:

```go
if errors.Is(err, os.ErrNotExist) {
    // Handle "not found"
}

if errors.Is(err, workspace.ErrWorkspaceNotInitialized) {
    // Handle workspace not initialized
}
```

## Sentinel Errors

Define reusable errors as package-level variables:

### Storage Errors

```go
var ErrWorkspaceNotInitialized = errors.New("workspace not initialized")
```

### Workspace Errors

```go
var (
    ErrWorkspaceNotInitialized = errors.New("workspace not initialized")
    ErrWorkspaceAlreadyExists  = errors.New("workspace already exists")
    ErrWorkspaceDoesNotExist   = errors.New("workspace does not exist")
)
```

## JSON Error Responses

When `--json` flag is used, errors are formatted as JSON:

```go
if jsonFlag {
    errorResponse := map[string]string{"error": "note not found"}
    jsonError, _ := json.Marshal(errorResponse)
    fmt.Println(string(jsonError))
} else {
    fmt.Println("note not found")
}
os.Exit(exitorrors.ExitNotFound)
```

Example output:
```json
{"error": "note not found"}
```

## Best Practices

### 1. Wrap Errors with Context

Add context at each layer:

```go
// CLI layer
store, err := storage.NewStore(wm)
if err != nil {
    return fmt.Errorf("create storage: %w", err)
}

// Storage layer
func (s *Store) Path(id string) (string, error) {
    _, path, err := s.wm.Current()
    if err != nil {
        return "", fmt.Errorf("storage: get current workspace: %w", err)
    }
    return filepath.Join(path, id+".md"), nil
}
```

### 2. Return Early on Errors

Avoid deep nesting with early returns:

```go
func processNote(id string) error {
    if err := validate(id); err != nil {
        return err  // Early return
    }
    
    // Main logic
    if err := doSomething(); err != nil {
        return err  // Early return
    }
    
    return nil
}
```

### 3. Use Descriptive Messages

Include enough context to debug:

```go
// Good
fmt.Printf("create note: %w", err)

// Better (in context)
fmt.Printf("create note '%s': %w", slug, err)
```

### 4. Handle Edge Cases

Check for common failure modes:

```go
// Check file exists
if _, err := os.Stat(path); os.IsNotExist(err) {
    return ErrNotFound
}

// Check permissions
if os.IsPermission(err) {
    return fmt.Errorf("permission denied: %w", err)
}
```

## Testing Errors

Use table-driven tests for error cases:

```go
func TestEnsure_AlreadyExists(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        wantErr error
    }{
        {"duplicate slug", "existing-note", os.ErrExist},
        {"empty id", "", ErrInvalidInput},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := store.Ensure(tt.id, content)
            assert.Error(t, err)
            assert.True(t, errors.Is(err, tt.wantErr))
        })
    }
}
```
