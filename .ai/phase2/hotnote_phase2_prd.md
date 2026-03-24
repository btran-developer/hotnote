# HotNote — Phase 2 PRD (TUI + Editor)

## 1. Objective

Build a Terminal User Interface (TUI) for navigating, viewing, and editing markdown notes, providing a two-pane layout with full markdown rendering and syntax highlighting support.

---

## 2. Non-Goals (Phase 2 Scope Control)

Phase 2 will NOT include:

- AI integration
- Tag filtering UI
- Full-text search UI
- Backlinks/graph visualization
- Multiple cursor editing
- Syntax highlighting in editor (Phase 3+)
- Tabs (future)

---

## 3. Target Users

Primary:
- Developers who prefer terminal-based workflows
- Users who want visual navigation of their note collection
- Users who want to edit notes without leaving the terminal

---

## 4. Core Use Cases

Launch TUI:
```
hotnote
hotnote tui
```

Navigate tree (left pane):
- Arrow keys to move between nodes
- Enter to expand/collapse (folder) or open (file)
- Shows folders + files in single tree

Preview/Edit note (right pane):
- Renders full markdown (headers, lists, code blocks, links, bold/italic, blockquotes)
- Toggle between raw markdown and rendered view (Ctrl+R)
- Enter edit mode for inline editing (Ctrl+E)
- Syntax highlighting for code blocks in preview (editor in Phase 3+)
- Hybrid: press `e` to open in $EDITOR instead

Note/Folder management:
- `n` = new note
- `Shift+N` = new folder
- `d` = delete note
- `Ctrl+M` = rename note

---

## 5. Functional Requirements

### 5.1 Layout

Two-column layout:

```
+-------------------+--------------------------------+
| TREE              | PREVIEW/EDITOR                 |
| (25% width)       | (75% width)                   |
|                   |                                |
| [-] workspace/   | # My Idea                      |
|   [+] project-a/ |                                |
|     * note.md    | Some content...                |
|     * design.md  |                                |
|   [+] journal/   |                                |
+-------------------+--------------------------------+
```

- Tree shows both folders and files
- Resizable panes (future - key-based)
- Minimum widths: 15 chars each pane

### 5.2 Part 1: Navigation + Preview

#### 5.2.1 Tree (Left Pane)

- Displays folder hierarchy + files from workspace root
- Folders: `[+]` collapsed, `[-]` expanded
- Files: `*` prefix
- Keyboard navigation: k/j to move, Enter to select

#### 5.2.2 Preview Pane (Right Pane)

- Renders full markdown to styled text
- Supports:
  - Headings (H1-H6)
  - Bold, italic, strikethrough
  - Ordered and unordered lists
  - Links (display URL in parentheses)
  - Code blocks with syntax highlighting (preview only)
  - Inline code
  - Blockquotes
  - Horizontal rules
  - Tables (basic)
- Toggle raw/rendered view with Ctrl+R
- Scroll support for long notes

### 5.3 Part 2: Editor + Syntax Highlighting

#### 5.3.1 In-TUI Editor

- Toggle between preview and edit mode with Ctrl+E
- Basic text editing:
  - Arrow keys to navigate (char + line)
  - Home/End for line start/end
  - Type to insert text
  - Backspace to delete
  - Enter for new lines
- Save with Ctrl+S
- Exit edit mode with Escape (discard changes) or Ctrl+Q (save and quit)

#### 5.3.2 External Editor (Hybrid)

- Press `e` while note selected to open in $EDITOR
- Same behavior as CLI `open` command
- Falls back to TUI edit mode if $EDITOR not set

#### 5.3.3 Syntax Highlighting

Using chroma library (based on highlight.go):

- **Supported languages**: Go, JavaScript, TypeScript, Python, Rust, Ruby, Java, C, C++, HTML, CSS, JSON, YAML, Bash, SQL
- **Scope**: Code blocks in preview only (editor in Phase 3+)
- **Theme**: Monokai (hardcoded)

### 5.4 Note/Folder Management (Part 3)

#### 5.4.1 Create Note

- Key: `n`
- Create in current tree selection (folder or parent of file)
- If nothing selected: create at root
- Opens in $EDITOR, or TUI edit mode if not set

#### 5.4.2 Create Folder

- Key: `Shift+N`
- Prompt for folder name
- Create in current tree selection (or root if nothing selected)

#### 5.4.3 Delete Note

- Key: `d`
- Confirmation prompt: "Delete <slug>? (y/n/c)"

#### 5.4.4 Rename Note

- Key: `Ctrl+M`
- Prompt for new title
- Handle slug collision (append -1, -2, etc.)

#### 5.4.5 Refresh

