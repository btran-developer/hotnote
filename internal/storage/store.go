package storage

import (
	"errors"
	"fmt"
	"path/filepath"

	"hotnotego/internal/fsutil"
)

var ErrWorkspaceNotInitialized = errors.New("workspace not initialized")

type WorkspaceManager interface {
	Current() (name string, path string, err error)
}

type Store struct {
	wm WorkspaceManager
}

func NewStore(wm WorkspaceManager) *Store {
	return &Store{wm: wm}
}

func (s *Store) Path(id string) (string, error) {
	_, workspacePath, err := s.wm.Current()
	if err != nil {
		return "", fmt.Errorf("storage: get current workspace: %w", err)
	}
	return filepath.Join(workspacePath, id+".md"), nil
}

func (s *Store) Ensure(id string, content []byte) error {
	path, err := s.Path(id)
	if err != nil {
		return err
	}
	if err := fsutil.AtomicWriteExclusive(path, content, 0644); err != nil {
		return fmt.Errorf("ensure: %w", err)
	}
	return nil
}
