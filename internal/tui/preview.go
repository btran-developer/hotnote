package tui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PreviewPane struct {
	*tview.TextView
	currentFile string
}

func NewPreviewPane() *PreviewPane {
	p := &PreviewPane{
		TextView:    tview.NewTextView(),
		currentFile: "",
	}

	p.SetDynamicColors(true).
		SetWrap(true).
		SetBorder(true).
		SetTitle(" Preview ").
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(p.handleInputCapture)

	return p
}

func (p *PreviewPane) LoadNote(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	p.currentFile = path
	p.SetTextColor(tcell.ColorDefault)
	p.SetText(string(content))
	return nil
}

func (p *PreviewPane) Clear() {
	p.currentFile = ""
	p.SetText("")
}

func (p *PreviewPane) GetCurrentFile() string {
	return p.currentFile
}

func (p *PreviewPane) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	action, matched := MatchKey(ContextPreview, event)
	if !matched {
		return event
	}

	switch action {
	case ActionSwitchPane:
		return event
	case ActionToggleRaw:
		return event
	case ActionEnterEdit:
		return event
	case ActionExternalEdit:
		return event
	case ActionNewNote:
		return event
	case ActionNewFolder:
		return event
	case ActionDelete:
		return event
	case ActionRename:
		return event
	case ActionRefresh:
		return event
	}

	return event
}
