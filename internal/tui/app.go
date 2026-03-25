package tui

import (
	"fmt"
	"sort"

	"hotnotego/internal/workspace"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	workspaceMgr *workspace.Manager

	flex              *tview.Flex
	tree              *tview.TreeView
	preview           *tview.TextView
	workspaceSelector *WorkspaceSelector
	workspaceInput    *StackedFormInput

	workspaceRoot string
	currentNote   string
	showRaw       bool
	initialized   bool
}

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

	if err := app.initWorkspace(); err != nil {
		return nil, fmt.Errorf("init workspace: %w", err)
	}

	if err := app.setupUI(); err != nil {
		return nil, fmt.Errorf("setup UI: %w", err)
	}

	return app, nil
}

func (a *App) initWorkspace() error {
	workspaces, current, err := a.workspaceMgr.List()
	if err != nil || len(workspaces) == 0 {
		if initErr := a.workspaceMgr.Init(); initErr != nil {
			return fmt.Errorf("workspace corrupted, re-init failed: %w", initErr)
		}
		workspaces, _, err = a.workspaceMgr.List()
		if err != nil {
			return fmt.Errorf("workspace re-init failed: %w", err)
		}
		current = ""
	}

	entries := sortWorkspaces(workspaces, current)
	a.workspaceSelector = NewWorkspaceSelector(
		entries,
		current,
		a.onWorkspaceSelected,
		a.onWorkspaceSelectorCancel,
	)

	return nil
}

func (a *App) onWorkspaceSelected(name, path string) {
	a.workspaceMgr.Use(name)
	a.workspaceRoot = path
	a.workspaceSelector.SetInputCapture(nil)
	a.workspaceInput.SetInputCapture(nil)
	a.SetRoot(a.flex, true)
	a.focusTree()
}

func (a *App) onWorkspaceSelectorCancel() {
	workspaces, current, err := a.workspaceMgr.List()
	if err != nil {
		return
	}

	if current != "" {
		a.workspaceRoot = workspaces[current]
	} else if len(workspaces) > 0 {
		names := make([]string, 0, len(workspaces))
		for name := range workspaces {
			names = append(names, name)
		}
		sort.Strings(names)
		a.workspaceRoot = workspaces[names[0]]
		a.workspaceMgr.Use(names[0])
	}

	a.SetRoot(a.flex, true)
	a.focusTree()
}

func (a *App) setupUI() error {
	a.setupStyles()
	a.createMainView()

	createPane := a.createWorkspaceInputPane()
	a.setupWorkspaceInputHandler()

	mainLayout := a.createMainLayout(createPane)
	a.SetRoot(mainLayout, true)

	a.setupGlobalInputCapture()

	return nil
}

func (a *App) setupStyles() {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault
	tview.Styles.PrimaryTextColor = tcell.ColorDefault
	tview.Styles.SecondaryTextColor = tcell.ColorDefault
}

func (a *App) createMainView() {
	a.tree = tview.NewTreeView()

	a.preview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)

	a.flex = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.tree, 0, 1, true).
		AddItem(a.preview, 0, 3, false)
}

func (a *App) createWorkspaceInputPane() *tview.Flex {
	a.workspaceInput = NewStackedFormInput("Name:", "new workspace name...").
		SetFieldWidth(30).
		SetBorder(true).
		SetErrorColor(tcell.ColorRed)

	horizontalCenter := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(a.workspaceInput, 0, 1, false).
		AddItem(nil, 0, 1, false)

	createPane := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(horizontalCenter, 0, 1, false).
		AddItem(nil, 0, 1, false)

	createPane.SetBorder(true).SetTitle(" New Workspace ")

	createPane.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.workspaceInput.SetText("")
			a.workspaceInput.SetError("")
			return nil
		}
		return event
	})

	return createPane
}

func (a *App) setupWorkspaceInputHandler() {
	a.workspaceInput.SetDoneFunc(a.handleWorkspaceInputDone)
}

func (a *App) handleWorkspaceInputDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	name := a.workspaceInput.GetText()
	if name == "" {
		return
	}

	if err := a.workspaceMgr.New(name, ""); err != nil {
		a.workspaceInput.SetError(fmt.Sprintf("Error: %v", err))
		return
	}

	workspaces, _, _ := a.workspaceMgr.List()
	path := workspaces[name]

	a.updateWorkspaceList(name)

	a.workspaceInput.SetText("")
	a.workspaceInput.SetError("")
	a.onWorkspaceSelected(name, path)
}

func (a *App) updateWorkspaceList(current string) {
	workspaces, _, _ := a.workspaceMgr.List()
	entries := sortWorkspaces(workspaces, current)
	a.workspaceSelector = NewWorkspaceSelector(
		entries,
		current,
		a.onWorkspaceSelected,
		a.onWorkspaceSelectorCancel,
	)
	a.workspaceSelector.SetBorder(true).SetTitle(" Workspaces ")
}

func (a *App) createMainLayout(createPane *tview.Flex) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.workspaceSelector, 0, 1, true).
		AddItem(createPane, 0, 3, false)
}

func (a *App) setupGlobalInputCapture() {
	a.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlQ {
			a.Stop()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			a.toggleFocus()
			return nil
		}
		return event
	})
}

func (a *App) toggleFocus() {
	currentFocus := a.Application.GetFocus()
	if currentFocus == a.workspaceInput.GetInputField() {
		a.Application.SetFocus(a.workspaceSelector)
	} else if currentFocus == a.workspaceSelector {
		a.Application.SetFocus(a.workspaceInput.GetInputField())
	}
}

func (a *App) focusTree() {
	a.Application.SetFocus(a.tree)
}

func (a *App) Run() error {
	return a.Application.Run()
}
