# Go Concepts

This guide explains Go patterns used in Hotnote for developers unfamiliar with the language.

## Packages and Imports

Go organizes code into packages. Each `.go` file starts with a `package` declaration.

```go
package main  // Executable package

import (
    "fmt"        // Standard library
    "os"         // Standard library
    "path/filepath"

    "github.com/spf13/cobra"  // Third-party
    "github.com/google/uuid"  // Third-party
)
```

### Package Visibility

In Go, names are exported (visible outside the package) if they start with an **uppercase letter**.

```go
package storage

type Store struct { ... }  // Exported - can be accessed from other packages

type store struct { ... }  // Unexported - only accessible within this package
```

## Interfaces

Interfaces define behavior without specifying implementation. They're implemented implicitly.

```go
// WorkspaceManager defines methods for workspace operations
type WorkspaceManager interface {
    Current() (string, string, error)  // Returns name, path, error
}
```

Any type that implements these methods satisfies the interface - no explicit declaration needed.

```go
// Manager implements WorkspaceManager
type Manager struct { ... }

func (m *Manager) Current() (string, string, error) {
    // Implementation here
    return m.current, m.rootPath, nil
}
```

## Error Handling

Go uses explicit error returns rather than exceptions.

### Sentinel Errors

Define reusable errors as package-level variables:

```go
var (
    ErrWorkspaceNotInitialized = errors.New("workspace not initialized")
    ErrWorkspaceDoesNotExist   = errors.New("workspace does not exist")
)
```

### Wrapping Errors

Add context when propagating errors:

```go
func (m *Manager) Use(name string) error {
    if !exists(name) {
        return fmt.Errorf("workspace use: %w", ErrWorkspaceDoesNotExist)
    }
    return nil
}
```

### Checking Errors

```go
if errors.Is(err, storage.ErrWorkspaceNotInitialized) {
    // Handle specific error
}
```

## Structs and Methods

Structs are Go's way of defining data structures. Methods are functions attached to structs.

```go
type Note struct {
    ID        string
    Title     string
    Path      string
    Tags      []string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Value receiver - receives a copy
func (n Note) FullPath() string {
    return n.Path + "/" + n.ID + ".md"
}

// Pointer receiver - modifies the struct
func (n *Note) UpdateTitle(newTitle string) {
    n.Title = newTitle
}
```

## Pointers

Go passes by value by default. Use pointers (`*T`) to modify values or avoid copying large structs.

```go
// Without pointer - gets a copy
func badTake(n Note) {
    n.Title = "changed"  // Doesn't affect original
}

// With pointer - gets reference
func goodTake(n *Note) {
    n.Title = "changed"  // Affects original
}
```

## YAML Configuration

Hotnote uses `gopkg.in/yaml.v3` for configuration.

```go
type Config struct {
    Current string `yaml:"current"`  // YAML key -> struct field
    Workspaces []WorkspaceConfig `yaml:"workspaces"`
}

func SaveConfig(cfg Config) error {
    data, err := yaml.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("marshal config: %w", err)
    }
    return os.WriteFile(configPath, data, 0644)
}
```

## The `init()` Function

Each package can have an `init()` function that runs before `main()`:

```go
package cmd

var RootCmd = &cobra.Command{...}

func init() {
    RootCmd.AddCommand(NewCmd)  // Register subcommands
}
```

This pattern is used in Cobra to register commands without explicit calls in `main()`.

## Defer

`defer` schedules a function to run when the surrounding function returns:

```go
func readNote(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // Ensures file is closed even if function returns early

    // Process file...
    return nil
}
```

## Project Structure Pattern

Hotnote uses the standard Go project layout:

```
cmd/                    # Application entry points
  ├── hotnote/
  │   └── main.go      # Calls RootCmd.Execute()
  ├── root.go          # Root command definition
  ├── new.go           # "new" subcommand
  └── ...

internal/               # Private packages (not importable)
  ├── core/            # Domain models
  ├── storage/         # Storage implementation
  └── workspace/       # Workspace management
```
