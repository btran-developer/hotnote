# HotNote TUI Specification (Phase 2)

## 1. Overview

Framework: tview (rivo/tview)
Terminal library: tcell (gdamore/tcell)

---

## 2. Launch

### Usage
```
hotnote
hotnote tui
```

### Workflow
1. TUI launches
2. If workspaces exist → show workspace selection overlay
3. If no workspaces → auto-create default workspace silently
4. User selects workspace (Enter) or creates new
5. Transition to main 2-pane view

### Workspace Selection Overlay
```
┌─────────────────────────────────────────────┐
│          Select Workspace                   │
├─────────────────────────────────────────────┤
│  > default                    ~/.local/... │
│    work                      ~/notes/work  │
│    [+ create new workspace]                │
│                                             │
│  [Enter] select  [Esc] cancel              │
└─────────────────────────────────────────────┘
```
- Arrow keys (k/j) to navigate
- Enter to select
- Escape to cancel (closes TUI)
- Show current active workspace marked with `>`

---

## 3. Layout

### Dimensions
- Full terminal width/height
- Tree: 25% width (min 15 chars)
- Preview/Editor: 75% width (remaining)

### Regions
```
+-------------------+--------------------------------+
| TREE              | PREVIEW/EDITOR                 |
| (left pane)       | (right pane)                  |
|                   |                                |
| [-] workspace/   | # Note content                 |
|   [+] project-a/ |                                |
|     * note.md    |  Rendered markdown...         |
|     * design.md  |                                |
|   [+] journal/   |                                |
+-------------------+--------------------------------+
| Status bar with context-sensitive key hints  |
+-------------------------------------------------+
```

### Visual Styles (tview + tcell)
- Tree: Dim color for folders, normal for files
- `[+]` = collapsed folder
- `[-]` = expanded folder
- `*` = file prefix
- Preview: Normal text, monospace in raw mode
- Editor: Editable background
- Selected: Highlighted background

### Status Bar
- Bottom of screen
- Context-sensitive hints based on current pane/mode
- Preview read mode: `h:← | l:→ | n:new | +N:folder | d:del | R:ren | e:ext | Ctrl+E:edit | Ctrl+R:raw | ?:help`
- Preview edit mode: `←→:char | ↑↓:line | Home/End | Ctrl+S:save | Ctrl+Q:quit | Esc:cancel | ?:help`

---

## 4. Tree (Left Pane)

### Component
- tview.TreeView

### Display
- Show folder hierarchy + files from workspace root
- Folders: `[+]` collapsed, `[-]` expanded
- Files: `*` prefix
- Indent with `  ` per level

### Selection
- Single node selected at a time
- Highlighted when selected

### Navigation
| Key | Action |
|-----|--------|
| k / ↑ | Move up |
| j / ↓ | Move down |
| Enter | Expand (folder) or open (file) |
| l | Move to preview pane |
| h | Exit pane (or stay at edge) |
| ? | Help |

### Virtual Scrolling
- tview handles automatically

---

## 5. Preview Pane (Right Pane)

### Component
- tview.TextView (read mode)
- tview.TextArea (edit mode)

### Display (Rendered Mode)
- Render markdown to styled text using goldmark
- Full markdown support:
  - H1-H6: Size variation + bold
  - Bold: **text** → bold style
  - Italic: *text* → italic style
  - Code blocks: Background color + syntax highlighting (chroma, monokai theme)
  - Inline code: Background color
  - Lists: Bullet/number with indent
  - Links: Text + (url) in dim color
  - Blockquotes: Indented + border left
  - Tables: Aligned columns
  - Horizontal rule: Dim line
- Vertical scroll for long notes

### Display (Raw Mode)
- Show raw markdown text
- No rendering
- Monospace font

### Toggle
- Ctrl+R: Toggle raw/rendered
- Default: Rendered mode

---

## 6. Editor Mode (Part 2)

### Activation
- Ctrl+E: Toggle from preview to edit mode

