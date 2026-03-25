package tui

import (
	"github.com/gdamore/tcell/v2"
)

type KeyContext string

const (
	ContextGlobal  KeyContext = "global"
	ContextTree    KeyContext = "tree"
	ContextPreview KeyContext = "preview"
	ContextEditor  KeyContext = "editor"
)

type KeyAction string

const (
	ActionSwitchPane   KeyAction = "switch_pane"
	ActionQuit         KeyAction = "quit"
	ActionHelp         KeyAction = "help"
	ActionCancel       KeyAction = "cancel"
	ActionScrollUp     KeyAction = "scroll_up"
	ActionScrollDown   KeyAction = "scroll_down"
	ActionExpandOpen   KeyAction = "expand_open"
	ActionMoveUp       KeyAction = "move_up"
	ActionMoveDown     KeyAction = "move_down"
	ActionExternalEdit KeyAction = "external_edit"
	ActionNewNote      KeyAction = "new_note"
	ActionNewFolder    KeyAction = "new_folder"
	ActionDelete       KeyAction = "delete"
	ActionRename       KeyAction = "rename"
	ActionRefresh      KeyAction = "refresh"
	ActionToggleRaw    KeyAction = "toggle_raw"
	ActionEnterEdit    KeyAction = "enter_edit"
	ActionSave         KeyAction = "save"
	ActionEditorQuit   KeyAction = "editor_quit"
	ActionTabIndent    KeyAction = "tab_indent"
)

type keyBinding struct {
	key    tcell.Key
	rune   rune
	action KeyAction
}

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
