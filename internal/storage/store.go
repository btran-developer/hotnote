package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hotnotego/internal/fsutil"
)

var (
	ErrWorkspaceNotInitialized = errors.New("workspace not initialized")
	ErrNoteNotFound            = errors.New("note not found")
	ErrMultipleMatches         = errors.New("multiple notes match slug")
	ErrNoteAlreadyExists       = errors.New("note already exists")
)

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
	id = strings.TrimSuffix(id, ".md")
	path, err := s.Path(id)
	if err != nil {
		return err
	}
	if err := fsutil.AtomicWriteExclusive(path, content, 0644); err != nil {
		return fmt.Errorf("ensure: %w", err)
	}
	return nil
}

type NoteInfo struct {
	Slug    string
	RelPath string
	ModTime time.Time
}

func (s *Store) List() ([]NoteInfo, error) {
	_, workspacePath, err := s.wm.Current()
	if err != nil {
		return nil, fmt.Errorf("list: get current workspace: %w", err)
	}

	var notes []NoteInfo
	err = filepath.WalkDir(workspacePath, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		relPath, err := filepath.Rel(workspacePath, path)
		if err != nil {
			return err
		}
		slug := strings.TrimSuffix(relPath, ".md")
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("list: get file info: %w", err)
		}
		notes = append(notes, NoteInfo{
			Slug:    slug,
			RelPath: relPath,
			ModTime: info.ModTime(),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list: walk directory: %w", err)
	}

	return notes, nil
}

func (s *Store) Delete(id string) error {
	id = strings.TrimSuffix(id, ".md")
	path, err := s.Path(id)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (s *Store) Rename(oldID, newID string) error {
	oldID = strings.TrimSuffix(oldID, ".md")
	newID = strings.TrimSuffix(newID, ".md")

	srcPath, err := s.Path(oldID)
	if err != nil {
		return err
	}

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("rename: %w", ErrNoteNotFound)
	}

	dstPath, err := s.Path(newID)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dstPath); err == nil {
		return fmt.Errorf("rename: %w", ErrNoteAlreadyExists)
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("rename: create directory: %w", err)
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	return nil
}

func (s *Store) Resolve(input string) (string, error) {
	if strings.Contains(input, "/") {
		slug := strings.TrimSuffix(input, ".md")
		path, err := s.Path(slug)
		if err != nil {
			return "", err
		}
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return "", ErrNoteNotFound
			}
			return "", fmt.Errorf("resolve: stat: %w", err)
		}
		if info.IsDir() {
			return "", ErrNoteNotFound
		}
		return slug, nil
	}

	notes, err := s.List()
	if err != nil {
		return "", fmt.Errorf("resolve: list notes: %w", err)
	}

	baseName := strings.TrimSuffix(input, ".md")
	var matches []string
	for _, n := range notes {
		if filepath.Base(n.Slug) == baseName {
			matches = append(matches, n.Slug)
		}
	}

	if len(matches) == 0 {
		return "", ErrNoteNotFound
	}
	if len(matches) > 1 {
		return "", ErrMultipleMatches
	}

	return matches[0], nil
}
