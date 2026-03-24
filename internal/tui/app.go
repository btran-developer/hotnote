package tui

import (
	"fmt"

	"github.com/rivo/tview"
	"hotnotego/internal/workspace"
)

// App represents the main TUI application
type App struct {
	*tview.Application
	workspaceMgr *workspace.Manager

	// UI Components
	flex    *tview.Flex
	tree    *tview.TreeView
	preview *tview.TextView

	// State
	workspaceRoot string
	currentNote   string
	showRaw       bool
}

// New creates a new TUI application
func New() (*App, error) {
	mgr, err := workspace.NewManager()
	if err != nil {
		return nil, fmt.Errorf("create workspace manager: %w", err)
	}

	app := &App{
		Application:  tview.NewApplication(),
		workspaceMgr: mgr,
		showRaw:      false,
	}

	if err := app.setupUI(); err != nil {
		return nil, fmt.Errorf("setup UI: %w", err)
	}

	return app, nil
}

// setupUI initializes the UI components
func (a *App) setupUI() error {
	// Create TreeView for file browser
	a.tree = tview.NewTreeView()

	// Create TextView for preview
	a.preview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)

	// Create Flex layout
	a.flex = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.tree, 0, 1, true).
		AddItem(a.preview, 0, 3, false)

	a.SetRoot(a.flex, true)

	return nil
}

// Run starts the TUI application
func (a *App) Run() error {
	return a.Application.Run()
}