- Key: `Ctrl+G`
- Reload tree and current note content

### 5.5 Keyboard Shortcuts

| Key | Action |
|-----|--------|
| k / ↑ | Move up (tree/scroll) |
| j / ↓ | Move down (tree/scroll) |
| h | Move to tree pane |
| l | Move to preview pane |
| Enter | Expand (folder) / Open (file) |
| Ctrl+R | Toggle raw/rendered markdown |
| Ctrl+E | Toggle edit mode |
| e | Open in $EDITOR (external) |
| n | New note |
| Shift+N | New folder |
| d | Delete note |
| Ctrl+M | Rename note |
| Ctrl+G | Refresh |
| Ctrl+S | Save (edit mode) |
| Escape | Cancel / Discard edits |
| Ctrl+Q | Save and quit edit mode |
| q / Ctrl+C | Quit TUI |
| ? | Show help |

### 5.6 Startup Behavior

- Fresh start: Workspace selection overlay → main 2-pane view
- No folder/file pre-selected
- User must manually navigate to view notes

---

## 6. Technical Requirements

### 6.1 Dependencies

- **tview** (rivo/tview) - Framework + components
- **tcell** (gdamore/tcell/v2) - Terminal handling
- **goldmark** (yuin/goldmark) - Markdown rendering
- **chroma** (alecthomas/chroma/v2) - Syntax highlighting

### 6.2 Architecture

```
internal/tui/
  ├── main.go        # Entry point
  ├── app.go         # tview Application setup
  ├── tree.go        # TreeView component
  ├── preview.go     # TextView for markdown preview
  ├── editor.go      # TextArea for editing
  ├── keys.go        # Key bindings
  └── workspace.go   # Workspace selection overlay
```

### 6.3 State Management

```go
type App struct {
    workspaceRoot string
    workspaces    []Workspace
    
    tree    *tview.TreeView
    preview *tview.TextView
    editor  *tview.TextArea
    status  *tview.TextView
    
    selectedNode string  // current tree node path
    currentNote  string  // current note path
    content      string // current note content
    showRaw      bool
    editMode     bool
}
```

---

## 7. Non-Functional Requirements

- Fast startup (<500ms)
- Responsive navigation (<50ms)
- Smooth scrolling
- Memory efficient (stream large notes)
- Works on terminals supporting 256 colors

---

## 8. Edge Cases

- Empty workspace: Show "No notes yet"
- Very long note titles: Truncate with `...`
- Deep folder nesting: Horizontal scroll or truncation
- Missing/corrupted markdown: Show raw content + error indicator
- No selection: Default to root for create
- Large notes: Scroll handled by tview
- Unsaved changes: Prompt before quit
- External editor fails: Show error, stay in TUI

---

## 9. Success Criteria

**Part 1**:
- TUI launches and displays 2-pane layout
- Tree navigation works with keyboard
- File selection updates preview pane
- Markdown renders correctly (full feature set)
- Raw/rendered toggle works

**Part 2**:
- In-TUI editing works (insert, delete, save)
- External editor opens with `e`
- Syntax highlighting appears in code blocks (preview)
- Unsaved changes prompt on quit

**Part 3**:
- Create note works
- Create folder works
- Delete note with confirmation works
- Rename note works
- Refresh works

---

## 10. Implementation Plan

### Part 1 (3 weeks)

**Week 1**:
- Project setup (tview + tcell)
- Workspace selection overlay
- TreeView component (folders + files)

**Week 2**:
- TextView for preview
- Integration with storage layer

**Week 3**:
- Raw/rendered toggle
- Help overlay
- Status bar

### Part 2 (2 weeks)

**Week 4**:
- TextArea for editor
- Basic text editing (insert, delete, navigation)

**Week 5**:
- Syntax highlighting integration (chroma, monokai)
- External editor integration
- Polish (save, quit, unsaved prompts)

### Part 3 (1-2 weeks)

**Week 6**:
- Create note (n)
- Create folder (Shift+N)
- Delete note (d) with confirmation
- Rename note (Ctrl+M)
- Manual refresh (Ctrl+G)

---

## 11. Principles

- Preview-first, editor second
- Hybrid editing (in-TUI or external)
- Keyboard-driven
- Markdown-native rendering
- Syntax highlighting in preview (Phase 2), editor (Phase 3+)
- Clean separation of components

---

## 12. Future (Phase 3+)

- Tabs (multiple notes open)
- Dual-view editor
- Syntax highlighting in editor
- Tag filtering
- Full-text search
- Backlinks visualization
- Multiple cursor editing
- Auto-save
- Vim/Emacs keybindings
- Customizable pane layout
- Settings/configurability