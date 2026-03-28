package tui

import (
	"github.com/gdamore/tcell/v2"
)

// KeyContext represents the context in which key bindings are active.
type KeyContext string

const (
	ContextGlobal  KeyContext = "global"  // ContextGlobal is the default context.
	ContextTree    KeyContext = "tree"    // ContextTree is active when the file tree has focus.
	ContextPreview KeyContext = "preview" // ContextPreview is active when the preview pane has focus.
	ContextEditor  KeyContext = "editor"  // ContextEditor is active when the text editor has focus.
)

// KeyAction represents an action triggered by a key binding.
type KeyAction string

const (
	ActionSwitchPane   KeyAction = "switch_pane"   // ActionSwitchPane switches focus between panels.
	ActionQuit         KeyAction = "quit"          // ActionQuit exits the application.
	ActionHelp         KeyAction = "help"          // ActionHelp shows help information.
	ActionCancel       KeyAction = "cancel"        // ActionCancel cancels the current operation.
	ActionScrollUp     KeyAction = "scroll_up"     // ActionScrollUp scrolls up.
	ActionScrollDown   KeyAction = "scroll_down"   // ActionScrollDown scrolls down.
	ActionExpandOpen   KeyAction = "expand_open"   // ActionExpandOpen expands or opens a folder.
	ActionMoveUp       KeyAction = "move_up"       // ActionMoveUp moves selection up.
	ActionMoveDown     KeyAction = "move_down"     // ActionMoveDown moves selection down.
	ActionExternalEdit KeyAction = "external_edit" // ActionExternalEdit opens the note in an external editor.
	ActionNewNote      KeyAction = "new_note"      // ActionNewNote creates a new note.
	ActionNewFolder    KeyAction = "new_folder"    // ActionNewFolder creates a new folder.
	ActionDelete       KeyAction = "delete"        // ActionDelete deletes the selected item.
	ActionRename       KeyAction = "rename"        // ActionRename renames the selected item.
	ActionRefresh      KeyAction = "refresh"       // ActionRefresh refreshes the current view.
	ActionToggleRaw    KeyAction = "toggle_raw"    // ActionToggleRaw toggles raw/preview mode.
	ActionEnterEdit    KeyAction = "enter_edit"    // ActionEnterEdit enters edit mode.
	ActionSave         KeyAction = "save"          // ActionSave saves the current content.
	ActionEditorQuit   KeyAction = "editor_quit"   // ActionEditorQuit quits the editor.
	ActionTabIndent    KeyAction = "tab_indent"    // ActionTabIndent inserts a tab or indentation.
)

type keyBinding struct {
	key    tcell.Key
	rune   rune
	action KeyAction
}

// Keymaps maps each KeyContext to its list of key bindings.
var Keymaps = map[KeyContext][]keyBinding{
	ContextGlobal: {
		{key: tcell.KeyCtrlQ, action: ActionQuit},
		{key: tcell.KeyTab, action: ActionSwitchPane},
		{rune: '?', action: ActionHelp},
		{key: tcell.KeyEscape, action: ActionCancel},
	},
	ContextTree: {
		{key: tcell.KeyUp, action: ActionMoveUp},
		{key: tcell.KeyCtrlP, action: ActionMoveUp},
		{key: tcell.KeyDown, action: ActionMoveDown},
		{key: tcell.KeyCtrlN, action: ActionMoveDown},
		{key: tcell.KeyEnter, action: ActionExpandOpen},
		{key: tcell.KeyTab, action: ActionSwitchPane},
		{rune: 'k', action: ActionMoveUp},
		{rune: 'j', action: ActionMoveDown},
	},
	ContextPreview: {
		{key: tcell.KeyUp, action: ActionScrollUp},
		{key: tcell.KeyCtrlP, action: ActionScrollUp},
		{key: tcell.KeyDown, action: ActionScrollDown},
		{key: tcell.KeyCtrlN, action: ActionScrollDown},
		{key: tcell.KeyTab, action: ActionSwitchPane},
		{rune: 'k', action: ActionScrollUp},
		{rune: 'j', action: ActionScrollDown},
		{rune: 'e', action: ActionExternalEdit},
		{rune: 'n', action: ActionNewNote},
		{rune: 'N', action: ActionNewFolder},
		{rune: 'd', action: ActionDelete},
		{key: tcell.KeyCtrlM, action: ActionRename},
		{key: tcell.KeyCtrlG, action: ActionRefresh},
		{key: tcell.KeyCtrlR, action: ActionToggleRaw},
		{key: tcell.KeyCtrlE, action: ActionEnterEdit},
	},
	ContextEditor: {
		{key: tcell.KeyLeft, action: ActionMoveUp},
		{key: tcell.KeyRight, action: ActionMoveDown},
		{key: tcell.KeyUp, action: ActionMoveUp},
		{key: tcell.KeyDown, action: ActionMoveDown},
		{key: tcell.KeyHome, action: ActionMoveUp},
		{key: tcell.KeyEnd, action: ActionMoveDown},
		{key: tcell.KeyCtrlS, action: ActionSave},
		{key: tcell.KeyCtrlQ, action: ActionEditorQuit},
		{key: tcell.KeyEscape, action: ActionCancel},
		{key: tcell.KeyTab, action: ActionTabIndent},
		{rune: '?', action: ActionHelp},
	},
}

// MatchKey returns the KeyAction for the given context and key event.
func MatchKey(ctx KeyContext, event *tcell.EventKey) (KeyAction, bool) {
	bindings, exists := Keymaps[ctx]
	if !exists {
		return "", false
	}

	key := event.Key()
	var rune rune
	if key == tcell.KeyRune {
		rune = event.Rune()
	}

	for _, b := range bindings {
		if b.rune != 0 {
			if key == tcell.KeyRune && b.rune == rune {
				return b.action, true
			}
			continue
		}
		if b.key == key {
			return b.action, true
		}
	}

	return "", false
}
