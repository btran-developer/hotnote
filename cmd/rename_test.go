package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRename_Success_JSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldNotePath := filepath.Join(wsPath, "old-note.md")
	err = os.WriteFile(oldNotePath, []byte("# Old Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "old-note", "New Note", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "renamed", response["status"])
	assert.Equal(t, "old-note", response["old_slug"])
	assert.Equal(t, "new-note", response["new_slug"])

	_, err = os.Stat(oldNotePath)
	assert.True(t, os.IsNotExist(err), "old note should not exist")

	newNotePath := filepath.Join(wsPath, "new-note.md")
	_, err = os.Stat(newNotePath)
	assert.NoError(t, err, "new note should exist")
}

func TestRename_Success_Human(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldNotePath := filepath.Join(wsPath, "old-note.md")
	err = os.WriteFile(oldNotePath, []byte("# Old Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "old-note", "New Note", "--force")

	assert.Contains(t, out, "Renamed note: old-note → new-note")

	_, err = os.Stat(oldNotePath)
	assert.True(t, os.IsNotExist(err), "old note should not exist")

	newNotePath := filepath.Join(wsPath, "new-note.md")
	_, err = os.Stat(newNotePath)
	assert.NoError(t, err, "new note should exist")
}

func TestRename_NestedPath(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldNotePath := filepath.Join(wsPath, "folder", "subfolder", "old-note.md")
	err = os.MkdirAll(filepath.Dir(oldNotePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(oldNotePath, []byte("# Old Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "folder/subfolder/old-note", "New Note", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "renamed", response["status"])
	assert.Equal(t, "folder/subfolder/old-note", response["old_slug"])
	assert.Equal(t, "folder/subfolder/new-note", response["new_slug"])

	_, err = os.Stat(oldNotePath)
	assert.True(t, os.IsNotExist(err), "old note should not exist")

	newNotePath := filepath.Join(wsPath, "folder", "subfolder", "new-note.md")
	_, err = os.Stat(newNotePath)
	assert.NoError(t, err, "new note should exist")
}

func TestRename_CrossFolder(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldNotePath := filepath.Join(wsPath, "folder1", "note.md")
	err = os.MkdirAll(filepath.Dir(oldNotePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(oldNotePath, []byte("# Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "folder1/note", "folder2/new-note", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "renamed", response["status"])
	assert.Equal(t, "folder1/note", response["old_slug"])
	assert.Equal(t, "folder2/new-note", response["new_slug"])

	_, err = os.Stat(oldNotePath)
	assert.True(t, os.IsNotExist(err), "old note should not exist")

	newNotePath := filepath.Join(wsPath, "folder2", "new-note.md")
	_, err = os.Stat(newNotePath)
	assert.NoError(t, err, "new note should exist")
}

func TestRename_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "rename", "nonexistent", "new-name", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "note not found: nonexistent", response["error"])
}

func TestRename_ExitCode_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "rename", "nonexistent", "new-name")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestRename_MultipleMatches(t *testing.T) {
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

	out := runHotnote(t, "rename", "duplicate", "new-name", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "multiple notes match")
}

func TestRename_DestExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	oldNotePath := filepath.Join(wsPath, "old-note.md")
	err = os.WriteFile(oldNotePath, []byte("# Old"), 0644)
	require.NoError(t, err)

	newNotePath := filepath.Join(wsPath, "existing.md")
	err = os.WriteFile(newNotePath, []byte("# Existing"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "old-note", "Existing", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "note already exists: existing", response["error"])
}

func TestRename_RequiresForce_JSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "testnote", "newnote", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "use --force to rename", response["error"])
}

func TestRename_SlugUnchanged(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "test-note.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "test-note", "test-note", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "unchanged", response["status"])
	assert.Equal(t, "test-note", response["slug"])
	assert.Equal(t, "slug unchanged", response["message"])

	_, err = os.Stat(notePath)
	assert.NoError(t, err, "note should still exist")
}

func TestRename_PrettyJSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldNotePath := filepath.Join(wsPath, "old-note.md")
	err = os.WriteFile(oldNotePath, []byte("# Old Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "old-note", "new-note", "--force", "--json", "--pretty")

	assert.Contains(t, out, "  \"status\":")
	assert.Contains(t, out, "  \"old_slug\":")
	assert.Contains(t, out, "  \"new_slug\":")
}

func TestRename_SubfolderResolution(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	notePath := filepath.Join(wsPath, "subfolder", "mynote.md")
	err = os.MkdirAll(filepath.Dir(notePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(notePath, []byte("# My Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "subfolder/mynote", "renamed", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "renamed", response["status"])
	assert.Equal(t, "subfolder/mynote", response["old_slug"])
	assert.Equal(t, "subfolder/renamed", response["new_slug"])

	_, err = os.Stat(notePath)
	assert.True(t, os.IsNotExist(err), "old note should not exist")

	newNotePath := filepath.Join(wsPath, "subfolder", "renamed.md")
	_, err = os.Stat(newNotePath)
	assert.NoError(t, err, "new note should exist")
}

func TestRename_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	out := runHotnote(t, "rename", "old-note", "new-note")

	assert.Contains(t, out, "create workspace manager")
}

func TestRename_ExitCode_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	code := runHotnoteWithExitCode(t, "rename", "old-note", "new-note")
	if code != 4 {
		t.Errorf("expected exit code 4, got %d", code)
	}
}

func TestRename_EmptySlug(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rename", "testnote", "!@#$", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid title: produces empty slug", response["error"])
}

func TestRename_EmptySlug_ExitCode(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "testnote.md")
	err = os.WriteFile(notePath, []byte("# Test Note"), 0644)
	require.NoError(t, err)

	code := runHotnoteWithExitCode(t, "rename", "testnote", "!@#$", "--force")
	if code != 3 {
		t.Errorf("expected exit code 3, got %d", code)
	}
}


func TestRename_Alias_rn(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	notePath := filepath.Join(wsPath, "oldname.md")
	err = os.WriteFile(notePath, []byte("# Old Name"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rn", "oldname", "New Name", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "renamed", response["status"])
	assert.Equal(t, "oldname", response["old_slug"])
	assert.Equal(t, "new-name", response["new_slug"])
}
