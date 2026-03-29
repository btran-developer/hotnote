// Package errors defines typed errors for hotnote CLI operations.
package errors

import (
	"errors"
	"fmt"
)

// WithContext appends context to a base error message for debugging.
func WithContext(base error, context string) string {
	return fmt.Sprintf("%s: %s", base.Error(), context)
}

var (
	// ErrWorkspaceNotInit is returned when no workspace is initialized.
	ErrWorkspaceNotInit = errors.New("workspace not initialized")
	// ErrWorkspaceInit is returned when workspace initialization fails.
	ErrWorkspaceInit = errors.New("workspace initialization failed")
	// ErrWorkspaceList is returned when listing workspaces fails.
	ErrWorkspaceList = errors.New("list workspaces failed")
	// ErrWorkspaceUse is returned when switching to a workspace fails.
	ErrWorkspaceUse = errors.New("use workspace failed")
	// ErrWorkspaceCreate is returned when creating a workspace fails.
	ErrWorkspaceCreate = errors.New("create workspace failed")
	// ErrWorkspaceDelete is returned when deleting a workspace fails.
	ErrWorkspaceDelete = errors.New("delete workspace failed")
	// ErrWorkspaceRename is returned when renaming a workspace fails.
	ErrWorkspaceRename = errors.New("rename workspace failed")
	// ErrWorkspaceNotFound is returned when the requested workspace does not exist.
	ErrWorkspaceNotFound = errors.New("workspace not found")
	// ErrWorkspaceExists is returned when a workspace with the same name already exists.
	ErrWorkspaceExists = errors.New("workspace already exists")
	// ErrWorkspaceIsCurrent is returned when attempting to modify the current workspace.
	ErrWorkspaceIsCurrent = errors.New("cannot modify current workspace")

	// ErrNoteNotFound is returned when the requested note does not exist.
	ErrNoteNotFound = errors.New("note not found")
	// ErrNoteExists is returned when a note with the same slug already exists.
	ErrNoteExists = errors.New("note already exists")
	// ErrNoteCreate is returned when creating a note fails.
	ErrNoteCreate = errors.New("create note failed")
	// ErrNoteDelete is returned when deleting a note fails.
	ErrNoteDelete = errors.New("delete note failed")
	// ErrNoteRename is returned when renaming a note fails.
	ErrNoteRename = errors.New("rename note failed")
	// ErrNoteResolve is returned when resolving a note slug fails.
	ErrNoteResolve = errors.New("resolve note failed")
	// ErrNoteRender is returned when rendering a note's markdown fails.
	ErrNoteRender = errors.New("render note failed")
	// ErrNoteOpen is returned when opening a note for editing fails.
	ErrNoteOpen = errors.New("open note failed")
	// ErrNotePath is returned when retrieving a note's file path fails.
	ErrNotePath = errors.New("get note path failed")
	// ErrNoteList is returned when listing notes fails.
	ErrNoteList = errors.New("list notes failed")
	// ErrMultipleMatches is returned when a slug matches multiple notes.
	ErrMultipleMatches = errors.New("multiple matches found")

	// ErrFolderNotFound is returned when the requested folder does not exist.
	ErrFolderNotFound = errors.New("folder not found")
	// ErrFolderExists is returned when a folder with the same path already exists.
	ErrFolderExists = errors.New("folder already exists")
	// ErrFolderCreate is returned when creating a folder fails.
	ErrFolderCreate = errors.New("create folder failed")
	// ErrFolderDelete is returned when deleting a folder fails.
	ErrFolderDelete = errors.New("delete folder failed")
	// ErrFolderRename is returned when renaming a folder fails.
	ErrFolderRename = errors.New("rename folder failed")
	// ErrFolderRead is returned when reading a folder's contents fails.
	ErrFolderRead = errors.New("read folder failed")
	// ErrInvalidFolderPath is returned when the folder path is invalid or outside the workspace.
	ErrInvalidFolderPath = errors.New("invalid folder path")

	// ErrInvalidSlug is returned when the provided slug is invalid.
	ErrInvalidSlug = errors.New("invalid slug")
	// ErrEmptySlug is returned when slug generation produces an empty value.
	ErrEmptySlug = errors.New("slug produces empty value")

	// ErrMarshalJSON is returned when JSON marshaling fails.
	ErrMarshalJSON = errors.New("marshal JSON failed")
	// ErrOpenEditor is returned when launching the text editor fails.
	ErrOpenEditor = errors.New("open editor failed")
)
