package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_WithPath_JSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "new", "My Note", "--path", "projects", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "projects/my-note", response["slug"])
}

func TestNew_WithPath_Human(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "new", "My Note", "--path", "projects")

	assert.Contains(t, out, "Created note: projects/my-note")
}

func TestNew_WithNestedPath(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	out := runHotnote(t, "new", "My Note", "--path", "projects/todo", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "projects/todo/my-note", response["slug"])

	expectedPath := filepath.Join(wsPath, "projects", "todo", "my-note.md")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err, "note should exist at expected path")
}

func TestNew_WithPath_AlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "projects", "existing.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# Existing"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "new", "Existing", "--path", "projects", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "note already exists")
}

func TestNew_WithPath_NestedDuplicate(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "folder", "note.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "new", "Note", "--path", "folder", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "note already exists: folder/note")
}

func TestNew_WithPath_EmptySlug(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "new", "!@#$", "--path", "folder", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid title: produces empty slug", response["error"])
}

func TestNew_WithPath_ExitCode_AlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "projects", "existing.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# Existing"), 0644)
	require.NoError(t, err)

	code := runHotnoteWithExitCode(t, "new", "Existing", "--path", "projects")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestNew_WithPath_ExitCode_EmptySlug(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "new", "!@#$", "--path", "folder")
	if code != 3 {
		t.Errorf("expected exit code 3, got %d", code)
	}
}
