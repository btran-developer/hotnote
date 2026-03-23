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

---

## Benchmarking

HotNote includes benchmarks to ensure operations meet the performance requirement of **<100ms** (PRD §163).

### What Are Benchmarks?

Benchmarks measure the performance of code operations:
- **Unlike tests**: Tests verify correctness; benchmarks measure speed and memory
- **Go benchmarks**: Use `testing.B` instead of `testing.T`
- **Automatic iteration**: Go runs benchmarks multiple times to get statistically significant results

Benchmarked operations:
- `new` - Creating notes
- `list` - Listing notes (10, 100, 500 notes)
- `open` - Opening notes for editing
- `render` - Converting markdown to HTML
- `workspace` - Workspace manager operations

### How Benchmarks Work

A benchmark function looks like this:

```go
func BenchmarkListNotes(b *testing.B) {
    // Setup (not timed)
    workspacePath, cleanup := setupTestWorkspace(b, 500)
    defer cleanup()
    
    b.ResetTimer()  // Start timing here
    
    for i := 0; i < b.N; i++ {
        // This loop runs multiple times
        listNotes(workspacePath)
    }
}
```

**Key concepts:**

| Function | Purpose |
|----------|---------|
| `b.N` | Number of iterations Go will run (automatically determined) |
| `b.ResetTimer()` | Excludes setup time from measurements |
| `b.StopTimer()` / `b.StartTimer()` | Pause timing for non-benchmark code |

### Running Benchmarks

#### Using the Benchmark Script (Recommended)

```bash
./scripts/benchmark.sh
```

This runs all benchmarks, generates a report in `.benchmarks/`, and highlights any operations exceeding the 100ms threshold.

**Profiling with the script:**

```bash
# Generate CPU profile
./scripts/benchmark.sh --cpuprofile
go tool pprof -http=:8080 cpu.prof

# Generate memory profile
./scripts/benchmark.sh --memprofile
go tool pprof -http=:8080 mem.prof
```

The profile files (`cpu.prof` or `mem.prof`) are generated in the project directory after running with profiling flags.

#### Manual Benchmark Runs

```bash
# Run all benchmarks with memory stats
go test -bench=. -benchmem ./cmd/

# Run specific benchmark
go test -bench=BenchmarkListNotes -benchmem ./cmd/

# Run benchmarks matching pattern
go test -bench=BenchmarkList -benchmem ./cmd/

# Run with CPU profiling
go test -bench=. -cpuprofile=.benchmarks/cpu.prof ./cmd/
go tool pprof .benchmarks/cpu.prof
```

**Common flags:**
- `-bench=.` - Run all benchmarks (regex match)
- `-benchmem` - Include memory allocation stats
- `-count=5` - Run each benchmark 5 times (for statistical significance)
- `-run=^$` - Skip tests, run only benchmarks

### Interpreting Results

Example benchmark output:

```
BenchmarkListNotes/notes_500-8    816   1358188 ns/op   328404 B/op   2681 allocs/op
```

Breaking it down:

| Field | Value | Meaning |
|-------|-------|---------|
| Name | `BenchmarkListNotes/notes_500-8` | Benchmark name with subtest and GOMAXPROCS (8) |
| Iterations | `816` | How many times the loop ran (higher = more reliable) |
| Time | `1358188 ns/op` | 1.36 milliseconds per operation |
| Memory | `328404 B/op` | Bytes allocated per operation (~328 KB) |
| Allocs | `2681 allocs/op` | Heap allocations per operation |

**Status indicators:**
- ✅ Under 100ms threshold
- ⚠️ Exceeds 100ms threshold (needs optimization)

### Performance Targets

Per PRD §163, all operations must complete in <100ms on SSD hardware.

| Operation | Target | Typical | Status |
|-----------|--------|---------|--------|
| `new` | <100ms | ~5ms | ✅ ~20x faster |
| `list` (500 notes) | <100ms | ~1ms | ✅ ~100x faster |
| `open` | <100ms | ~0.002ms | ✅ ~50000x faster |
| `render` | <100ms | ~0.005ms | ✅ ~20000x faster |
| `workspace` init | <100ms | ~0.04ms | ✅ ~2500x faster |

