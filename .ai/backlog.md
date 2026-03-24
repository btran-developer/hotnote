# Backlog — HotNote

## CLI (Phase 1)
- [x] All items completed

## TUI (Phase 2)

### Part 1: Navigation + Preview
- [ ] Workspace selection on launch (overlay with list)
- [ ] 2-pane layout (tree 25%, preview 75%)
- [ ] TreeView (tview) showing folders + files
- [ ] Virtual scrolling (lazy load on scroll)
- [ ] Pane toggles (Ctrl+T tree, Ctrl+N list)
- [ ] Raw/rendered toggle (Ctrl+R)
- [ ] Help overlay (?) with k/j scroll, Escape close
- [ ] Status bar (bottom, context-sensitive hints)
- [ ] Escape = universal cancel

### Part 2: Editor + Syntax Highlighting
- [ ] In-TUI editing (Ctrl+E toggle)
- [ ] External editor (e) with TUI fallback if no $EDITOR
- [ ] Syntax highlighting in preview (monokai theme, chroma library)
- [ ] Save (Ctrl+S), discard (Escape), save+quit (Ctrl+Q)
- [ ] Unsaved changes prompt on quit

### Part 3: Note/Folder Management
- [ ] Create note (n) → $EDITOR or TUI fallback
- [ ] Create folder (Shift+N) → prompt for name
- [ ] Delete note (d) → confirmation prompt (y/n/c)
- [ ] Rename note (Ctrl+M) → prompt for new title
- [ ] Manual refresh (Ctrl+G) → reload all

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