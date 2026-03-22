package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func getBinaryPathExitCode(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if filepath.Base(wd) == "cmd" {
		return filepath.Join("..", "hotnote")
	}
	return "./hotnote"
}

func runHotnoteWithExitCode(t *testing.T, args ...string) int {
	binaryPath := getBinaryPathExitCode(t)
	cmd := exec.Command(binaryPath, args...)
	home := os.Getenv("HOME")
	cmd.Env = append(os.Environ(), "HOME="+home)

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}
	return 0
}

func setupTestConfigExitCode(t *testing.T) string {
	configDir, err := os.MkdirTemp("", "hotnote-test-exitcode-*")
	if err != nil {
		t.Fatal(err)
	}

	workspaceDir := filepath.Join(configDir, "workspaces", "default")
	err = os.MkdirAll(workspaceDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	err = os.MkdirAll(configDirPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	configContent := "current_workspace: default\nworkspaces:\n  default: " + workspaceDir + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("HOME", configDir)
	return configDir
}

func TestExitCode_workspaceUse_NotFound(t *testing.T) {
	configDir := setupTestConfigExitCode(t)
	defer os.RemoveAll(configDir)

	code := runHotnoteWithExitCode(t, "workspace", "use", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestExitCode_render_NotFound(t *testing.T) {
	configDir := setupTestConfigExitCode(t)
	defer os.RemoveAll(configDir)

	code := runHotnoteWithExitCode(t, "render", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestExitCode_open_NotFound(t *testing.T) {
	configDir := setupTestConfigExitCode(t)
	defer os.RemoveAll(configDir)

	code := runHotnoteWithExitCode(t, "open", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestExitCode_list_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	os.Unsetenv("HOME")
	defer os.Setenv("HOME", oldHome)

	code := runHotnoteWithExitCode(t, "list")
	if code != 4 {
		t.Errorf("expected exit code 4, got %d", code)
	}
}

func TestExitCode_new_MissingArg(t *testing.T) {
	configDir := setupTestConfigExitCode(t)
	defer os.RemoveAll(configDir)

	code := runHotnoteWithExitCode(t, "new")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}
