// Package tui implements the terminal user interface for hotnote.
package tui

// ViewMode determines which TUI panel has focus.
type ViewMode string

const (
	ViewModeWorkspace ViewMode = "workspace" // ViewModeWorkspace shows the workspace selector.
	ViewModeMain      ViewMode = "main"      // ViewModeMain shows the file tree and preview.
)

const (
	LabelName         = "Name:"           // LabelName is the label for the name input field.
	TitleNewWorkspace = " New Workspace " // TitleNewWorkspace is the title for the new workspace dialog.
	TitleWorkspaces   = " Workspaces "    // TitleWorkspaces is the title for the workspace list.

	DefaultWorkspaceName = "default" // DefaultWorkspaceName is the name of the default workspace.

	PrefixSelected   = " >" // PrefixSelected is the prefix for the selected item.
	PrefixUnselected = "  " // PrefixUnselected is the prefix for unselected items.

	TreePrefixFolderClosed = "[+]" // TreePrefixFolderClosed is the prefix for collapsed folders.
	TreePrefixFolderOpen   = "[-]" // TreePrefixFolderOpen is the prefix for expanded folders.
	TreePrefixFile         = "*"   // TreePrefixFile is the prefix for files.
)
