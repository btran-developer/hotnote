package tui

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TreeView struct {
	*tview.TreeView
	rootPath     string
	onFileSelect func(string)
}

func NewTreeView(root string, onFileSelect func(string)) *TreeView {
	tv := &TreeView{
		TreeView:     tview.NewTreeView(),
		rootPath:     root,
		onFileSelect: onFileSelect,
	}

	tv.SetRoot(tv.buildTree(root))
	tv.SetSelectedFunc(tv.handleSelect)
	tv.SetInputCapture(tv.handleInputCapture)

	tv.SetBorder(true).
		SetTitle(" " + filepath.Base(root) + " ").
		SetBackgroundColor(tcell.ColorDefault)

	return tv
}

func (tv *TreeView) buildTree(root string) *tview.TreeNode {
	rootNode := tview.NewTreeNode(TreePrefixFolderOpen + " " + filepath.Base(root)).
		SetReference(root)

	if err := tv.addChildren(rootNode, root); err != nil {
		return rootNode
	}

	return rootNode
}

func (tv *TreeView) addChildren(parent *tview.TreeNode, dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		if len(name) == 0 || name[0] == '.' {
			continue
		}

		entryPath := filepath.Join(dirPath, name)

		if entry.IsDir() {
			folderNode := tview.NewTreeNode(TreePrefixFolderClosed + " " + name).
				SetReference(entryPath).
				SetSelectable(true)
			if err := tv.addChildren(folderNode, entryPath); err != nil {
				continue
			}
			parent.AddChild(folderNode)
		} else if filepath.Ext(name) == ".md" {
			fileNode := tview.NewTreeNode(TreePrefixFile + " " + name[:len(name)-3]).
				SetReference(entryPath).
				SetSelectable(true)
			parent.AddChild(fileNode)
		}
	}

	return nil
}

func (tv *TreeView) handleSelect(node *tview.TreeNode) {
	ref := node.GetReference()
	if ref == nil {
		return
	}

	path, ok := ref.(string)
	if !ok {
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.IsDir() {
		if node.IsExpanded() {
			node.Collapse()
			offset := len(TreePrefixFolderClosed) + 1
			node.SetText(TreePrefixFolderClosed + " " + node.GetText()[offset:])
		} else {
			node.Expand()
			offset := len(TreePrefixFolderOpen) + 1
			node.SetText(TreePrefixFolderOpen + " " + node.GetText()[offset:])
		}
	} else {
		tv.onFileSelect(path)
	}
}

func (tv *TreeView) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	action, matched := MatchKey(ContextTree, event)
	if !matched {
		return event
	}

	switch action {
	case ActionSwitchPane:
		return event
	case ActionExpandOpen:
		return event
	}

	return event
}
