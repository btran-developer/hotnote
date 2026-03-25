package tui

type ViewMode string

const (
	ViewModeWorkspace ViewMode = "workspace"
	ViewModeMain      ViewMode = "main"
)

const (
	LabelName         = "Name:"
	TitleNewWorkspace = " New Workspace "
	TitleWorkspaces   = " Workspaces "

	DefaultWorkspaceName = "default"

	PrefixSelected   = " >"
	PrefixUnselected = "  "

	TreePrefixFolderClosed = "[+]"
	TreePrefixFolderOpen   = "[-]"
	TreePrefixFile         = "*"
)
