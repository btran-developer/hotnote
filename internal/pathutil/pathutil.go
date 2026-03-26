package pathutil

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidPath         = errors.New("invalid folder path: must be relative to workspace")
	ErrPathOutsideWS       = errors.New("invalid folder path: must be inside workspace")
	ErrPathIsWorkspaceRoot = errors.New("cannot delete workspace root")
)

type ValidationResult struct {
	FolderPath    string
	AbsWsPath     string
	AbsFolderPath string
}

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
