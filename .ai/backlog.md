# Backlog — HotNote

## CLI (Phase 1)
- [x] All items completed

## TUI (Phase 2)

### Issue 1: TUI Project Setup
- [x] Add dependencies (tview, tcell, goldmark, chroma)
- [x] Create entry point with tview Application
- [x] Set up basic app struct

### Issue 2: Workspace Selection
- [ ] Add "tui" command to CLI
- [ ] Create workspace selection overlay
- [ ] Handle no-workspace case (auto-create default)
- [ ] Transition to main view

### Issue 3: 2-Pane Layout + TreeView
- [ ] Create Flex layout (25% tree, 75% preview)
- [ ] Implement TreeView with folders + files
- [ ] Load directory structure from workspace
- [ ] Expand/collapse functionality

### Issue 4: Preview Pane
- [ ] Create TextView component
- [ ] Display markdown content on file selection
- [ ] Scroll support for long notes

### Issue 5: Raw/Rendered Toggle
- [ ] Integrate goldmark for markdown rendering
- [ ] Implement Ctrl+R toggle
- [ ] Apply markdown styling (headings, lists, etc.)

### Issue 6: Status Bar + Help Overlay
- [ ] Context-sensitive key hints in status bar
- [ ] Help overlay on ? key
- [ ] Scrollable help (k/j)
- [ ] Escape as universal cancel

### Issue 7: Editor Mode
- [ ] TextArea component for editing
- [ ] Cursor movement (char + line)
- [ ] Text insertion/deletion
- [ ] Save (Ctrl+S), discard (Escape), save+quit (Ctrl+Q)

### Issue 8: External Editor
- [ ] 'e' key to open in $EDITOR
- [ ] TUI fallback if $EDITOR not set
- [ ] Reload note after return

### Issue 9: Syntax Highlighting
- [ ] Integrate chroma library
- [ ] Apply to code blocks in preview
- [ ] Use monokai theme

### Issue 10: Create + Delete Note
- [ ] 'n' key to create new note
- [ ] 'd' key to delete with confirmation (y/n/c)
- [ ] Handle file operations

### Issue 11: Create Folder + Rename
- [ ] Shift+N to create new folder
- [ ] Prompt for folder name
- [ ] Ctrl+M to rename note
- [ ] Handle slug collision

### Issue 12: Refresh + Polish
- [ ] Ctrl+G to manual refresh
- [ ] Edge case handling
- [ ] Final testing

## Future
- Tabs (multiple notes open simultaneously)
- Dual-view editor (edit + preview side-by-side or toggle)
- Syntax highlighting in editor
- Tag filtering UI
- Full-text search UI
- Backlinks visualization
- Settings/configurability
- Multiple workspace management UI
- Customizable syntax highlighting theme
- Vim/Emacs keybindings
- Auto-save
- Multiple cursor editing
- Customizable pane layout
- Note templates
- Export (HTML, PDF)
- Import from other formats