### Comparing Benchmarks Over Time

#### Method 1: Simple File Comparison

Save benchmark results and compare with `diff`:

```bash
# Before changes
go test -bench=. -benchmem ./cmd/ > .benchmarks/baseline.txt

# ... make code changes ...

# After changes
go test -bench=. -benchmem ./cmd/ > .benchmarks/current.txt

# Compare
diff .benchmarks/baseline.txt .benchmarks/current.txt
```

#### Method 2: Using benchstat (Recommended)

`benchstat` provides statistical analysis of benchmark results:

```bash
# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Run benchmarks multiple times before changes
go test -bench=. -count=5 ./cmd/ > .benchmarks/old.txt

# ... make code changes ...

# Run benchmarks multiple times after changes
go test -bench=. -count=5 ./cmd/ > .benchmarks/new.txt

# Compare with statistical significance
benchstat .benchmarks/old.txt .benchmarks/new.txt
```

**Sample benchstat output:**

```
name                    old time/op    new time/op    delta
BenchmarkListNotes-8    1.36ms ± 3%    1.25ms ± 2%   -8.12%  (p=0.002 n=5+5)

name                    old alloc/op   new alloc/op   delta
BenchmarkListNotes-8     328kB ± 0%     310kB ± 0%   -5.49%  (p=0.008 n=5+5)
```

**Understanding the output:**
- `± 3%` - Variation between runs (lower is more consistent)
- `-8.12%` - Performance improvement (negative = faster)
- `(p=0.002)` - P-value < 0.05 means the change is statistically significant
- `n=5+5` - 5 runs before, 5 runs after

#### Method 3: Benchmark Report History

The benchmark script saves dated reports:

```bash
# List all reports
ls -la .benchmarks/

# Compare two specific reports
diff .benchmarks/report_20260322_120000.txt .benchmarks/report_20260322_180000.txt
```

### Writing New Benchmarks

To add a benchmark to `cmd/benchmark_test.go`:

```go
func BenchmarkMyOperation(b *testing.B) {
    // Setup test data
    workspacePath, configPath, cleanup := setupBenchmarkWorkspace(b, 100)
    defer cleanup()
    
    // Set up environment
    restoreConfig := setBenchmarkConfig(b, configPath)
    defer restoreConfig()
    
    b.ResetTimer()  // Start timing
    
    for i := 0; i < b.N; i++ {
        // Code to benchmark
        err := myOperation(workspacePath)
        if err != nil {
            b.Fatalf("operation failed: %v", err)
        }
    }
}
```

**Best practices:**
1. Use `setupBenchmarkWorkspace()` to create isolated test data
2. Use `setBenchmarkConfig()` to configure environment
3. Call `b.ResetTimer()` after setup
4. Check for errors inside the loop
5. Clean up with `defer cleanup()`

### Benchmark Reports

Benchmark reports are saved to `.benchmarks/report_YYYYMMDD_HHMMSS.txt`:

```
.benchmarks/
├── report_20260322_120000.txt
├── report_20260322_180000.txt
└── cpu.prof  (if profiling enabled)
```

The `.benchmarks/` directory is gitignored to avoid committing large/temporary files.

### Troubleshooting

**Benchmarks fail with "workspace not initialized"**
- Ensure `setBenchmarkConfig()` is called before operations
- Check that config file is created at correct path

**Inconsistent results**
- Close other applications to reduce system load
- Run with `-count=5` or higher for statistical significance
- Check for background processes (indexing, backups)

**High memory usage**
- Large note counts (500+) will use more memory
- This is expected; check `allocs/op` for optimization opportunities

**"No benchmarks to run"**
- Ensure benchmark files end with `_test.go`
- Benchmark functions must start with `Benchmark`
- Use `-run=^$` to skip tests that might interfere
