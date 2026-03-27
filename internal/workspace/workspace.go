package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"hotnotego/internal/fsutil"
)

var (
	// ErrWorkspaceNotInitialized is returned when trying to use workspace before initialization
	ErrWorkspaceNotInitialized = errors.New("workspace not initialized")
	// ErrWorkspaceAlreadyExists is returned when trying to create a workspace that already exists
	ErrWorkspaceAlreadyExists = errors.New("workspace already exists")
	// ErrWorkspaceDoesNotExist is returned when trying to use a workspace that doesn't exist
	ErrWorkspaceDoesNotExist = errors.New("workspace does not exist")
	// ErrCannotDeleteCurrent is returned when trying to delete the current workspace
	ErrCannotDeleteCurrent = errors.New("cannot delete current workspace")
	// ErrCannotDeleteDefaultStructure is returned when trying to delete the default workspace structure
	ErrCannotDeleteDefaultStructure = errors.New("cannot delete default workspace structure")
)

// Config represents the hotnote configuration
type Config struct {
	CurrentWorkspace string            `yaml:"current_workspace"`
	Workspaces       map[string]string `yaml:"workspaces"`
}

// Manager handles workspace operations
type Manager struct {
	configPath string
	config     *Config
}

// NewManager creates a new workspace manager
func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("workspace: get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "hotnote")
	configPath := filepath.Join(configDir, "config.yaml")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("workspace: create config directory: %w", err)
	}

	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}

	// Load existing config if it exists
	if err := m.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("workspace: load config: %w", err)
	}

	return m, nil
}

// load reads the configuration from file
func (m *Manager) load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, that's okay. We'll leave the config as the zero value.
			// Ensure Workspaces map is initialized
			if m.config.Workspaces == nil {
				m.config.Workspaces = make(map[string]string)
			}
			return nil
		}
		return fmt.Errorf("workspace: read config file: %w", err)
	}
	if err := yaml.Unmarshal(data, m.config); err != nil {
		return fmt.Errorf("workspace: unmarshal config: %w", err)
	}
	// Ensure Workspaces map is initialized after unmarshaling
	if m.config.Workspaces == nil {
		m.config.Workspaces = make(map[string]string)
	}
	return nil
}

// save writes the configuration to file atomically
func (m *Manager) save() error {
	var workspacesYAML string
	for name, path := range m.config.Workspaces {
		if len(workspacesYAML) > 0 {
			workspacesYAML += "\n"
		}
		workspacesYAML += fmt.Sprintf("  %s: %s", name, path)
	}
	configYAML := fmt.Sprintf("current_workspace: %s\nworkspaces:\n%s\n", m.config.CurrentWorkspace, workspacesYAML)
	if err := fsutil.AtomicWrite(m.configPath, []byte(configYAML), 0644); err != nil {
		return fmt.Errorf("workspace: write config: %w", err)
	}
	return nil
}

// Init initializes the default workspace
func (m *Manager) Init() error {
	if err := m.load(); err != nil {
		return fmt.Errorf("workspace: load config: %w", err)
	}

	// Check if already initialized
	if m.config.CurrentWorkspace != "" {
		return ErrWorkspaceAlreadyExists
	}

	// Set up default workspace
	m.config.CurrentWorkspace = "default"
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("workspace: get home directory: %w", err)
	}
	defaultPath := filepath.Join(home, ".local", "share", "hotnote", "workspaces", "default")
	m.config.Workspaces["default"] = defaultPath

	// Create the directory
	if err := os.MkdirAll(defaultPath, 0755); err != nil {
		return fmt.Errorf("workspace: create workspace directory: %w", err)
	}

	// Save config
	if err := m.save(); err != nil {
		return fmt.Errorf("workspace: save config: %w", err)
	}

	return nil
}

// List returns the list of workspaces
func (m *Manager) List() (map[string]string, string, error) {
	if err := m.load(); err != nil {
		return nil, "", fmt.Errorf("workspace: load config: %w", err)
	}
	return m.config.Workspaces, m.config.CurrentWorkspace, nil
}

// Use sets the current workspace
func (m *Manager) Use(name string) error {
	if err := m.load(); err != nil {
		return fmt.Errorf("workspace: load config: %w", err)
	}

	if _, exists := m.config.Workspaces[name]; !exists {
		return ErrWorkspaceDoesNotExist
	}
	m.config.CurrentWorkspace = name
	if err := m.save(); err != nil {
		return fmt.Errorf("workspace: save config: %w", err)
	}
	return nil
}

// New creates a new workspace
func (m *Manager) New(name string, path string) error {
	if err := m.load(); err != nil {
		return fmt.Errorf("workspace: load config: %w", err)
	}

	if _, exists := m.config.Workspaces[name]; exists {
		return ErrWorkspaceAlreadyExists
	}

	// Determine the path for the new workspace
	var workspacePath string
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("workspace: get home directory: %w", err)
		}
		workspacePath = filepath.Join(home, ".local", "share", "hotnote", "workspaces", name)
	} else {
		workspacePath = path
	}

	// Create the workspace directory
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return fmt.Errorf("workspace: create workspace directory: %w", err)
	}

	// Add to workspaces
	m.config.Workspaces[name] = workspacePath

	// Save config
	if err := m.save(); err != nil {
		return fmt.Errorf("workspace: save config: %w", err)
	}

	return nil
}

// Current returns the current workspace name and path
func (m *Manager) Current() (string, string, error) {
	if err := m.load(); err != nil {
		return "", "", fmt.Errorf("workspace: load config: %w", err)
	}
	if m.config.CurrentWorkspace == "" {
		return "", "", ErrWorkspaceNotInitialized
	}
	path := m.config.Workspaces[m.config.CurrentWorkspace]
	return m.config.CurrentWorkspace, path, nil
}

// Exists checks if a workspace exists
func (m *Manager) Exists(name string) (bool, error) {
	if err := m.load(); err != nil {
		return false, fmt.Errorf("workspace: load config: %w", err)
	}
	_, exists := m.config.Workspaces[name]
	return exists, nil
}

// Delete removes a workspace entirely (non-default workspaces only)
func (m *Manager) Delete(name string) error {
	if err := m.load(); err != nil {
		return fmt.Errorf("workspace: load config: %w", err)
	}

	workspacePath, exists := m.config.Workspaces[name]
	if !exists {
		return ErrWorkspaceDoesNotExist
	}

	if name == "default" {
		return ErrCannotDeleteDefaultStructure
	}

	if name == m.config.CurrentWorkspace {
		return ErrCannotDeleteCurrent
	}

	if err := os.RemoveAll(workspacePath); err != nil {
		return fmt.Errorf("workspace: delete directory: %w", err)
	}

	delete(m.config.Workspaces, name)

	if err := m.save(); err != nil {
		return fmt.Errorf("workspace: save config: %w", err)
	}

	return nil
}

// ClearDefaultWorkspace removes all contents from the default workspace
func (m *Manager) ClearDefaultWorkspace() error {
	if err := m.load(); err != nil {
		return fmt.Errorf("workspace: load config: %w", err)
	}

	workspacePath, exists := m.config.Workspaces["default"]
	if !exists {
		return ErrWorkspaceDoesNotExist
	}

	entries, err := os.ReadDir(workspacePath)
	if err != nil {
		return fmt.Errorf("workspace: read directory: %w", err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(workspacePath, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return fmt.Errorf("workspace: delete %s: %w", entry.Name(), err)
		}
	}

	return nil
}
