# Workspace Management

Workspaces allow organizing notes into separate directories, each with its own set of notes.

## Concept

A **workspace** is simply a directory containing markdown notes. You can have multiple workspaces for different purposes:

```
~/.local/share/hotnote/workspaces/
├── default/          # Default workspace
│   ├── note1.md
│   └── note2.md
├── work/             # Work-related notes
│   ├── meeting-notes.md
│   └── project-ideas.md
└── personal/         # Personal notes
    ├── journal.md
    └── recipes.md
```

## Why Workspaces?

- **Organization**: Separate notes by project, context, or topic
- **Isolation**: Work notes don't mix with personal notes
- **Portability**: Entire workspace is just a directory (easy to backup, sync)

## Configuration

### Config File Location

```
~/.config/hotnote/config.yaml
```

### Config Structure

```yaml
current: default
workspaces:
  - name: default
    path: /home/user/.local/share/hotnote/workspaces/default
  - name: work
    path: /home/user/.local/share/hotnote/workspaces/work
```

### Storage Locations

| Item | Location |
|------|----------|
| Config | `~/.config/hotnote/config.yaml` |
| Workspace data | `~/.local/share/hotnote/workspaces/` |

## Architecture

### Manager Interface

```go
type WorkspaceManager interface {
    Current() (name string, path string, err error)
}
```

Any type implementing `Current()` can be used as a workspace manager.

### Manager Implementation

```go
type Manager struct {
    configPath string  // ~/.config/hotnote/config.yaml
    dataPath   string  // ~/.local/share/hotnote/workspaces/
    config     Config
}

type Config struct {
    Current     string               `yaml:"current"`
    Workspaces  []WorkspaceConfig   `yaml:"workspaces"`
}
```

### Error Handling

```go
var (
    ErrWorkspaceNotInitialized = errors.New("workspace not initialized")
    ErrWorkspaceAlreadyExists  = errors.New("workspace already exists")
    ErrWorkspaceDoesNotExist   = errors.New("workspace does not exist")
)
```

## Operations

### Initialize

Creates the default workspace if it doesn't exist:

```go
func (m *Manager) Init() error {
    if m.config.Current != "" {
        return nil  // Already initialized
    }

    // Create default workspace directory
    defaultPath := filepath.Join(m.dataPath, "default")
    if err := os.MkdirAll(defaultPath, 0755); err != nil {
        return fmt.Errorf("init: %w", err)
    }

    // Save config
    m.config.Current = "default"
    m.config.Workspaces = []WorkspaceConfig{
        {Name: "default", Path: defaultPath},
    }
    return m.save()
}
```

### List Workspaces

Returns all workspaces, marking the current one:

```go
func (m *Manager) List() ([]WorkspaceInfo, error) {
    var workspaces []WorkspaceInfo
    for _, ws := range m.config.Workspaces {
        workspaces = append(workspaces, WorkspaceInfo{
            Name:     ws.Name,
            Path:     ws.Path,
            IsCurrent: ws.Name == m.config.Current,
        })
    }
    return workspaces, nil
}
```

### Switch Workspace

Changes the current workspace:

```go
func (m *Manager) Use(name string) error {
    for _, ws := range m.config.Workspaces {
        if ws.Name == name {
            m.config.Current = name
            return m.save()
        }
    }
    return fmt.Errorf("workspace use: %w", ErrWorkspaceDoesNotExist)
}
```

### Create Workspace

Adds a new workspace:

```go
func (m *Manager) New(name string, customPath string) error {
    // Check if exists
    for _, ws := range m.config.Workspaces {
        if ws.Name == name {
            return fmt.Errorf("workspace new: %w", ErrWorkspaceAlreadyExists)
        }
    }

    // Determine path
    path := customPath
    if path == "" {
        path = filepath.Join(m.dataPath, name)
    }

    // Create directory
    if err := os.MkdirAll(path, 0755); err != nil {
        return fmt.Errorf("workspace new: %w", err)
    }

    // Update config
    m.config.Workspaces = append(m.config.Workspaces, WorkspaceConfig{
        Name: name,
        Path: path,
    })
    return m.save()
}
```

## CLI Integration

The `Store` type depends on `WorkspaceManager` interface:

```go
type Store struct {
    wm WorkspaceManager  // Can be any implementation
}

func NewStore(wm WorkspaceManager) *Store {
    return &Store{wm: wm}
}

func (s *Store) Path(slug string) (string, error) {
    _, path, err := s.wm.Current()
    if err != nil {
        return "", err
    }
    return filepath.Join(path, slug+".md"), nil
}
```

This dependency injection pattern allows:
- Easy testing (mock WorkspaceManager)
- Flexibility in implementation
- Decoupled components

## Best Practices

1. **Naming**: Use lowercase, hyphenated names (`work-notes`, `research`)
2. **Organization**: One workspace per major context (work, personal, projects)
3. **Backup**: Workspace directories can be synced with cloud storage
4. **Migration**: Copy workspace directory to move notes between machines
