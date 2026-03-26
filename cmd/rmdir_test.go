package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRmdirJSON_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "emptydir")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "rmdir", "emptydir", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])
	assert.Equal(t, "emptydir", response["folder"])

	_, err = os.Stat(folderPath)
	assert.True(t, os.IsNotExist(err), "folder should be deleted")
}

func TestRmdir_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "testfolder")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "rmdir", "testfolder")

	assert.Contains(t, out, "Deleted folder: testfolder")

	_, err = os.Stat(folderPath)
	assert.True(t, os.IsNotExist(err), "folder should be deleted")
}

func TestRmdir_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "rmdir", "nonexistent", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "folder not found: nonexistent", response["error"])
}

func TestRmdir_ExitCode_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "rmdir", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestRmdir_NonEmpty_WithForce(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "nonempty")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(folderPath, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rmdir", "nonempty", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])

	_, err = os.Stat(folderPath)
	assert.True(t, os.IsNotExist(err), "folder should be deleted")
}

func TestRmdir_NonEmpty_RequiresForce(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "nonempty")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(folderPath, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rmdir", "nonempty", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "folder not empty: use --force to delete", response["error"])
}

func TestRmdir_NonEmpty_Force_Human(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "nonempty")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(folderPath, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rmdir", "nonempty", "--force")

	assert.Contains(t, out, "Deleted folder: nonempty")

	_, err = os.Stat(folderPath)
	assert.True(t, os.IsNotExist(err), "folder should be deleted")
}

func TestRmdir_NonEmpty_VerifyContents(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "nested")
	err = os.MkdirAll(filepath.Join(folderPath, "subdir1", "subdir2"), 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(folderPath, "file1.txt"), []byte("content1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(folderPath, "subdir1", "file2.txt"), []byte("content2"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(folderPath, "subdir1", "subdir2", "file3.txt"), []byte("content3"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "rmdir", "nested", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])

	_, err = os.Stat(folderPath)
	assert.True(t, os.IsNotExist(err), "folder and all contents should be deleted")

	_, err = os.Stat(filepath.Join(wsPath, "nested", "subdir1"))
	assert.True(t, os.IsNotExist(err), "subdirectories should be deleted")
}

func TestRmdir_CannotDeleteRoot(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "rmdir", ".", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "cannot delete workspace root", response["error"])
}

func TestRmdir_ExitCode_CannotDeleteRoot(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	code := runHotnoteWithExitCode(t, "rmdir", ".")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRmdir_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	out := runHotnote(t, "rmdir", "test")

	assert.Contains(t, out, "create workspace manager")
}

func TestRmdir_ExitCode_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	code := runHotnoteWithExitCode(t, "rmdir", "test")
	if code != 4 {
		t.Errorf("expected exit code 4, got %d", code)
	}
}

func TestRmdir_Nested_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "projects/2024")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "rmdir", "projects/2024", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])
	assert.Equal(t, "projects/2024", response["folder"])
}

func TestRmdir_PrettyJSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "prettytest")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "rmdir", "prettytest", "--json", "--pretty")

	assert.Contains(t, out, "  \"status\":")
	assert.Contains(t, out, "  \"folder\":")
}

func TestRmdir_PathTraversal_Parent(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "rmdir", "../outside", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid folder path: must be inside workspace", response["error"])
}

func TestRmdir_PathTraversal_Absolute(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "rmdir", "/tmp/outside", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid folder path: must be relative to workspace", response["error"])
}
