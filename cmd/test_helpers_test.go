package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupTestWorkspace(t *testing.T) string {
	configDir, err := os.MkdirTemp("", "hotnote-test-workspace-*")
	require.NoError(t, err)

	workspaceDir := filepath.Join(configDir, "workspaces", "default")
	err = os.MkdirAll(workspaceDir, 0755)
	require.NoError(t, err)

	configDirPath := filepath.Join(configDir, ".config", "hotnote")
	err = os.MkdirAll(configDirPath, 0755)
	require.NoError(t, err)

	configContent := "current_workspace: default\nworkspaces:\n  default: " + workspaceDir + "\n"
	configPath := filepath.Join(configDirPath, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Setenv("HOME", configDir)
	return configDir
}

func findWorkspacePath(configDir string) (string, string, error) {
	workspaceDir := filepath.Join(configDir, "workspaces", "default")
	return "default", workspaceDir, nil
}
