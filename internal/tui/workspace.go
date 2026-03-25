package tui

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type WorkspaceSelector struct {
	*tview.List
	workspaces []workspaceEntry
	current    string
	onSelect   func(string, string)
	onCancel   func()
}

type workspaceEntry struct {
	name    string
	path    string
	modTime time.Time
}

func sortWorkspaces(workspaces map[string]string, current string) []workspaceEntry {
	entries := make([]workspaceEntry, 0, len(workspaces))

	for name, path := range workspaces {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		entries = append(entries, workspaceEntry{
			name:    name,
			path:    path,
			modTime: info.ModTime(),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].name == current {
			return true
		}
		if entries[j].name == current {
			return false
		}
		if entries[i].name == "default" {
			return true
		}
		if entries[j].name == "default" {
			return false
		}
		return entries[i].modTime.After(entries[j].modTime)
	})

	return entries
}

func NewWorkspaceSelector(
	entries []workspaceEntry,
	current string,
	onSelect func(string, string),
	onCancel func(),
) *WorkspaceSelector {
	w := &WorkspaceSelector{
		List:       tview.NewList(),
		workspaces: entries,
		current:    current,
		onSelect:   onSelect,
		onCancel:   onCancel,
	}

	w.buildList()
	w.SetSelectedFunc(w.handleSelect)
	w.SetDoneFunc(w.handleCancel)
	w.SetInputCapture(w.handleInputCapture)

	w.SetBorder(true).
		SetTitle(" Select Workspace ").
		SetBackgroundColor(tcell.ColorDefault)

	return w
}

func (w *WorkspaceSelector) buildList() {
	w.Clear()

	first := true
	for _, entry := range w.workspaces {
		prefix := "  "
		if entry.name == w.current {
			prefix = " >"
		} else if first && w.current == "" {
			prefix = " >"
		}

		display := fmt.Sprintf("%s %s", prefix, entry.name)
		w.AddItem(display, entry.path, 0, nil)
		first = false
	}
}

func (w *WorkspaceSelector) handleSelect(index int, mainText, secondaryText string, rune rune) {
	entry := w.workspaces[index]
	w.onSelect(entry.name, entry.path)
}

func (w *WorkspaceSelector) handleCancel() {
	w.onCancel()
}

func (w *WorkspaceSelector) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp, tcell.KeyCtrlP:
		// Let List handle its own up/down
		return event
	case tcell.KeyRune:
		switch event.Rune() {
		case 'k':
			current := w.GetCurrentItem()
			if current > 0 {
				w.SetCurrentItem(current - 1)
			}
			return nil
		case 'j':
			current := w.GetCurrentItem()
			w.SetCurrentItem(current + 1)
			return nil
		}
	}
	return event
}

func (w *WorkspaceSelector) HandleKey(event *tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape {
		w.onCancel()
		return true
	}
	if event.Key() == tcell.KeyEnter {
		idx := w.GetCurrentItem()
		w.handleSelect(idx, "", "", 0)
		return true
	}
	return false
}
