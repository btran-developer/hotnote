package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFolderRename_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldPath := filepath.Join(wsPath, "old-folder")
	err = os.MkdirAll(oldPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "rename", "old-folder", "new-folder")

	assert.Contains(t, out, "Renamed folder: old-folder → new-folder")

	newPath := filepath.Join(wsPath, "new-folder")
	_, err = os.Stat(newPath)
	assert.NoError(t, err, "new folder should exist")

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err), "old folder should not exist")
}

func TestFolderRename_JSON_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldPath := filepath.Join(wsPath, "old-folder")
	err = os.MkdirAll(oldPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "rename", "old-folder", "new-folder", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "renamed", response["status"])
	assert.Equal(t, "old-folder", response["old"])
	assert.Equal(t, "new-folder", response["new"])
}

func TestFolderRename_Nested_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldPath := filepath.Join(wsPath, "parent", "old-child")
	err = os.MkdirAll(oldPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "rename", "parent/old-child", "parent/new-child")

	assert.Contains(t, out, "Renamed folder: parent/old-child → parent/new-child")

	newPath := filepath.Join(wsPath, "parent", "new-child")
	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestFolderRename_SourceNotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "rename", "nonexistent", "new-folder", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "folder not found: nonexistent", response["error"])
}

func TestFolderRename_DestinationExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldPath := filepath.Join(wsPath, "old-folder")
	err = os.MkdirAll(oldPath, 0755)
	require.NoError(t, err)
	newPath := filepath.Join(wsPath, "new-folder")
	err = os.MkdirAll(newPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "rename", "old-folder", "new-folder", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "folder already exists: new-folder", response["error"])
}

func TestFolderRename_RejectsWorkspaceRoot(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "rename", ".", "new-folder", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "cannot rename workspace root", response["error"])
}

func TestFolderRename_PathTraversal(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "rename", "../outside", "new-folder", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid folder path: must be inside workspace", response["error"])
}

func TestFolderRename_CreatesParentDirs(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldPath := filepath.Join(wsPath, "old-folder")
	err = os.MkdirAll(oldPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "rename", "old-folder", "parent/new-folder")

	assert.Contains(t, out, "Renamed folder: old-folder → parent/new-folder")

	newPath := filepath.Join(wsPath, "parent", "new-folder")
	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestFolderRename_ExitCode_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "folder", "rename", "nonexistent", "new-folder")
	assert.Equal(t, 2, code) // ExitNotFound
}

func TestFolderRename_ExitCode_InvalidPath(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "folder", "rename", "../outside", "new-folder")
	assert.Equal(t, 3, code) // ExitInvalidInput
}
