package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestList_Empty(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	notes, err := store.List()
	require.NoError(t, err)
	assert.Empty(t, notes)
}

func TestList_TopLevel(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.WriteFile(filepath.Join(workspaceDir, "note1.md"), []byte("# Note 1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "note2.md"), []byte("# Note 2"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	notes, err := store.List()
	require.NoError(t, err)
	assert.Len(t, notes, 2)

	slugs := []string{notes[0].Slug, notes[1].Slug}
	assert.Contains(t, slugs, "note1")
	assert.Contains(t, slugs, "note2")
}

func TestList_Nested(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "projects"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "projects", "my-idea.md"), []byte("# My Idea"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	notes, err := store.List()
	require.NoError(t, err)
	assert.Len(t, notes, 1)
	assert.Equal(t, "projects/my-idea", notes[0].Slug)
}

func TestList_Mixed(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "folder1"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "root-note.md"), []byte("# Root Note"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "folder1", "nested.md"), []byte("# Nested"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	notes, err := store.List()
	require.NoError(t, err)
	assert.Len(t, notes, 2)

	slugs := []string{notes[0].Slug, notes[1].Slug}
	assert.Contains(t, slugs, "root-note")
	assert.Contains(t, slugs, "folder1/nested")
}

func TestList_SkipsDirectories(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "empty-folder"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "note.md"), []byte("# Note"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	notes, err := store.List()
	require.NoError(t, err)
	assert.Len(t, notes, 1)
	assert.Equal(t, "note", notes[0].Slug)
}

func TestList_CrTime_FromFrontmatter(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	frontmatter := "---\nid: abc-123\ntitle: Test\ncreated_at: 2025-01-15T10:00:00Z\nupdated_at: 2025-01-15T10:00:00Z\ntags: []\n---\n\n# Test\n"
	err = os.WriteFile(filepath.Join(workspaceDir, "test.md"), []byte(frontmatter), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	notes, err := store.List()
	require.NoError(t, err)
	require.Len(t, notes, 1)

	expectedCrTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	assert.Equal(t, expectedCrTime, notes[0].CrTime)
	assert.NotEqual(t, time.Time{}, notes[0].ModTime)
}

func TestList_CrTime_FallbackToModTime(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.WriteFile(filepath.Join(workspaceDir, "plain.md"), []byte("# No Frontmatter"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	notes, err := store.List()
	require.NoError(t, err)
	require.Len(t, notes, 1)

	assert.Equal(t, notes[0].ModTime, notes[0].CrTime)
}

func TestDelete_Success(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	notePath := filepath.Join(workspaceDir, "to-delete.md")
	err = os.WriteFile(notePath, []byte("# To Delete"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Delete("to-delete")
	require.NoError(t, err)

	_, err = os.Stat(notePath)
	assert.True(t, os.IsNotExist(err))
}

func TestDelete_NotExists(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Delete("nonexistent")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrNotExist), "expected not exists error, got: %v", err)
}

func TestDelete_Nested(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "projects"), 0755)
	require.NoError(t, err)
	notePath := filepath.Join(workspaceDir, "projects", "note.md")
	err = os.WriteFile(notePath, []byte("# Note"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Delete("projects/note")
	require.NoError(t, err)

	_, err = os.Stat(notePath)
	assert.True(t, os.IsNotExist(err))
}

func TestRename_Success(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	oldPath := filepath.Join(workspaceDir, "old-note.md")
	err = os.WriteFile(oldPath, []byte("# Old Note"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Rename("old-note", "new-note")
	require.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	newPath := filepath.Join(workspaceDir, "new-note.md")
	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestRename_SourceNotFound(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Rename("nonexistent", "new-name")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNoteNotFound))
}

func TestRename_DestExists(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.WriteFile(filepath.Join(workspaceDir, "old-note.md"), []byte("# Old"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "existing.md"), []byte("# Existing"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Rename("old-note", "existing")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNoteAlreadyExists))
}

func TestRename_NestedToTop(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "folder"), 0755)
	require.NoError(t, err)
	oldPath := filepath.Join(workspaceDir, "folder", "note.md")
	err = os.WriteFile(oldPath, []byte("# Note"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Rename("folder/note", "note")
	require.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	newPath := filepath.Join(workspaceDir, "note.md")
	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestRename_TopToNested(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	oldPath := filepath.Join(workspaceDir, "note.md")
	err = os.WriteFile(oldPath, []byte("# Note"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Rename("note", "folder/new-note")
	require.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	newPath := filepath.Join(workspaceDir, "folder", "new-note.md")
	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestRename_WithDotMdSuffix(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	oldPath := filepath.Join(workspaceDir, "old.md")
	err = os.WriteFile(oldPath, []byte("# Old"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	err = store.Rename("old.md", "new.md")
	require.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	newPath := filepath.Join(workspaceDir, "new.md")
	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestResolve_DirectPath(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "projects"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "projects", "my-idea.md"), []byte("# My Idea"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	slug, err := store.Resolve("projects/my-idea")
	require.NoError(t, err)
	assert.Equal(t, "projects/my-idea", slug)
}

func TestResolve_DirectPath_NotFound(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	_, err = store.Resolve("nonexistent/path/note")
	assert.ErrorIs(t, err, ErrNoteNotFound)
}

func TestResolve_Recursive_Single(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "folder1"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "folder1", "my-note.md"), []byte("# My Note"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	slug, err := store.Resolve("my-note")
	require.NoError(t, err)
	assert.Equal(t, "folder1/my-note", slug)
}

func TestResolve_Recursive_NotFound(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.WriteFile(filepath.Join(workspaceDir, "note.md"), []byte("# Note"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	_, err = store.Resolve("nonexistent")
	assert.ErrorIs(t, err, ErrNoteNotFound)
}

func TestResolve_Recursive_Multiple(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "folder1"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(workspaceDir, "folder2"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "folder1", "duplicate.md"), []byte("# Note 1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "folder2", "duplicate.md"), []byte("# Note 2"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	_, err = store.Resolve("duplicate")
	assert.ErrorIs(t, err, ErrMultipleMatches)
}

func TestResolve_DirectPath_Nested(t *testing.T) {
	workspaceDir, err := os.MkdirTemp("", "hotnote-workspace-*")
	require.NoError(t, err)
	defer os.RemoveAll(workspaceDir)

	err = os.MkdirAll(filepath.Join(workspaceDir, "folder1"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(workspaceDir, "folder2"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "folder1", "note.md"), []byte("# Note in folder1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workspaceDir, "folder2", "note.md"), []byte("# Note in folder2"), 0644)
	require.NoError(t, err)

	wm := &mockWorkspaceManager{
		currentName:  "default",
		currentPath:  workspaceDir,
		currentError: nil,
	}

	store := NewStore(wm)

	slug, err := store.Resolve("folder1/note")
	require.NoError(t, err)
	assert.Equal(t, "folder1/note", slug)
}
