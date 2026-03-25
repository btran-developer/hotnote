package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusBar struct {
	*tview.TextView
}

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

func (sb *StatusBar) SetStatus(text string) {
	sb.TextView.SetText(text)
}
