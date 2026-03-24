# User Journey — HotNote CLI

## Journey 1: Create Note
1. Run `hotnote new "My Note Title"`
2. Slug generated from title (e.g., "my-note-title")
3. Note created in workspace notes directory
4. Opens in $EDITOR
5. User writes content and saves/exits

## Journey 2: List Notes
1. Run `hotnote list`
2. Shows all notes in workspace
3. Display: `<path>...........<date>`
4. Sorted by modified date (newest first)

## Journey 3: Open Note
1. Run `hotnote open <slug>`
2. Finds note by slug (supports partial match)
3. Opens in $EDITOR

## Journey 4: Render Note
1. Run `hotnote render <slug>`
2. Renders markdown to terminal
3. Supports full markdown features

## Journey 5: Workspace Management
- `hotnote workspace init` - Initialize default workspace
- `hotnote workspace list` - List all workspaces
- `hotnote workspace use <name>` - Switch to workspace
- `hotnote workspace new <name>` - Create new workspace

## Journey 6: Help
- `hotnote --help` - Show help
- `hotnote <command> --help` - Show command help

## Journey 7: Output Modes
- Default: Human-readable output
- `--json` flag: JSON output
- `--pretty` flag: Pretty-printed JSON