package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StatusBar displays status information at the bottom of the TUI.
type StatusBar struct {
	*tview.TextView
}

// NewStatusBar creates a new status bar.
func NewStatusBar() *StatusBar {
	sb := &StatusBar{
		TextView: tview.NewTextView(),
	}

	sb.SetDynamicColors(true).
		SetWrap(false).
		SetTextColor(DefaultPalette.Info).
		SetBackgroundColor(tcell.ColorDefault).
		SetBorder(false)

	return sb
}

// SetStatus sets the status bar text.
func (sb *StatusBar) SetStatus(text string) {
	sb.TextView.SetText(text)
}
