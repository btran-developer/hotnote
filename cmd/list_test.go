package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList_JSON_Subfolders(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "projects", "my-idea.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# My Idea"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "list", "--json")

	var response []map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Len(t, response, 1)
	assert.Equal(t, "projects/my-idea", response[0]["slug"])
}

func TestList_Human_Subfolders(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "projects", "my-idea.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# My Idea"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "list")

	assert.Contains(t, out, "projects/my-idea")
}

func TestList_JSON_Mixed(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(wsPath, "folder1"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(wsPath, "root-note.md"), []byte("# Root Note"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(wsPath, "folder1", "nested.md"), []byte("# Nested"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "list", "--json")

	var response []map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
}

func TestList_JSON_Empty(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "list", "--json")

	var response []map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Len(t, response, 0)
}

func TestList_PrettyJSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "list", "--json", "--pretty")

	assert.Contains(t, out, "  \"slug\":")
}

func TestList_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	out := runHotnote(t, "list")

	assert.Contains(t, out, "create workspace manager")
}
