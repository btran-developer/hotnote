package tui

import (
	"fmt"
	"sort"
	"time"

	"hotnotego/internal/workspace"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	PendingActionSwitchWorkspace = "switch_workspace" // PendingActionSwitchWorkspace indicates a pending workspace switch.
	PendingActionSwitchMain      = "switch_main"      // PendingActionSwitchMain indicates a pending main view switch.
	PendingTimeoutSeconds        = 3                  // PendingTimeoutSeconds is the timeout for pending actions.

	StatusDefaultWorkspace = "view: workspaces · Esc: switch · Tab: switch · ?: help" // StatusDefaultWorkspace is the default status in workspace view.
	StatusDefaultMain      = "view: main · Esc: switch · Tab: switch · ?: help"       // StatusDefaultMain is the default status in main view.
	StatusToMain           = "Esc: to main · Tab: switch · ?: help"                   // StatusToMain is the status when switching to main view.
	StatusToWorkspace      = "Esc: to workspaces · Tab: switch · ?: help"             // StatusToWorkspace is the status when switching to workspace view.
)

// App is the main TUI application.
type App struct {
	*tview.Application
	workspaceMgr *workspace.Manager

	flex              *tview.Flex
	workspaceLayout   *tview.Flex
	treeView          *TreeView
	previewPane       *PreviewPane
	workspaceSelector *WorkspaceSelector
	workspaceInput    *StackedFormInput
	statusBar         *StatusBar

	workspaceRoot string
	currentNote   string
	showRaw       bool
	viewMode      ViewMode

	pendingAction string
	pendingTimer  *time.Timer
}

// New creates a new TUI application.
func New() (*App, error) {
	mgr, err := workspace.NewManager()
	if err != nil {
		return nil, fmt.Errorf("create workspace manager: %w", err)
	}

	app := &App{
		Application:  tview.NewApplication(),
		workspaceMgr: mgr,
		showRaw:      false,
		viewMode:     ViewModeWorkspace,
	}

	Apply()

	if err := app.initWorkspace(); err != nil {
		return nil, fmt.Errorf("init workspace: %w", err)
	}

	app.statusBar = NewStatusBar()

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
	a.viewMode = ViewModeMain

	a.resetPendingAction()
	a.loadMainView()

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

	a.viewMode = ViewModeMain
	a.resetPendingAction()
	a.loadMainView()

	a.SetRoot(a.flex, true)
	a.focusTree()
}

func (a *App) showWorkspaceView() {
	workspaces, _, _ := a.workspaceMgr.List()
	entries := sortWorkspaces(workspaces, a.workspaceRoot)

	a.workspaceSelector = NewWorkspaceSelector(
		entries,
		"",
		a.onWorkspaceSelected,
		a.onWorkspaceSelectorCancel,
	)

	a.workspaceLayout = a.buildWorkspaceLayout()

	a.viewMode = ViewModeWorkspace
	a.resetPendingAction()
	a.SetRoot(a.workspaceLayout, true)
	a.Application.SetFocus(a.workspaceSelector)
}

func (a *App) showMainView() {
	if a.workspaceRoot == "" {
		workspaces, current, _ := a.workspaceMgr.List()
		if current != "" {
			a.workspaceRoot = workspaces[current]
		} else if len(workspaces) > 0 {
			for _, path := range workspaces {
				a.workspaceRoot = path
				break
			}
		}
	}

	a.viewMode = ViewModeMain
	a.resetPendingAction()
	a.loadMainView()

	a.workspaceSelector.SetInputCapture(nil)
	a.workspaceInput.SetInputCapture(nil)
	a.SetRoot(a.flex, true)
	a.focusTree()
}

func (a *App) loadMainView() {
	a.treeView = NewTreeView(a.workspaceRoot, func(path string) {
		if err := a.previewPane.LoadNote(path); err != nil {
			a.previewPane.TextView.SetTextColor(DefaultPalette.Error)
			a.previewPane.TextView.SetText("Error: " + err.Error())
		} else {
			a.previewPane.TextView.SetTextColor(tcell.ColorDefault)
		}
		a.currentNote = path
	})

	a.previewPane = NewPreviewPane()

	mainContent := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.treeView.TreeView, 0, 1, true).
		AddItem(a.previewPane.TextView, 0, 3, false)

	a.flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainContent, 0, 1, true).
		AddItem(a.statusBar.TextView, 1, 0, false)

	a.setupGlobalInputCapture()
}

