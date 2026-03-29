package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate_WithPath_JSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "create", "My Note", "--path", "projects", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "projects/my-note", response["slug"])
}

func TestCreate_WithPath_Human(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "create", "My Note", "--path", "projects")

	assert.Contains(t, out, "Created note: projects/my-note")
}

func TestCreate_WithNestedPath(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	out := runHotnote(t, "create", "My Note", "--path", "projects/todo", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "projects/todo/my-note", response["slug"])

	expectedPath := filepath.Join(wsPath, "projects", "todo", "my-note.md")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err, "note should exist at expected path")
}

func TestCreate_WithPath_AlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "projects", "existing.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# Existing"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "create", "Existing", "--path", "projects", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "note already exists")
}

func TestCreate_WithPath_NestedDuplicate(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "folder", "note.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "create", "Note", "--path", "folder", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "note already exists")
}

func TestCreate_WithPath_EmptySlug(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "create", "!@#$", "--path", "folder", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "slug produces empty value", response["error"])
}

func TestCreate_WithPath_ExitCode_AlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "projects", "existing.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# Existing"), 0644)
	require.NoError(t, err)

	code := runHotnoteWithExitCode(t, "create", "Existing", "--path", "projects")
	if code != 3 {
		t.Errorf("expected exit code 3, got %d", code)
	}
}

func TestCreate_WithPath_ExitCode_EmptySlug(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "create", "!@#$", "--path", "folder")
	if code != 3 {
		t.Errorf("expected exit code 3, got %d", code)
	}
}

func TestCreate_Alias_New(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "new", "My Note", "--path", "projects", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "projects/my-note", response["slug"])
}
