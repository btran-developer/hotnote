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
current_workspace: default
workspaces:
  default: /home/user/.local/share/hotnote/workspaces/default
  work: /home/user/.local/share/hotnote/workspaces/work
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
    config     *Config
}

type Config struct {
    CurrentWorkspace string            `yaml:"current_workspace"`
    Workspaces       map[string]string `yaml:"workspaces"`
}
```

### Error Handling

```go
var (
    ErrWorkspaceNotInitialized    = errors.New("workspace not initialized")
    ErrWorkspaceAlreadyExists    = errors.New("workspace already exists")
    ErrWorkspaceDoesNotExist     = errors.New("workspace does not exist")
    ErrCannotDeleteCurrent       = errors.New("cannot delete current workspace")
    ErrCannotDeleteDefaultStructure = errors.New("cannot delete default workspace structure")
    ErrCannotRenameDefault       = errors.New("cannot rename default workspace")
    ErrEmptyWorkspaceName        = errors.New("workspace name cannot be empty")
)
```

## Operations

### Initialize

Creates the default workspace if it doesn't exist:

```go
func (m *Manager) Init() error {
    if m.config.CurrentWorkspace != "" {
        return ErrWorkspaceAlreadyExists
    }

    home, _ := os.UserHomeDir()
    defaultPath := filepath.Join(home, ".local", "share", "hotnote", "workspaces", "default")
    
    if err := os.MkdirAll(defaultPath, 0755); err != nil {
        return fmt.Errorf("create workspace directory: %w", err)
    }

    m.config.CurrentWorkspace = "default"
    m.config.Workspaces["default"] = defaultPath
    return m.save()
}
```

### List Workspaces

Returns all workspaces, marking the current one:

```go
func (m *Manager) List() (map[string]string, string, error) {
    return m.config.Workspaces, m.config.CurrentWorkspace, nil
}
```

### Switch Workspace

Changes the current workspace:

```go
func (m *Manager) Use(name string) error {
    if _, exists := m.config.Workspaces[name]; !exists {
        return ErrWorkspaceDoesNotExist
    }
    m.config.CurrentWorkspace = name
    return m.save()
}
```

### Create Workspace

Adds a new workspace:

```go
func (m *Manager) New(name string, customPath string) error {
    if _, exists := m.config.Workspaces[name]; exists {
        return ErrWorkspaceAlreadyExists
    }

    var workspacePath string
    if customPath == "" {
        home, _ := os.UserHomeDir()
        workspacePath = filepath.Join(home, ".local", "share", "hotnote", "workspaces", name)
    } else {
        workspacePath = customPath
    }

    if err := os.MkdirAll(workspacePath, 0755); err != nil {
        return fmt.Errorf("create workspace directory: %w", err)
    }

    m.config.Workspaces[name] = workspacePath
    return m.save()
}
```

### Rename Workspace

Renames a workspace by updating the config map. The directory path remains unchanged.

```go
func (m *Manager) Rename(oldName, newName string) error {
    if oldName == "" || newName == "" {
        return ErrEmptyWorkspaceName
    }

    workspacePath, exists := m.config.Workspaces[oldName]
    if !exists {
        return ErrWorkspaceDoesNotExist
    }

    if _, exists := m.config.Workspaces[newName]; exists {
        return ErrWorkspaceAlreadyExists
    }

    if oldName == "default" {
        return ErrCannotRenameDefault
    }

    m.config.Workspaces[newName] = workspacePath
    delete(m.config.Workspaces, oldName)

    if m.config.CurrentWorkspace == oldName {
        m.config.CurrentWorkspace = newName
    }

    return m.save()
}
```

**Safety rules:**
- Cannot rename the `default` workspace
- New name must not already exist
- If renaming the current workspace, `current_workspace` is updated

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

### Workspace Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `workspace init` | - | Initialize default workspace |
| `workspace list` | `workspace ls` | List all workspaces |
| `workspace use <name>` | - | Switch to workspace |
| `workspace create <name>` | `workspace new` | Create workspace |
| `workspace delete <name>` | `workspace del` | Delete workspace |
| `workspace rename <old> <new>` | `workspace rn` | Rename workspace |

**Rename command:**
```
hotnote workspace rename <old> <new>   # Rename workspace
hotnote workspace rn work personal     # Using alias
```

## Best Practices

1. **Naming**: Use lowercase, hyphenated names (`work-notes`, `research`)
2. **Organization**: One workspace per major context (work, personal, projects)
3. **Backup**: Workspace directories can be synced with cloud storage
4. **Migration**: Copy workspace directory to move notes between machines