func (a *App) setupUI() error {
	a.workspaceLayout = a.buildWorkspaceLayout()
	a.SetRoot(a.workspaceLayout, true)

	a.setupGlobalInputCapture()

	return nil
}

func (a *App) buildWorkspaceLayout() *tview.Flex {
	a.workspaceInput = NewStackedFormInput(LabelName, "").
		SetFieldWidth(30).
		SetBorder(true).
		SetLabelColor(DefaultPalette.Accent).
		SetErrorColor(DefaultPalette.Error)

	a.workspaceInput.SetInputCapture(a.handleWorkspaceInputEscape)
	a.workspaceInput.SetDoneFunc(a.handleWorkspaceInputDone)

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

	createPane.SetBorder(true).SetTitle(TitleNewWorkspace)

	workspaceContent := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.workspaceSelector, 0, 1, true).
		AddItem(createPane, 0, 3, false)

	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(workspaceContent, 0, 1, true).
		AddItem(a.statusBar.TextView, 1, 0, false)
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

func (a *App) handleWorkspaceInputEscape(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEscape {
		if a.pendingAction == PendingActionSwitchMain {
			a.showMainView()
		} else {
			input := a.workspaceInput.GetText()
			if input != "" {
				a.workspaceInput.SetText("")
				a.workspaceInput.SetError("")
			} else {
				a.setPendingAction(PendingActionSwitchMain)
			}
		}
		return nil
	}
	return event
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
	a.workspaceSelector.SetBorder(true).SetTitle(TitleWorkspaces)
}

func (a *App) setPendingAction(action string) {
	a.pendingAction = action

	if a.pendingTimer != nil {
		a.pendingTimer.Stop()
	}

	if action == PendingActionSwitchMain {
		a.statusBar.SetStatus(StatusToMain)
	} else if action == PendingActionSwitchWorkspace {
		a.statusBar.SetStatus(StatusToWorkspace)
	}

	a.pendingTimer = time.AfterFunc(PendingTimeoutSeconds*time.Second, func() {
		a.Application.QueueUpdate(func() {
			a.resetPendingAction()
		})
	})
}

func (a *App) resetPendingAction() {
	if a.pendingTimer != nil {
		a.pendingTimer.Stop()
		a.pendingTimer = nil
	}

	if a.viewMode == ViewModeWorkspace {
		a.statusBar.SetStatus(StatusDefaultWorkspace)
	} else if a.viewMode == ViewModeMain {
		a.statusBar.SetStatus(StatusDefaultMain)
	}

	a.pendingAction = ""
}

func (a *App) setupGlobalInputCapture() {
	a.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		action, matched := MatchKey(ContextGlobal, event)
		if !matched {
			return event
		}

		switch action {
		case ActionQuit:
			a.Stop()
			return nil
		case ActionSwitchPane:
			a.toggleFocus()
			return nil
		case ActionCancel:
			if a.viewMode == ViewModeMain {
				if a.pendingAction == PendingActionSwitchWorkspace {
					a.showWorkspaceView()
				} else {
					a.setPendingAction(PendingActionSwitchWorkspace)
				}
				return nil
			}
		}

		return event
	})
}

func (a *App) toggleFocus() {
	if a.viewMode == ViewModeWorkspace {
		currentFocus := a.Application.GetFocus()
		if currentFocus == a.workspaceInput.GetInputField() {
			a.Application.SetFocus(a.workspaceSelector)
		} else if currentFocus == a.workspaceSelector {
			a.Application.SetFocus(a.workspaceInput.GetInputField())
		}
		return
	}

	currentFocus := a.Application.GetFocus()
	if currentFocus == a.treeView.TreeView {
		a.Application.SetFocus(a.previewPane.TextView)
	} else if currentFocus == a.previewPane.TextView {
		a.Application.SetFocus(a.treeView.TreeView)
	}
}

func (a *App) focusTree() {
	if a.treeView != nil {
		a.Application.SetFocus(a.treeView.TreeView)
	}
}

func (a *App) focusPreview() {
	if a.previewPane != nil {
		a.Application.SetFocus(a.previewPane.TextView)
	}
}

// Run starts the TUI application.
func (a *App) Run() error {
	return a.Application.Run()
}
