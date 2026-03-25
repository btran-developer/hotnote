package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestMatchKey_GlobalContext(t *testing.T) {
	tests := []struct {
		name     string
		key      tcell.Key
		rune     rune
		expected KeyAction
		found    bool
	}{
		{
			name:     "CtrlQ triggers quit",
			key:      tcell.KeyCtrlQ,
			expected: ActionQuit,
			found:    true,
		},
		{
			name:     "Tab triggers switch pane",
			key:      tcell.KeyTab,
			expected: ActionSwitchPane,
			found:    true,
		},
		{
			name:     "Escape triggers cancel",
			key:      tcell.KeyEscape,
			expected: ActionCancel,
			found:    true,
		},
		{
			name:     "rune ? triggers help",
			key:      tcell.KeyRune,
			rune:     '?',
			expected: ActionHelp,
			found:    true,
		},
		{
			name:  "unmatched key returns false",
			key:   tcell.KeyCtrlA,
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var event *tcell.EventKey
			if tt.rune != 0 {
				event = tcell.NewEventKey(tcell.KeyRune, tt.rune, 0)
			} else {
				event = tcell.NewEventKey(tt.key, 0, 0)
			}

			action, found := MatchKey(ContextGlobal, event)
			assert.Equal(t, tt.found, found)
			if tt.found {
				assert.Equal(t, tt.expected, action)
			}
		})
	}
}

func TestMatchKey_TreeContext(t *testing.T) {
	tests := []struct {
		name     string
		key      tcell.Key
		rune     rune
		expected KeyAction
		found    bool
	}{
		{
			name:     "Up arrow triggers move up",
			key:      tcell.KeyUp,
			expected: ActionMoveUp,
			found:    true,
		},
		{
			name:     "k rune triggers move up",
			key:      tcell.KeyRune,
			rune:     'k',
			expected: ActionMoveUp,
			found:    true,
		},
		{
			name:     "j rune triggers move down",
			key:      tcell.KeyRune,
			rune:     'j',
			expected: ActionMoveDown,
			found:    true,
		},
		{
			name:     "Enter triggers expand open",
			key:      tcell.KeyEnter,
			expected: ActionExpandOpen,
			found:    true,
		},
		{
			name:     "Tab triggers switch pane",
			key:      tcell.KeyTab,
			expected: ActionSwitchPane,
			found:    true,
		},
		{
			name:  "unmatched key returns false",
			key:   tcell.KeyCtrlZ,
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var event *tcell.EventKey
			if tt.rune != 0 {
				event = tcell.NewEventKey(tcell.KeyRune, tt.rune, 0)
			} else {
				event = tcell.NewEventKey(tt.key, 0, 0)
			}

			action, found := MatchKey(ContextTree, event)
			assert.Equal(t, tt.found, found)
			if tt.found {
				assert.Equal(t, tt.expected, action)
			}
		})
	}
}

func TestMatchKey_PreviewContext(t *testing.T) {
	tests := []struct {
		name     string
		key      tcell.Key
		rune     rune
		expected KeyAction
		found    bool
	}{
		{
			name:     "e rune triggers external edit",
			key:      tcell.KeyRune,
			rune:     'e',
			expected: ActionExternalEdit,
			found:    true,
		},
		{
			name:     "n rune triggers new note",
			key:      tcell.KeyRune,
			rune:     'n',
			expected: ActionNewNote,
			found:    true,
		},
		{
			name:     "d rune triggers delete",
			key:      tcell.KeyRune,
			rune:     'd',
			expected: ActionDelete,
			found:    true,
		},
		{
			name:     "CtrlR triggers toggle raw",
			key:      tcell.KeyCtrlR,
			expected: ActionToggleRaw,
			found:    true,
		},
		{
			name:     "CtrlE triggers enter edit",
			key:      tcell.KeyCtrlE,
			expected: ActionEnterEdit,
			found:    true,
		},
		{
			name:  "unmatched key returns false",
			key:   tcell.KeyF1,
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var event *tcell.EventKey
			if tt.rune != 0 {
				event = tcell.NewEventKey(tcell.KeyRune, tt.rune, 0)
			} else {
				event = tcell.NewEventKey(tt.key, 0, 0)
			}

			action, found := MatchKey(ContextPreview, event)
			assert.Equal(t, tt.found, found)
			if tt.found {
				assert.Equal(t, tt.expected, action)
			}
		})
	}
}

func TestMatchKey_EditorContext(t *testing.T) {
	tests := []struct {
		name     string
		key      tcell.Key
		rune     rune
		expected KeyAction
		found    bool
	}{
		{
			name:     "CtrlS triggers save",
			key:      tcell.KeyCtrlS,
			expected: ActionSave,
			found:    true,
		},
		{
			name:     "Escape triggers cancel",
			key:      tcell.KeyEscape,
			expected: ActionCancel,
			found:    true,
		},
		{
			name:     "Tab triggers tab indent",
			key:      tcell.KeyTab,
			expected: ActionTabIndent,
			found:    true,
		},
		{
			name:     "Left arrow triggers move up",
			key:      tcell.KeyLeft,
			expected: ActionMoveUp,
			found:    true,
		},
		{
			name:  "unmatched key returns false",
			key:   tcell.KeyInsert,
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var event *tcell.EventKey
			if tt.rune != 0 {
				event = tcell.NewEventKey(tcell.KeyRune, tt.rune, 0)
			} else {
				event = tcell.NewEventKey(tt.key, 0, 0)
			}

			action, found := MatchKey(ContextEditor, event)
			assert.Equal(t, tt.found, found)
			if tt.found {
				assert.Equal(t, tt.expected, action)
			}
		})
	}
}

func TestMatchKey_RuneVsKeyDistinction(t *testing.T) {
	event := tcell.NewEventKey(tcell.KeyRune, 'j', 0)
	action, found := MatchKey(ContextTree, event)
	assert.True(t, found)
	assert.Equal(t, ActionMoveDown, action)
}
