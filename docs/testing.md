# Testing

Hotnote uses Go's standard testing package with testify for assertions.

## Running Tests

### Run all tests

```bash
go test -v ./...
```

### Run tests for a specific package

```bash
go test -v ./internal/storage/...
```

### Run tests with coverage

```bash
go test -cover ./...
```

### Run tests matching a pattern

```bash
go test -v -run TestEnsure ./internal/storage/...
```

## Test Structure

### Unit Tests

Located alongside source files with `_test.go` suffix:

```
internal/storage/
├── store.go
└── store_test.go
```

### Integration Tests

CLI commands are tested by building the binary and running commands:

```go
func TestWorkspaceInit(t *testing.T) {
    // Build the binary
    tmpDir := t.TempDir()
    bin := filepath.Join(tmpDir, "hotnote")
    cmd := exec.Command("go", "build", "-o", bin, "./cmd/hotnote")
    require.NoError(t, cmd.Run())

    // Run commands
    out, err := exec.Command(bin, "workspace", "init").CombinedOutput()
    require.NoError(t, err)
    assert.Contains(t, string(out), "Initialized workspace: default")
}
```

## Test Patterns

### Table-Driven Tests

For testing multiple input/output combinations:

```go
func TestSlugify(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"simple", "Hello World", "hello-world"},
        {"special chars", "Test@123!", "test123"},
        {"spaces", "a  b   c", "a-b-c"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := slugify(tt.input)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Mocking Dependencies

Use interfaces to inject mock implementations:

```go
type mockWorkspaceManager struct {
    currentName  string
    currentPath  string
    currentError error
}

func (m *mockWorkspaceManager) Current() (string, string, error) {
    return m.currentName, m.currentPath, m.currentError
}

// Use in tests
wm := &mockWorkspaceManager{
    currentName:  "test",
    currentPath:  tmpDir,
    currentError: nil,
}
store := NewStore(wm)
```

### Temporary Files

Use `t.TempDir()` for automatic cleanup:

```go
func TestEnsure(t *testing.T) {
    tmpDir := t.TempDir()  // Automatically cleaned up
    
    wm := &mockWorkspaceManager{
        currentName: "test",
        currentPath: tmpDir,
    }
    store := NewStore(wm)
    
    err := store.Ensure("test-note", []byte("# Test\n"))
    require.NoError(t, err)
    
    // File automatically cleaned up when test ends
}
```

### Cleanup with defer

For resources that don't have built-in cleanup:

```go
func TestIntegration(t *testing.T) {
    // Create temp directory
    tmpDir, err := os.MkdirTemp("", "hotnote-test-*")
    require.NoError(t, err)
    
    defer os.RemoveAll(tmpDir)  // Clean up on test end
    
    // ... rest of test
}
```

## Verification Approach

After implementing features, verify:

1. **Run all tests**: `go test ./...`
2. **Test edge cases**: Empty inputs, special characters, concurrent access
3. **Verify error paths**: Ensure errors are returned correctly
4. **Clean artifacts**: Remove config files and workspace directories after tests

### Manual Verification

Clean up before manual testing:

```bash
# Remove config
rm -f ~/.config/hotnote/config.yaml

# Remove workspaces
rm -rf ~/.local/share/hotnote/workspaces/

# Remove built binary
rm -f ./hotnote
```

Build and test:

```bash
go build -o hotnote ./cmd/hotnote
./hotnote workspace init
./hotnote new "Test Note"
./hotnote list
```

## Code Coverage

Check coverage for untested code:

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

## Continuous Integration

Tests should be deterministic and not depend on:
- External services
- Specific timing
- Global state
- User-specific directories

Each test should set up its own fixtures and clean up after itself.
