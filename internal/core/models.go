// Package core defines the domain models for hotnote.
package core

import "time"

// Note represents a markdown note with its metadata.
type Note struct {
	ID        string    // ID is the unique identifier derived from the filename.
	Title     string    // Title is the display title extracted from the filename.
	Path      string    // Path is the absolute filesystem path to the note file.
	Tags      []string  // Tags are extracted from YAML frontmatter.
	CreatedAt time.Time // CreatedAt is the file creation time.
	UpdatedAt time.Time // UpdatedAt is the last content modification time.
}

// Workspace represents a collection of notes at a filesystem path.
type Workspace struct {
	RootPath string // RootPath is the absolute path to the workspace root directory.
}
