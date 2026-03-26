package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMkdirJSON_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "mkdir", "projects", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "projects", response["folder"])
	assert.Contains(t, response, "path")
}

func TestMkdir_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "mkdir", "testfolder")

	assert.Contains(t, out, "Created folder: testfolder")

	_, wsPath, _ := findWorkspacePath(configDir)
	folderPath := filepath.Join(wsPath, "testfolder")
	_, err := os.Stat(folderPath)
	assert.NoError(t, err, "folder should exist on disk")
}

func TestMkdir_Nested_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "mkdir", "projects/2024", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "projects/2024", response["folder"])
}

func TestMkdir_DeepNested(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "mkdir", "a/b/c/d/e/f", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])

	_, wsPath, _ := findWorkspacePath(configDir)
	folderPath := filepath.Join(wsPath, "a/b/c/d/e/f")
	_, err = os.Stat(folderPath)
	assert.NoError(t, err, "deep nested folder should exist on disk")
}

func TestMkdir_AlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "existing")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "mkdir", "existing", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "folder already exists: existing", response["error"])
}

func TestMkdir_NestedAlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "a/b/c")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "mkdir", "a/b/c", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "folder already exists: a/b/c", response["error"])
}

func TestMkdir_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	out := runHotnote(t, "mkdir", "test")

	assert.Contains(t, out, "create workspace manager")
}

func TestMkdir_ExitCode_AlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "existing")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	code := runHotnoteWithExitCode(t, "mkdir", "existing")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestMkdir_ExitCode_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	code := runHotnoteWithExitCode(t, "mkdir", "test")
	if code != 4 {
		t.Errorf("expected exit code 4, got %d", code)
	}
}

func TestMkdir_PrettyJSON(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "mkdir", "prettytest", "--json", "--pretty")

	assert.Contains(t, out, "  \"status\":")
	assert.Contains(t, out, "  \"folder\":")
	assert.Contains(t, out, "  \"path\":")
}

func TestMkdir_PathTraversal_Parent(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "mkdir", "../outside", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid folder path: must be inside workspace", response["error"])
}

func TestMkdir_PathTraversal_Absolute(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "mkdir", "/etc/malicious", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid folder path: must be relative to workspace", response["error"])
}

func TestMkdir_PathTraversal_MidPath(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "mkdir", "a/../../b", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid folder path: must be inside workspace", response["error"])
}
