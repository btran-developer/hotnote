package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	Root string
}

func NewStore(root string) *Store {
	return &Store{Root: root}
}

func (s *Store) Path(id string) string {
	return filepath.Join(s.Root, id+".md")
}

func (s *Store) Ensure(id string) (*os.File, error) {
	path := s.Path(id)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("ensure: mkdir: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return nil, fmt.Errorf("ensure: open file: %w", err)
	}
	return file, nil
}
