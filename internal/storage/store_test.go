package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockWorkspaceManager is a mock implementation for testing
type mockWorkspaceManager struct {
	currentName  string
	currentPath  string
	currentError error
}

func (m *mockWorkspaceManager) Current() (string, string, error) {
	return m.currentName, m.currentPath, m.currentError
}

func TestPath(t *testing.T) {
	// Create a temp workspace directory
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	// Create a mock workspace manager
	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	path, err := store.Path("test-note")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(workspaceDir, "test-note.md"), path)
}

func TestPath_NoWorkspace(t *testing.T) {
	// Create a mock workspace manager that returns an error
	wm := &mockWorkspaceManager{
		currentName:  "",
		currentPath:  "",
		currentError: ErrWorkspaceNotInitialized,
	}

	store := NewStore(wm)

	_, err := store.Path("test-note")
	assert.ErrorIs(t, err, ErrWorkspaceNotInitialized)
}

func TestEnsure_NewFile(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Ensure("new-note", []byte("# New Note\n\nContent here"))
	require.NoError(t, err)

	expectedPath := filepath.Join(workspaceDir, "new-note.md")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err)

	content, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	assert.Equal(t, "# New Note\n\nContent here", string(content))
}

func TestEnsure_AlreadyExists(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	notePath := filepath.Join(workspaceDir, "existing-note.md")
	err = os.WriteFile(notePath, []byte("# Existing Note\n"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Ensure("existing-note", []byte("# Another Note\n"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrExist) || strings.Contains(err.Error(), "already exists"))
}

func TestEnsure_CreatesDirectory(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Ensure("nested/path/note", []byte("# Nested Note\n"))
	require.NoError(t, err)

	expectedPath := filepath.Join(workspaceDir, "nested", "path", "note.md")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err)
}
