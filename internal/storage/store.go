package storage

import (
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
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
}