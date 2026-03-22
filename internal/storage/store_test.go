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

	file, err := store.Ensure("new-note")
	require.NoError(t, err)
	require.NotNil(t, file)
	file.Close()

	// Verify the file was created
	expectedPath := filepath.Join(workspaceDir, "new-note.md")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err)

	// Verify content is empty (no frontmatter written by Ensure)
	content, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestEnsure_AlreadyExists(t *testing.T) {
	// Create a temp workspace directory with a pre-existing note
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	// Create the note file
	notePath := filepath.Join(workspaceDir, "existing-note.md")
	err = os.WriteFile(notePath, []byte("# Existing Note\n"), 0644)
	require.NoError(t, err)

	// Create a mock workspace manager
	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	// Try to ensure a note that already exists
	_, err = store.Ensure("existing-note")
	assert.Error(t, err)
	// The error is wrapped, check that it contains "already exists"
	assert.True(t, errors.Is(err, os.ErrExist) || strings.Contains(err.Error(), "already exists"))
}

func TestEnsure_CreatesDirectory(t *testing.T) {
	// Create a temp workspace directory (empty)
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

	file, err := store.Ensure("nested/path/note")
	require.NoError(t, err)
	require.NotNil(t, file)
	file.Close()

	// Verify the file was created with nested directories
	expectedPath := filepath.Join(workspaceDir, "nested", "path", "note.md")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err)
}
