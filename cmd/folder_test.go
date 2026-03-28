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

func TestFolderCreate_JSON_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "create", "projects", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "projects", response["folder"])
}

func TestFolderCreate_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "create", "testfolder")

	assert.Contains(t, out, "Created folder: testfolder")

	_, wsPath, _ := findWorkspacePath(configDir)
	folderPath := filepath.Join(wsPath, "testfolder")
	_, err := os.Stat(folderPath)
	assert.NoError(t, err, "folder should exist on disk")
}

func TestFolderCreate_Nested_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "create", "projects/2024", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "projects/2024", response["folder"])
}

func TestFolderCreate_DeepNested(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "create", "a/b/c/d/e/f", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])

	folderPath := filepath.Join(wsPath, "a/b/c/d/e/f")
	_, err = os.Stat(folderPath)
	assert.NoError(t, err, "deep nested folder should exist on disk")
}

func TestFolderCreate_Alias_cr(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "cr", "test-alias", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "test-alias", response["folder"])
}

func TestFolderCreate_Alias_new(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "new", "test-alias-new", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "created", response["status"])
	assert.Equal(t, "test-alias-new", response["folder"])
}

func TestFolderCreate_AlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "existing")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "create", "existing", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "folder already exists: existing", response["error"])
}

func TestFolderDelete_JSON_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "emptydir")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "delete", "emptydir", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])
	assert.Equal(t, "emptydir", response["folder"])

	_, err = os.Stat(folderPath)
	assert.True(t, os.IsNotExist(err), "folder should be deleted")
}

func TestFolderDelete_Success(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "testfolder")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "delete", "testfolder")

	assert.Contains(t, out, "Deleted folder: testfolder")

	_, err = os.Stat(folderPath)
	assert.True(t, os.IsNotExist(err), "folder should be deleted")
}

func TestFolderDelete_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "delete", "nonexistent", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "folder not found: nonexistent", response["error"])
}

func TestFolderDelete_NonEmpty_WithForce(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "nonempty")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(folderPath, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "delete", "nonempty", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])

	_, err = os.Stat(folderPath)
	assert.True(t, os.IsNotExist(err), "folder should be deleted")
}

func TestFolderDelete_Alias_del(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	folderPath := filepath.Join(wsPath, "test-del")
	err = os.MkdirAll(folderPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "del", "test-del", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Equal(t, "deleted", response["status"])
	assert.Equal(t, "test-del", response["folder"])
}

func TestFolderDelete_CannotDeleteRoot(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "delete", ".", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "cannot delete workspace root", response["error"])
}

func TestFolderList_JSON_Empty(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "list", "--json")

	var response []map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Empty(t, response)
}

func TestFolderList_JSON_WithContents(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(wsPath, "folder1"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(wsPath, "folder2"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(wsPath, "note1.md"), []byte("# Note 1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(wsPath, "note2.md"), []byte("# Note 2"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "list", "--json")

	var response []map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Len(t, response, 4)

	var names []string
	for _, item := range response {
		names = append(names, item["name"])
	}
	assert.Contains(t, names, "folder1")
	assert.Contains(t, names, "folder2")
	assert.Contains(t, names, "note1.md")
	assert.Contains(t, names, "note2.md")

	for _, item := range response {
		if item["name"] == "folder1" || item["name"] == "folder2" {
			assert.Equal(t, "folder", item["type"])
		} else {
			assert.Equal(t, "file", item["type"])
		}
	}
}

func TestFolderList_Human_Output(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(wsPath, "myfolder"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(wsPath, "mynote.md"), []byte("# Note"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "list")

	assert.Contains(t, out, "myfolder/")
	assert.Contains(t, out, "mynote.md")
}

func TestFolderList_NestedPath(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(wsPath, "parent", "child"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(wsPath, "parent", "file.md"), []byte("# File"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "list", "parent", "--json")

	var response []map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
}

func TestFolderList_InvalidPath(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "list", "nonexistent", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "path not found")
}

func TestFolderList_PathTraversal(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "folder", "list", "../outside", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid folder path: must be inside workspace", response["error"])
}

func TestFolderList_Alias_ls(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(wsPath, "testfolder"), 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "ls", "--json")

	var response []map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Len(t, response, 1)
	assert.Equal(t, "testfolder", response[0]["name"])
}

func TestFolderRename_Alias_rn(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	_, wsPath, err := findWorkspacePath(configDir)
	require.NoError(t, err)
	oldPath := filepath.Join(wsPath, "old-name")
	err = os.MkdirAll(oldPath, 0755)
	require.NoError(t, err)

	out := runHotnote(t, "folder", "rn", "old-name", "new-name")

	assert.Contains(t, out, "Renamed folder: old-name → new-name")
}
