package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete_WithForce_JSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "delete", "testnote", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])

	_, err = os.Stat(notePath)
	assert.True(t, os.IsNotExist(err), "note should be deleted")
}

func TestDelete_WithForce_Human(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "delete", "testnote", "--force")

	assert.Contains(t, out, "Deleted note: testnote")

	_, err = os.Stat(notePath)
	assert.True(t, os.IsNotExist(err), "note should be deleted")
}

func TestDelete_NestedPath(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "folder", "subfolder", "testnote.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "delete", "folder/subfolder/testnote", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])

	_, err = os.Stat(notePath)
	assert.True(t, os.IsNotExist(err), "note should be deleted")
}

func TestDelete_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "delete", "nonexistent", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "note not found", response["error"])
}

func TestDelete_ExitCode_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "delete", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestDelete_MultipleMatches(t *testing.T) {
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

	out := runHotnote(t, "delete", "duplicate", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "multiple matches found")
}

func TestDelete_RequiresForce_JSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "delete", "testnote", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "use --force to delete", response["error"])
}

func TestDelete_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	out := runHotnote(t, "delete", "testnote")

	assert.Contains(t, out, "workspace not initialized")
}

func TestDelete_ExitCode_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	code := runHotnoteWithExitCode(t, "delete", "testnote")
	if code != 4 {
		t.Errorf("expected exit code 4, got %d", code)
	}
}

func TestDelete_PrettyJSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "delete", "testnote", "--force", "--json", "--pretty")

	assert.Contains(t, out, "  \"status\":")
	assert.Contains(t, out, "  \"slug\":")
}

func TestDelete_SubfolderResolution(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "subfolder", "mynote.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# My Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "delete", "subfolder/mynote", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])

	_, err = os.Stat(notePath)
	assert.True(t, os.IsNotExist(err), "note should be deleted")
}

func TestDelete_Alias_del(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "del", "testnote", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])
	assert.Equal(t, "testnote", response["slug"])
}
