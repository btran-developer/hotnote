package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender_SubfolderResolution(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "subfolder", "mynote.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# My Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "render", "subfolder/mynote", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response["content"], "My Note")
}

func TestRender_RecursiveResolution(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "folder", "mynote.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# My Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "render", "mynote", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response["content"], "My Note")
}

func TestRender_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "render", "nonexistent", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "note not found", response["error"])
}

func TestRender_MultipleMatches(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath1 := filepath.Join(wsPath, "folder1", "duplicate.md")
	err = os.MkdirAll(filepath.Dir(notePath1), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath1, []byte("# Note 1"), 0644)
	require.NoError(t, err)

	notePath2 := filepath.Join(wsPath, "folder2", "duplicate.md")
	err = os.MkdirAll(filepath.Dir(notePath2), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath2, []byte("# Note 2"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "render", "duplicate", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "multiple matches found")
}

func TestRender_ExitCode_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "render", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestRender_ExitCode_MultipleMatches(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath1 := filepath.Join(wsPath, "folder1", "dup.md")
	err = os.MkdirAll(filepath.Dir(notePath1), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath1, []byte("# Note 1"), 0644)
	require.NoError(t, err)

	notePath2 := filepath.Join(wsPath, "folder2", "dup.md")
	err = os.MkdirAll(filepath.Dir(notePath2), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath2, []byte("# Note 2"), 0644)
	require.NoError(t, err)

	code := runHotnoteWithExitCode(t, "render", "dup")
	if code != 3 {
		t.Errorf("expected exit code 3, got %d", code)
	}
}

func TestRender_Alias_rdr(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	content := "# Test Note\n\nContent"
	err = os.WriteFile(notePath, []byte(content), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rdr", "testnote", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response["content"], "Test Note")
}
