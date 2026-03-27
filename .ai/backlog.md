# Backlog — HotNote

## CLI (Phase 1) - Expanded

### Issue A: Subfolder Support - Storage Layer
- [x] Add `Store.List()` with `filepath.WalkDir`
- [x] Add `Store.Delete()` method
- [x] Add `Store.Rename()` method
- [x] Update `Store.Path()` for hybrid resolution (direct + recursive)

### Issue B: Subfolder Support - CLI Commands
- [x] Update `list` command to show subfolder notes
- [x] Fix `open` command slugify inconsistency
- [x] Fix `render` command slugify inconsistency
- [x] Fix `new` command to support `--path` flag

### Issue C: AI Agent Compatibility
- [ ] Use `ExitInvalidInput` for "title required", "note exists" errors
- [ ] Fix `ai` command JSON output
- [ ] Fix `workspace list` schema to match spec
- [ ] Make error messages deterministic (remove `%v` formatting)
- [ ] Add `--no-open` flag to `new` command
- [ ] Add `--slug` flag to `new` command

### Issue D: Documentation Fixes
- [ ] Remove or implement `--data-dir` flag
- [ ] Document `hotnote tui` command
- [ ] Fix `workspace new` help text (positional vs flag)
- [ ] Add error documentation for all commands
- [ ] Fix `workspace list` output examples
- [ ] Remove `$DEBUG` documentation if not implemented

### Issue E: mkdir Command
- [x] Create `cmd/mkdir.go`
- [x] `hotnote mkdir <folder>` creates folder
- [x] Support nested folder creation
- [x] JSON output support
- [x] Error handling

### Issue F: rmdir Command
- [x] Create `cmd/rmdir.go`
- [x] `hotnote rmdir <folder>` deletes folder
- [x] Confirmation prompt if folder not empty
- [x] `--force` flag to skip prompt
- [x] JSON output support

### Issue G: delete Command
- [x] Add `Delete(id string) error` to `storage.Store`
- [x] Create `cmd/delete.go`
- [x] `hotnote delete <slug>` deletes note
- [x] Hybrid path resolution support
- [x] Confirmation prompt by default
- [x] `--force` flag to skip prompt
- [x] JSON output support

### Issue H: workspace delete Command
- [x] Add `workspaceDeleteCmd` to `cmd/workspace.go`
- [x] `hotnote workspace delete <name>` deletes workspace
- [x] Add `Delete(name string) error` to `workspace.Manager`
- [x] Recursively delete workspace directory and all contents
- [x] Confirmation prompt by default
- [x] `--force` flag to skip prompt
- [x] Safety checks (current, default workspace)

### Issue I: Unit Tests
- [x] Test hybrid path resolution
- [x] Test `Store.Delete()`, `Store.Rename()`, `Store.List()`
- [x] Test `workspace.Manager.Delete()`
- [x] Test all new commands
- [ ] Test safety checks
- [ ] Test conflict resolution

### Issue J: Documentation Updates
- [ ] Update `docs/features/commands.md`
- [ ] Update CLI examples
- [ ] Add error scenarios
- [ ] Document subfolder support
- [ ] Document AI agent compatibility

### Issue K: rename Command
- [x] Create `cmd/rename.go`
- [x] `hotnote rename <old> <new>` renames note
- [x] Hybrid path resolution for old-slug
- [x] Uses existing `Store.Rename()` method
- [x] Validates new-slug doesn't conflict with existing note
- [x] Export `slugify()` or move to internal package
- [x] JSON output support

### Issue L: folder rename Command
- [x] Create `cmd/folder.go` with `folder` parent command
- [x] `hotnote folder rename <old> <new>` renames folder
- [x] Validate both paths with `pathutil.ValidateFolderPath`
- [x] Reject if destination folder exists
- [x] Reject renaming workspace root
- [x] JSON output support

### Issue M: workspace rename Command
- [ ] Add `workspaceRenameCmd` to `cmd/workspace.go`
- [ ] `hotnote workspace rename <old> <new>` renames workspace
- [ ] Add `Rename(oldName, newName string) error` to `workspace.Manager`
- [ ] Update config.workspaces map (move value, delete old key)
- [ ] Update current_workspace if renaming current
- [ ] Safety: reject if new-name already exists
- [ ] JSON output support

### Issue N: Add CrTime to NoteInfo
- [ ] Add CrTime field to storage.NoteInfo struct
- [ ] Populate from frontmatter created_at when available
- [ ] Fall back to filesystem birth time or ModTime
- [ ] Display creation time in list output
- [ ] Document behavior in commands.md

### Issue O: CLI Command Restructure - Breaking Changes
- [ ] Refactor `mkdir` → `folder make` with `mk` alias
- [ ] Refactor `rmdir` → `folder remove` with `rmv` alias
- [ ] Add aliases to existing commands:
  - [ ] `list` → `ls`
  - [ ] `delete` → `del`
  - [ ] `open` → `opn`
  - [ ] `render` → `rdr`
  - [ ] `rename` → `rn`
  - [ ] `folder rename` → `rn`
  - [ ] `workspace list` → `workspace ls`
  - [ ] `workspace delete` → `workspace del`
- [ ] Delete old `cmd/mkdir.go` and `cmd/rmdir.go` files
- [ ] Update all test files to use new command names
- [ ] Add tests for alias functionality
- [ ] Update `docs/features/commands.md`:
  - [ ] Document new `folder` subcommand structure
  - [ ] Add alias information to all commands
  - [ ] Update examples with new syntax
- [ ] Update `docs/architecture/overview.md` command table
- [ ] Update `docs/index.md` quick reference
- [ ] Update `.ai/phase1/hotnote_cli_spec.md` with new command structure

### Issue P: Add folder list Command
- [ ] Create `cmd/folder_list.go` with `folder list [path]` command
- [ ] Support optional path argument (defaults to current directory)
- [ ] List all files and folders in specified path
- [ ] JSON output: array of objects with `name`, `path`, `type` fields
- [ ] Human output: tree-like format showing files and folders
- [ ] Add `ls` alias for `folder list`
- [ ] Handle errors: invalid path, path outside workspace, workspace not initialized
- [ ] Unit tests for folder list command
- [ ] Integration tests with JSON and human output
- [ ] Update documentation in `docs/features/commands.md`
- [ ] Update AI agent spec if needed

---

### Legacy: CLI (Phase 1) - Original

#### Core Commands
- [x] new command
- [x] list command
- [x] open command
- [x] render command

#### Workspace Management
- [x] workspace init
- [x] workspace list
- [x] workspace use
- [x] workspace new

---

## TUI (Phase 2)

### Issue 1: TUI Project Setup
- [x] Add dependencies (tview, tcell, goldmark, chroma)
- [x] Create entry point with tview Application
- [x] Set up basic app struct

### Issue 2: Workspace Selection
- [x] Add "tui" command to CLI
- [x] Create workspace selection overlay
- [x] Handle no-workspace case (auto-create default)
- [x] Handle corrupted config (re-init, show error if fails)
- [x] Transition to main view

### Issue 3: 2-Pane Layout + TreeView
- [x] Create Flex layout (25% tree, 75% preview)
- [x] Implement TreeView with folders + files
- [x] Load directory structure from workspace
- [x] Expand/collapse functionality

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

---

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