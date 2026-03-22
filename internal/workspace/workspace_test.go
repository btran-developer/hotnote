package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")

	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}

	assert.NotNil(t, m)
	assert.Equal(t, configPath, m.configPath)
}

func TestInit_Default(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}

	err = m.Init()
	require.NoError(t, err)

	// Verify config was saved
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Verify workspace directory was created
	assert.Equal(t, "default", m.config.CurrentWorkspace)
	assert.Contains(t, m.config.Workspaces, "default")

	// Cleanup: remove the workspace directory
	workspacePath := m.config.Workspaces["default"]
	defer os.RemoveAll(filepath.Dir(workspacePath))
}

func TestInit_AlreadyExists(t *testing.T) {
	// Create a temp config directory with existing config
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")

	// Pre-initialize config
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "default",
			Workspaces:       map[string]string{"default": "/some/path"},
		},
	}

	// Save the config
	m.save()
	defer os.RemoveAll(configDir)

	// Try to init again
	err = m.Init()
	assert.ErrorIs(t, err, ErrWorkspaceAlreadyExists)
}

func TestList_Empty(t *testing.T) {
	// Create a temp config directory with empty config
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}
	m.save()

	workspaces, current, err := m.List()
	require.NoError(t, err)
	assert.Empty(t, workspaces)
	assert.Empty(t, current)
}

func TestList_WithWorkspaces(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "work",
			Workspaces: map[string]string{
				"default": "/tmp/default",
				"work":    "/tmp/work",
			},
		},
	}
	m.save()

	workspaces, current, err := m.List()
	require.NoError(t, err)
	assert.Equal(t, "work", current)
	assert.Len(t, workspaces, 2)
	assert.Contains(t, workspaces, "default")
	assert.Contains(t, workspaces, "work")
}

func TestUse_Valid(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "default",
			Workspaces: map[string]string{
				"default": "/tmp/default",
				"work":    "/tmp/work",
			},
		},
	}
	m.save()

	err = m.Use("work")
	require.NoError(t, err)
	assert.Equal(t, "work", m.config.CurrentWorkspace)
}

func TestUse_NotFound(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "default",
			Workspaces: map[string]string{
				"default": "/tmp/default",
			},
		},
	}
	m.save()

	err = m.Use("nonexistent")
	assert.ErrorIs(t, err, ErrWorkspaceDoesNotExist)
}

func TestNew_Simple(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")

	// Create a temp home directory for the workspace path
	tempHome, err := os.MkdirTemp("", "hotnote-home-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempHome)

	expectedPath := filepath.Join(tempHome, ".local", "share", "hotnote", "workspaces", "testws")

	// Create manager and initialize Workspaces map
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}
	m.save()

	// Now load into a fresh manager and create workspace
	m2 := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}
	m2.save()

	err = m2.New("testws", expectedPath)
	require.NoError(t, err)

	assert.Contains(t, m2.config.Workspaces, "testws")
	assert.Equal(t, expectedPath, m2.config.Workspaces["testws"])

	// Verify the directory was created
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err)
	defer os.RemoveAll(expectedPath)
}

func TestNew_CustomPath(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}
	m.save()

	// Create a custom path
	customPath, err := os.MkdirTemp("", "hotnote-custom-*")
	require.NoError(t, err)
	defer os.RemoveAll(customPath)

	err = m.New("custom", customPath)
	require.NoError(t, err)

	assert.Contains(t, m.config.Workspaces, "custom")
	assert.Equal(t, customPath, m.config.Workspaces["custom"])
}

func TestNew_AlreadyExists(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: map[string]string{
				"existing": "/tmp/existing",
			},
		},
	}
	m.save()

	err = m.New("existing", "")
	assert.ErrorIs(t, err, ErrWorkspaceAlreadyExists)
}

func TestCurrent_Set(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "work",
			Workspaces: map[string]string{
				"work": "/tmp/work",
			},
		},
	}
	m.save()

	name, path, err := m.Current()
	require.NoError(t, err)
	assert.Equal(t, "work", name)
	assert.Equal(t, "/tmp/work", path)
}

func TestCurrent_NotInitialized(t *testing.T) {
	// Create a temp config directory
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: make(map[string]string),
		},
	}
	m.save()

	_, _, err = m.Current()
	assert.ErrorIs(t, err, ErrWorkspaceNotInitialized)
}
