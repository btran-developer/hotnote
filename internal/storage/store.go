// Package storage manages note files on disk.
package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hotnotego/internal/frontmatter"
	"hotnotego/internal/fsutil"
)

var (
	// ErrWorkspaceNotInitialized is returned when no workspace is set.
	ErrWorkspaceNotInitialized = errors.New("workspace not initialized")
	// ErrNoteNotFound is returned when a note slug does not match any file.
	ErrNoteNotFound = errors.New("note not found")
	// ErrMultipleMatches is returned when a slug is ambiguous.
	ErrMultipleMatches = errors.New("multiple notes match slug")
	// ErrNoteAlreadyExists is returned when trying to create a note that already exists.
	ErrNoteAlreadyExists = errors.New("note already exists")
)

// WorkspaceManager provides access to the current workspace.
type WorkspaceManager interface {
	// Current returns the name and filesystem path of the active workspace.
	Current() (name string, path string, err error)
}

// Store manages note files within a workspace.
type Store struct {
	wm WorkspaceManager
}

// NewStore creates a Store using the given WorkspaceManager.
func NewStore(wm WorkspaceManager) *Store {
	return &Store{wm: wm}
}

// Path returns the filesystem path for note id.
func (s *Store) Path(id string) (string, error) {
	_, workspacePath, err := s.wm.Current()
	if err != nil {
		return "", fmt.Errorf("storage: get current workspace: %w", err)
	}
	return filepath.Join(workspacePath, id+".md"), nil
}

// Ensure writes content for note id, creating the file if it doesn't exist.
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

// NoteInfo holds metadata about a stored note.
type NoteInfo struct {
	Slug    string    // Slug is the URL-safe identifier derived from the filename.
	RelPath string    // RelPath is the path relative to the workspace root.
	CrTime  time.Time // CrTime is the creation time derived from frontmatter, falls back to ModTime.
	ModTime time.Time // ModTime is the last content modification time.
}

// List returns all notes in the current workspace.
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
		crTime := parseCreatedTime(path, info.ModTime())
		notes = append(notes, NoteInfo{
			Slug:    slug,
			RelPath: relPath,
			CrTime:  crTime,
			ModTime: info.ModTime(),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list: walk directory: %w", err)
	}

	return notes, nil
}

// Delete removes the note file for id.
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

// Rename moves a note from oldID to newID.
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

// Resolve resolves input to a unique note slug.
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

// parseCreatedTime reads frontmatter created_at from a file, falling back to modTime.
func parseCreatedTime(path string, modTime time.Time) time.Time {
	f, err := os.Open(path)
	if err != nil {
		return modTime
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return modTime
	}

	data, ok := frontmatter.Extract(buf[:n])
	if !ok {
		return modTime
	}

	crTime, ok := frontmatter.ParseCreatedAt(data)
	if !ok {
		return modTime
	}

	return crTime
}
