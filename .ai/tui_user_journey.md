# User Journey — HotNote TUI

## Journey 0: Launch + Workspace Selection
1. TUI launches
2. If workspaces exist → show selection overlay
3. If no workspaces → auto-create default workspace
4. User selects workspace (Enter) or creates new
5. Transition to main 2-pane view

## Journey 1: Navigate Tree
1. Focus in tree (left pane)
2. k/j = move up/down
3. Enter = expand/collapse (folder) or open (file)
4. l = move to preview pane

## Journey 2: Preview Note
1. View rendered markdown (default)
2. k/j = scroll up/down
3. Ctrl+R = toggle raw/rendered
4. h = exit to tree

## Journey 3: Edit Note
1. Ctrl+E = enter edit mode
2. Arrow keys = move cursor (←→ char, ↑↓ line)
3. Type = insert text
4. Backspace = delete
5. Tab = 4 spaces
6. Home/End = line start/end
7. Ctrl+S = save and stay
8. Ctrl+Q = save and quit to preview
9. Escape = discard and exit to preview
10. e = open in $EDITOR (or TUI fallback)

## Journey 4: Note/Folder Management
- **Create note**: n = new note → $EDITOR (or TUI fallback)
- **Create folder**: Shift+N = new folder → prompt for name
- **Delete**: d = delete → confirm prompt → y/n/c
- **Rename**: Ctrl+M = rename → prompt for new title

## Journey 5: Refresh
- Ctrl+G = refresh tree and current note

## Journey 6: Help
- ? = show help overlay
- k/j = scroll up/down
- Escape = close help

## Journey 8: Quit
- q (not in edit mode) = exit to shell
- q (in edit mode with unsaved) = prompt "Save changes? (y/n/c)"