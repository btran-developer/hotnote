// Package pathutil validates and resolves workspace-relative paths.
package pathutil

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	// ErrInvalidPath is returned when the path is absolute instead of relative.
	ErrInvalidPath = errors.New("invalid folder path: must be relative to workspace")
	// ErrPathOutsideWS is returned when the resolved path escapes the workspace.
	ErrPathOutsideWS = errors.New("invalid folder path: must be inside workspace")
	// ErrPathIsWorkspaceRoot is returned when the path is the workspace root itself.
	ErrPathIsWorkspaceRoot = errors.New("cannot delete workspace root")
)

// ValidationResult holds the resolved paths after validation.
type ValidationResult struct {
	FolderPath    string // FolderPath is the resolved path relative to workspace.
	AbsWsPath     string // AbsWsPath is the absolute workspace path.
	AbsFolderPath string // AbsFolderPath is the absolute folder path.
}

// ValidateFolderPath validates that folder is a safe, workspace-relative path.
func ValidateFolderPath(wsPath, folder string) (*ValidationResult, error) {
	if filepath.IsAbs(folder) {
		return nil, ErrInvalidPath
	}

	folderPath := filepath.Join(wsPath, folder)

	absWsPath, err := filepath.Abs(wsPath)
	if err != nil {
		return nil, err
	}

	absFolderPath, err := filepath.Abs(folderPath)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(absFolderPath, absWsPath+string(filepath.Separator)) && absFolderPath != absWsPath {
		return nil, ErrPathOutsideWS
	}

	return &ValidationResult{
		FolderPath:    folderPath,
		AbsWsPath:     absWsPath,
		AbsFolderPath: absFolderPath,
	}, nil
}
