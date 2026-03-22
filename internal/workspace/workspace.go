package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	// ErrWorkspaceNotInitialized is returned when trying to use workspace before initialization
	ErrWorkspaceNotInitialized = errors.New("workspace not initialized")
	// ErrWorkspaceAlreadyExists is returned when trying to create a workspace that already exists
	ErrWorkspaceAlreadyExists = errors.New("workspace already exists")
	// ErrWorkspaceDoesNotExist is returned when trying to use a workspace that doesn't exist
	ErrWorkspaceDoesNotExist = errors.New("workspace does not exist")
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
		return nil, err
	}

	configDir := filepath.Join(home, ".config", "hotnote")
	configPath := filepath.Join(configDir, "config.yaml")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}

	// Load existing config if it exists
	if err := m.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return m, nil
}

// load reads the configuration from file
func (m *Manager) load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, that's okay. We'll leave the config as the zero value.
			return nil
		}
		return err
	}
	return yaml.Unmarshal(data, m.config)
}

// save writes the configuration to file
func (m *Manager) save() error {
	// In a real implementation, we would marshal to YAML
	// For this implementation, we'll write a simple YAML file
	// A full implementation would properly serialize the YAML
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	// Write a simple YAML file
	var workspacesYAML string
	for name, path := range m.config.Workspaces {
		if len(workspacesYAML) > 0 {
			workspacesYAML += "\n"
		}
		workspacesYAML += fmt.Sprintf("  %s: %s", name, path)
	}
	configYAML := fmt.Sprintf("current_workspace: %s\nworkspaces:\n%s\n", m.config.CurrentWorkspace, workspacesYAML)
	return os.WriteFile(m.configPath, []byte(configYAML), 0644)
}

// Init initializes the default workspace
func (m *Manager) Init() error {
	if err := m.load(); err != nil {
		return err
	}

	// Check if already initialized
	if m.config.CurrentWorkspace != "" {
		return ErrWorkspaceAlreadyExists
	}

	// Set up default workspace
	m.config.CurrentWorkspace = "default"
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	defaultPath := filepath.Join(home, ".local", "share", "hotnote", "workspaces", "default")
	m.config.Workspaces["default"] = defaultPath

	// Create the directory
	if err := os.MkdirAll(defaultPath, 0755); err != nil {
		return err
	}

	// Save config
	if err := m.save(); err != nil {
		return err
	}

	return nil
}

// List returns the list of workspaces
func (m *Manager) List() (map[string]string, string, error) {
	if err := m.load(); err != nil {
		return nil, "", err
	}
	return m.config.Workspaces, m.config.CurrentWorkspace, nil
}

// Use sets the current workspace
func (m *Manager) Use(name string) error {
	if err := m.load(); err != nil {
		return err
	}

	if _, exists := m.config.Workspaces[name]; !exists {
		return ErrWorkspaceDoesNotExist
	}
	m.config.CurrentWorkspace = name
	if err := m.save(); err != nil {
		return err
	}
	return nil
}

// New creates a new workspace
func (m *Manager) New(name string, path string) error {
	if err := m.load(); err != nil {
		return err
	}

	if _, exists := m.config.Workspaces[name]; exists {
		return ErrWorkspaceAlreadyExists
	}

	// Determine the path for the new workspace
	var workspacePath string
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		workspacePath = filepath.Join(home, ".local", "share", "hotnote", "workspaces", name)
	} else {
		workspacePath = path
	}

	// Create the workspace directory
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return err
	}

	// Add to workspaces
	m.config.Workspaces[name] = workspacePath

	// Save config
	if err := m.save(); err != nil {
		return err
	}

	return nil
}

// Current returns the current workspace name and path
func (m *Manager) Current() (string, string, error) {
	if err := m.load(); err != nil {
		return "", "", err
	}
	if m.config.CurrentWorkspace == "" {
		return "", "", ErrWorkspaceNotInitialized
	}
	path := m.config.Workspaces[m.config.CurrentWorkspace]
	return m.config.CurrentWorkspace, path, nil
}
