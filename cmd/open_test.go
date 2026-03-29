package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen_SubfolderResolution(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "subfolder", "mynote.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# My Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "open", "subfolder/mynote", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "opened", response["status"])
	assert.Contains(t, response["path"], "subfolder/mynote.md")
}

func TestOpen_RecursiveResolution(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "folder", "mynote.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# My Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "open", "mynote", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "opened", response["status"])
}

func TestOpen_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "open", "nonexistent", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "note not found", response["error"])
}

func TestOpen_MultipleMatches(t *testing.T) {
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

	out := runHotnote(t, "open", "duplicate", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "multiple matches found")
}

func TestOpen_ExitCode_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "open", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestOpen_ExitCode_MultipleMatches(t *testing.T) {
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

	code := runHotnoteWithExitCode(t, "open", "dup")
	if code != 3 {
		t.Errorf("expected exit code 3, got %d", code)
	}
}

func TestOpen_Alias_op(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "op", "testnote", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "opened", response["status"])
	assert.Contains(t, response["path"], "testnote.md")
}