### Component
- tview.TextArea

### Edit Mode Behavior
- Replace preview with editable text area
- Show cursor
- Enable text input
- Load note content for editing

### Editing Keys
| Key | Action |
|-----|--------|
| ← / → | Move cursor (character) |
| ↑ / ↓ | Move cursor (line) |
| Home | Move to line start |
| End | Move to line end |
| Type | Insert character |
| Backspace | Delete character |
| Enter | New line |
| Tab | Insert 4 spaces |
| Ctrl+S | Save and stay |
| Ctrl+Q | Save and exit edit mode |
| Escape | Exit without saving |

### Save Behavior
- Write to file immediately
- Update updated_at in frontmatter
- Show "Saved" indicator briefly
- On error: show error, stay in edit mode

### Exit Edit Mode
- Escape: Discard changes, return to preview
- Ctrl+Q: Save changes, return to preview

---

## 7. External Editor

### Activation
- Press `e` while note selected

### Behavior
- If `$EDITOR` set: spawn $EDITOR process, wait for close
- If `$EDITOR` not set: open in TUI edit mode (fallback)
- After return: reload note in preview

### Fallback
- If $EDITOR fails: show error, stay in TUI

---

## 8. Note/Folder Management (Part 3)

### Create Note
- Key: `n`
- Behavior: Create in current tree selection (folder or parent of file)
- If nothing selected: create at root
- Opens in $EDITOR, if not set → TUI edit mode

### Create Folder
- Key: `Shift+N`
- Behavior: Show prompt for folder name
- Create in current tree selection (or root if nothing selected)
- Show in tree immediately

### Delete Note
- Key: `d`
- Behavior: Show confirmation prompt "Delete <slug>? (y/n/c)"
- y: Delete file, refresh tree, clear preview
- n: Cancel, return to preview
- c: Cancel (same as n)

### Rename Note
- Key: `Ctrl+M`
- Behavior: Show prompt for new title
- Generate new slug from title
- Update filename and frontmatter title
- Handle slug collision (append -1, -2, etc.)

### Refresh
- Key: Ctrl+G
- Behavior: Reload tree and current note content
- Useful after external changes

---

## 9. Help Overlay

### Activation
- Key: `?`

### Display
- Full-screen overlay with scrollable list
- Show all keyboard shortcuts grouped by context

### Navigation
- k / ↑: Scroll up
- j / ↓: Scroll down
- Escape: Close help

---

## 10. Keyboard Shortcuts Summary

### Global Keys (Always Active)
| Key | Action | Disabled in Edit? |
|-----|--------|-------------------|
| ? | Help (scrollable overlay) | No |
| Escape | Universal cancel | No |
| q | Quit app | Yes |

### Tree (Left Pane)
| Key | Action | Disabled in Edit? |
|-----|--------|-------------------|
| k / ↑ | Move up | Yes |
| j / ↓ | Move down | Yes |
| Enter | Expand (folder) / Open (file) | Yes |
| l | Move to preview | Yes |
| ? | Help | No |

### Preview Pane (Read Mode)
| Key | Action | Disabled in Edit? |
|-----|--------|-------------------|
| k / ↑ | Scroll up | N/A (not edit mode) |
| j / ↓ | Scroll down | N/A |
| h | Move to tree | N/A |
| e | Open in $EDITOR | N/A |
| n | New note | N/A |
| Shift+N | New folder | N/A |
| d | Delete note | N/A |
| Ctrl+M | Rename note | N/A |
| Ctrl+G | Refresh | N/A |
| Ctrl+R | Toggle raw/rendered | N/A |
| Ctrl+E | Enter edit mode | N/A |
| ? | Help | No |

