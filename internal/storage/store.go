package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"hotnotego/internal/workspace"
)

// Store represents a storage backend for notes
type Store struct {
	wm *workspace.Manager
}

// NewStore creates a new store with the given workspace manager
func NewStore(wm *workspace.Manager) *Store {
	return &Store{wm: wm}
}

// Path returns the full path for a note ID in the current workspace
func (s *Store) Path(id string) (string, error) {
	_, workspacePath, err := s.wm.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(workspacePath, id+".md"), nil
}

// Ensure creates or opens a note file for the given ID in the current workspace
func (s *Store) Ensure(id string) (*os.File, error) {
	path, err := s.Path(id)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("ensure: mkdir: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return nil, fmt.Errorf("ensure: open file: %w", err)
	}
	return file, nil
}
