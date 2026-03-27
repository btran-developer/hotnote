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

	_, _, err = m.Current()
	assert.ErrorIs(t, err, ErrWorkspaceNotInitialized)
}

func TestDelete_NonExistent(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "default",
			Workspaces: map[string]string{
				"default": filepath.Join(configDir, "workspaces", "default"),
			},
		},
	}
	m.save()

	err = m.Delete("nonexistent")
	assert.ErrorIs(t, err, ErrWorkspaceDoesNotExist)
}

func TestDelete_DefaultWorkspace(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	wsPath := filepath.Join(configDir, "workspaces", "default")
	err = os.MkdirAll(wsPath, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "default",
			Workspaces: map[string]string{
				"default": wsPath,
			},
		},
	}
	m.save()

	err = m.Delete("default")
	assert.ErrorIs(t, err, ErrCannotDeleteDefaultStructure)
}

func TestDelete_CurrentWorkspace(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	wsPath := filepath.Join(configDir, "workspaces", "work")
	err = os.MkdirAll(wsPath, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "work",
			Workspaces: map[string]string{
				"default": filepath.Join(configDir, "workspaces", "default"),
				"work":    wsPath,
			},
		},
	}
	m.save()

	err = m.Delete("work")
	assert.ErrorIs(t, err, ErrCannotDeleteCurrent)
}

func TestDelete_Success(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	defaultPath := filepath.Join(configDir, "workspaces", "default")
	err = os.MkdirAll(defaultPath, 0755)
	require.NoError(t, err)

	workPath := filepath.Join(configDir, "workspaces", "work")
	err = os.MkdirAll(workPath, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			CurrentWorkspace: "default",
			Workspaces: map[string]string{
				"default": defaultPath,
				"work":    workPath,
			},
		},
	}
	m.save()

	err = m.Delete("work")
	require.NoError(t, err)

	assert.NotContains(t, m.config.Workspaces, "work")

	_, err = os.Stat(workPath)
	assert.True(t, os.IsNotExist(err))
}

func TestClearDefaultWorkspace_Empty(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	defaultPath := filepath.Join(configDir, "workspaces", "default")
	err = os.MkdirAll(defaultPath, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: map[string]string{
				"default": defaultPath,
			},
		},
	}
	m.save()

	err = m.ClearDefaultWorkspace()
	require.NoError(t, err)

	entries, err := os.ReadDir(defaultPath)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestClearDefaultWorkspace_WithContents(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(configDir)

	defaultPath := filepath.Join(configDir, "workspaces", "default")
	err = os.MkdirAll(defaultPath, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(defaultPath, "note1.md"), []byte("# Note 1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(defaultPath, "note2.md"), []byte("# Note 2"), 0644)
	require.NoError(t, err)

	subDir := filepath.Join(defaultPath, "subfolder")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(subDir, "note3.md"), []byte("# Note 3"), 0644)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, "config.yaml")
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Workspaces: map[string]string{
				"default": defaultPath,
			},
		},
	}
	m.save()

	err = m.ClearDefaultWorkspace()
	require.NoError(t, err)

	entries, err := os.ReadDir(defaultPath)
	require.NoError(t, err)
	assert.Empty(t, entries)

	_, err = os.Stat(defaultPath)
	assert.NoError(t, err)
}
