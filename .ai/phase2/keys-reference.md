# Keys Reference — HotNote TUI

## Global Keys (Always Active)
| Key | Action | Disabled in Edit? |
|-----|--------|-------------------|
| ? | Help (scrollable overlay) | No |
| Escape | Universal cancel | No |
| Ctrl+Q | Quit app | Yes |
| Tab | Switch focus (pane cycling) | No |

## Tree (Left Pane)
| Key | Action | Disabled in Edit? |
|-----|--------|-------------------|
| k / ↑ | Move up | Yes |
| j / ↓ | Move down | Yes |
| Enter | Expand (folder) / Open (file) | Yes |
| Tab | Switch to preview pane | Yes |
| ? | Help | No |

## Preview Pane (Read Mode)
| Key | Action | Disabled in Edit? |
|-----|--------|-------------------|
| k / ↑ | Scroll up | N/A (not edit mode) |
| j / ↓ | Scroll down | N/A |
| Tab | Switch to tree pane | N/A |
| e | Open in $EDITOR | N/A |
| n | New note | N/A |
| Shift+N | New folder | N/A |
| d | Delete note | N/A |
| Ctrl+M | Rename note | N/A |
| Ctrl+G | Refresh | N/A |
| Ctrl+R | Toggle raw/rendered | N/A |
| Ctrl+E | Enter edit mode | N/A |
| ? | Help | No |

## Preview Pane (Edit Mode)
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

Note: In edit mode, all tree navigation (h/j/k/l), note actions (e, n, Shift+N, d, Ctrl+M), and view toggles (Ctrl+R, Ctrl+E, Ctrl+G) are disabled.