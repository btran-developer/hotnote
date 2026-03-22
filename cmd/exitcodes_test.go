package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runHotnoteWithExitCode(t *testing.T, args ...string) int {
	binaryPath := getBinaryPath(t)
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

func TestExitCode_workspaceUse_NotFound(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-exitcode-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces", "default")
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		t.Fatal(err)
	}

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := "current_workspace: default\nworkspaces:\n  default: " + workspaceDir + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", configDir)

	code := runHotnoteWithExitCode(t, "workspace", "use", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestExitCode_render_NotFound(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-exitcode-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces", "default")
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		t.Fatal(err)
	}

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := "current_workspace: default\nworkspaces:\n  default: " + workspaceDir + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", configDir)

	code := runHotnoteWithExitCode(t, "render", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestExitCode_open_NotFound(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-exitcode-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces", "default")
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		t.Fatal(err)
	}

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := "current_workspace: default\nworkspaces:\n  default: " + workspaceDir + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", configDir)

	code := runHotnoteWithExitCode(t, "open", "nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestExitCode_list_NoWorkspace(t *testing.T) {
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	})
	os.Unsetenv("HOME")

	code := runHotnoteWithExitCode(t, "list")
	if code != 4 {
		t.Errorf("expected exit code 4, got %d", code)
	}
}

func TestExitCode_new_MissingArg(t *testing.T) {
	configDir, err := os.MkdirTemp("", "hotnote-test-exitcode-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(configDir) })

	workspaceDir := filepath.Join(configDir, "workspaces", "default")
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		t.Fatal(err)
	}

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := "current_workspace: default\nworkspaces:\n  default: " + workspaceDir + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", configDir)

	code := runHotnoteWithExitCode(t, "new")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}