### Preview Pane (Edit Mode)
| Key | Action | Disabled in Edit? |
|-----|--------|-------------------|
| ← / → | Move cursor (char) | N/A (edit mode) |
| ↑ / ↓ | Move cursor (line) | N/A |
| Home | Move to line start | N/A |
| End | Move to line end | N/A |
| Type | Insert character | N/A |
| Backspace | Delete character | N/A |
| Enter | New line | N/A |
| Tab | Insert 4 spaces | N/A |
| Ctrl+S | Save and stay | N/A |
| Ctrl+Q | Save and quit to preview | N/A |
| Escape | Discard and exit to preview | No |
| ? | Help | No |

---

## 11. State Management

```go
type App struct {
    // Workspace
    workspaceRoot string
    workspaces    []Workspace
    
    // UI Components (tview)
    tree     *tview.TreeView
    preview  *tview.TextView
    editor   *tview.TextArea
    status   *tview.TextView
    
    // State
    selectedNode   string  // current tree node path
    currentNote    string  // current note path
    content        string  // current note content
    showRaw        bool
    editMode       bool
    
    // UI state
    helpVisible    bool
    promptActive   bool
    promptType     string
}
```

---

## 12. Error Handling

| Scenario | Behavior |
|----------|----------|
| No workspace | Auto-create default silently |
| Empty workspace | Show "No notes yet" |
| Note deleted externally | Show error, refresh tree |
| Save fails | Show error, stay in edit mode |
| $EDITOR fails | Show error, stay in TUI |
| Corrupted markdown | Show raw content + error indicator |

---

## 13. Exit

### Quit Keys
- q: Quit if not in edit mode
- Ctrl+C: Always quit (force)

### Unsaved Changes
- If edit mode with unsaved changes:
  - q: Prompt "Save changes? (y/n/c)" 
    - y: Save and quit
    - n: Discard and quit
    - c: Cancel (stay in edit mode)

### Return
- Exit to shell, no output

---

## 14. Dependencies

### Required Packages
- github.com/rivo/tview (Framework + components)
- github.com/gdamore/tcell/v2 (Terminal handling)
- github.com/yuin/goldmark (Markdown rendering)
- github.com/alecthomas/chroma/v2 (Syntax highlighting)

### Version Constraints
- Go 1.22+
- tview latest
- tcell v2
- Goldmark v1.x
- Chroma v2.x

---

## 15. Implementation Order

### Part 1: Navigation + Preview
1. Setup tview project structure
2. Workspace selection overlay on launch
3. TreeView for tree (folders + files)
4. TextView for preview
5. Raw/rendered toggle (Ctrl+R)
6. Status bar with key hints
7. Help overlay (?) with scroll

### Part 2: Editor + Syntax Highlighting
8. TextArea for editor
9. Cursor movement (char + line)
10. Text insertion/deletion
11. Save (Ctrl+S), discard (Escape), save+quit (Ctrl+Q)
12. External editor (e) with TUI fallback
13. Syntax highlighting in preview (chroma, monokai)
14. Unsaved changes prompt

### Part 3: Note/Folder Management
15. Create note (n)
16. Create folder (Shift+N)
17. Delete note (d) with confirmation
18. Rename note (Ctrl+M)
19. Manual refresh (Ctrl+G)
20. Final polish and edge cases

---

## 16. Testing

### Manual Verification
1. Launch TUI → workspace selection shown
2. Select workspace → 2 panes visible
3. Navigate tree → expand/collapse folders
4. Select file → preview shows content
5. Press Ctrl+R → toggle raw/rendered
6. Press Ctrl+E → enter edit mode
7. Type text → text appears
8. Press Ctrl+S → save (file updated)
9. Press e → $EDITOR opens (or TUI fallback)
10. Press n → create new note
11. Press Shift+N → create new folder
12. Press d → delete prompt
13. Press Ctrl+M → rename prompt
14. Press Ctrl+G → refresh
15. Press ? → help overlay, scroll, Escape to close
16. Press q → exit to shell

### Edge Cases
- Empty workspace
- Deep folder nesting
- Very long note
- Special characters in title
- Multiple markdown features
- No $EDITOR set (TUI fallback)
- Unsaved changes on quit