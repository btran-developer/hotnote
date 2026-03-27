package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	binaryPath     string
	binaryBuildErr error
	buildOnce      sync.Once
)

func getBinaryPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	projectRoot := wd
	if filepath.Base(wd) == "cmd" {
		projectRoot = filepath.Dir(wd)
	}

	buildOnce.Do(func() {
		binaryDir, err := os.MkdirTemp("", "hotnote-test-binary-*")
		if err != nil {
			binaryBuildErr = err
			return
		}

		binaryPath = filepath.Join(binaryDir, "hotnote")
		cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/hotnote")
		cmd.Dir = projectRoot

		env := os.Environ()
		hasHome := false
		for _, e := range env {
			if strings.HasPrefix(e, "HOME=") {
				hasHome = true
				break
			}
		}
		if !hasHome {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				env = append(env, "HOME="+homeDir)
			}
		}
		cmd.Env = env
		output, err := cmd.CombinedOutput()
		if err != nil {
			binaryBuildErr = err
			t.Logf("Failed to build hotnote: %v\n%s", err, output)
		}
	})

	if binaryBuildErr != nil {
		t.Fatalf("Failed to build hotnote: %v", binaryBuildErr)
	}

	return binaryPath
}

func runHotnote(t *testing.T, args ...string) string {
	binaryPath := getBinaryPath(t)
	cmd := exec.Command(binaryPath, args...)
	home := os.Getenv("HOME")
	cmd.Env = append(os.Environ(), "HOME="+home)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command failed: %v\nOutput: %s", err, string(out))
	}
	return string(out)
}

func TestWorkspaceInitJSON_Success(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-init-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	t.Setenv("HOME", configDir)

	out := runHotnote(t, "workspace", "init", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Equal(t, "Initialized workspace: default", response["message"])
}

func TestWorkspaceInitJSON_AlreadyInitialized(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "workspace", "init", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "workspace already initialized", response["error"])
}

func TestWorkspaceListJSON_Single(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "workspace", "list", "--json")

	var workspaces []map[string]interface{}
	err := json.Unmarshal([]byte(out), &workspaces)
	require.NoError(t, err)

	require.Len(t, workspaces, 1)

	ws := workspaces[0]
	assert.Equal(t, "default", ws["name"])
	assert.Contains(t, ws, "path")
	assert.Equal(t, true, ws["current"])
}

func TestWorkspaceListJSON_Multiple(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-list-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces")
	err = os.MkdirAll(filepath.Join(workspaceDir, "default"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(workspaceDir, "work"), 0755)
	require.NoError(t, err)

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	err = os.MkdirAll(configDirPath, 0755)
	require.NoError(t, err)

	configContent := "current_workspace: default\nworkspaces:\n  default: " + filepath.Join(workspaceDir, "default") + "\n  work: " + filepath.Join(workspaceDir, "work") + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Setenv("HOME", configDir)
	out := runHotnote(t, "workspace", "list", "--json")

	var workspaces []map[string]interface{}
	err = json.Unmarshal([]byte(out), &workspaces)
	require.NoError(t, err)

	require.Len(t, workspaces, 2)

	currentFound := false
	for _, ws := range workspaces {
		assert.Contains(t, ws, "name")
		assert.Contains(t, ws, "path")
		assert.Contains(t, ws, "current")

		if ws["name"] == "default" {
			assert.Equal(t, true, ws["current"])
			currentFound = true
		} else {
			assert.Equal(t, false, ws["current"])
		}
	}
	assert.True(t, currentFound, "current workspace should be marked")
}

func TestWorkspaceUseJSON_Success(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-use-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces")
	err = os.MkdirAll(filepath.Join(workspaceDir, "default"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(workspaceDir, "testws"), 0755)
	require.NoError(t, err)

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	err = os.MkdirAll(configDirPath, 0755)
	require.NoError(t, err)

	configContent := "current_workspace: default\nworkspaces:\n  default: " + filepath.Join(workspaceDir, "default") + "\n  testws: " + filepath.Join(workspaceDir, "testws") + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Setenv("HOME", configDir)
	out := runHotnote(t, "workspace", "use", "testws", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Equal(t, "Switched to workspace: testws", response["message"])
}

func TestWorkspaceUseJSON_NotFound(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "workspace", "use", "nonexistent", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "workspace not found", response["error"])
}

func TestWorkspaceNewJSON_Success(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-new-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces")
	err = os.MkdirAll(filepath.Join(workspaceDir, "default"), 0755)
	require.NoError(t, err)

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	err = os.MkdirAll(configDirPath, 0755)
	require.NoError(t, err)

	configContent := "current_workspace: default\nworkspaces:\n  default: " + filepath.Join(workspaceDir, "default") + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Setenv("HOME", configDir)
	out := runHotnote(t, "workspace", "new", "newtest", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Equal(t, "Created workspace: newtest", response["message"])
}

func TestWorkspaceNewJSON_MissingName(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "workspace", "new", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "workspace name required", response["error"])
}

func TestWorkspaceNewJSON_AlreadyExists(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "workspace", "new", "default", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "workspace already exists", response["error"])
}

func TestWorkspaceJSON_ValidOutput(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "workspace", "list", "--json")

	var data json.RawMessage
	err := json.Unmarshal([]byte(out), &data)
	assert.NoError(t, err)
}

func TestWorkspaceDeleteJSON_NonExistent(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "workspace", "delete", "nonexistent", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "workspace not found", response["error"])
}

func TestWorkspaceDeleteJSON_DefaultWithoutForce(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	out := runHotnote(t, "workspace", "delete", "default", "--json")

	var response map[string]string
	err := json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "use --force to delete", response["error"])
}

func TestWorkspaceDeleteJSON_NonDefaultCurrent(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-current-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces")
	err = os.MkdirAll(filepath.Join(workspaceDir, "default"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(workspaceDir, "work"), 0755)
	require.NoError(t, err)

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	err = os.MkdirAll(configDirPath, 0755)
	require.NoError(t, err)

	configContent := "current_workspace: work\nworkspaces:\n  default: " + filepath.Join(workspaceDir, "default") + "\n  work: " + filepath.Join(workspaceDir, "work") + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Setenv("HOME", configDir)
	out := runHotnote(t, "workspace", "delete", "work", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "cannot delete current workspace: switch to another workspace first", response["error"])
}

func TestWorkspaceDeleteJSON_Success(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-delete-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces")
	err = os.MkdirAll(filepath.Join(workspaceDir, "default"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(workspaceDir, "work"), 0755)
	require.NoError(t, err)

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	err = os.MkdirAll(configDirPath, 0755)
	require.NoError(t, err)

	configContent := "current_workspace: default\nworkspaces:\n  default: " + filepath.Join(workspaceDir, "default") + "\n  work: " + filepath.Join(workspaceDir, "work") + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Setenv("HOME", configDir)
	out := runHotnote(t, "workspace", "delete", "work", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "status")
	assert.Equal(t, "deleted", response["status"])
	assert.Equal(t, "work", response["workspace"])
}

func TestWorkspaceDeleteDefaultJSON_WithForce(t *testing.T) {
	configDir := setupTestWorkspace(t)
	t.Cleanup(func() { os.RemoveAll(configDir) })

	err := os.WriteFile(filepath.Join(configDir, "workspaces", "default", "test.md"), []byte("# Test"), 0644)
	require.NoError(t, err)

	out := runHotnote(t, "workspace", "delete", "default", "--force", "--json")

	var response map[string]string
	err = json.Unmarshal([]byte(out), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "status")
	assert.Equal(t, "cleared", response["status"])
	assert.Equal(t, "default", response["workspace"])
}
