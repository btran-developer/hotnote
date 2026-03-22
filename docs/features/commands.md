# Commands

Hotnote CLI commands for managing notes.

## Global Flags

These flags work with all commands:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--data-dir` | `-d` | `notes` | Data directory for notes |
| `--json` | | `false` | Output in JSON format |
| `--version` | | `false` | Show version information |
| `--help` | `-h` | `false` | Show help for any command |

## Exit Codes

Hotnote returns the following exit codes:

| Code | Meaning | Example |
|------|---------|---------|
| 0 | Success | Command executed successfully |
| 1 | General error | Unexpected error occurred |
| 2 | Not found | Note or workspace not found |
| 3 | Invalid input | Missing required arguments |
| 4 | Config error | Workspace not initialized |

## hotnote new

Create a new note with YAML frontmatter.

### Usage

```bash
hotnote new [title] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `title` | Title of the note (will be slugified) |

### Examples

```bash
# Create a note titled "My Research"
hotnote new "My Research"

# Create a note with special characters
hotnote new "Q3 2024 Goals!"
```

### How It Works

1. Converts title to slug (URL-safe string)
   - `"My Research"` → `"my-research"`
   - Special characters removed
   - Spaces converted to hyphens
   - Lowercased

2. Creates file at `~/.local/share/hotnote/workspaces/<current>/<slug>.md`

3. Writes YAML frontmatter:
   ```yaml
   ---
   id: <uuid>
   title: My Research
   created_at: 2024-01-15T10:30:00Z
   updated_at: 2024-01-15T10:30:00Z
   tags: []
   ---
   ```

## hotnote list

List all notes in the current workspace.

### Usage

```bash
hotnote list [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--flat` | Output plain filenames only |
| `--sort updated` | Sort by last modified time |
| `--sort created` | Sort by creation time |
| `--json` | Output in JSON format |

### Examples

```bash
# List all notes
hotnote list

# List in JSON format
hotnote list --json

# Sort by updated time
hotnote list --sort updated
```

### Output Format

**Plain text:**
```
note1.md
note2.md
project-ideas.md
```

**JSON (`--json`):**
```json
["note1.md", "note2.md", "project-ideas.md"]
```

### Notes

- Only `.md` files are listed
- Files are read from the current workspace directory
- Does not show content, only filenames

## hotnote open

Open a note in the default editor.

### Usage

```bash
hotnote open [title] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `title` | Title or slug of the note |

### Examples

```bash
# Open note by title
hotnote open "My Research"

# Open note by slug
hotnote open my-research
```

### How It Works

1. Looks up the note file path from the slug
2. Verifies the file exists
3. Launches `$EDITOR` with the file path
4. Waits for the editor to close

### Editor Configuration

| Environment | Default |
|-------------|---------|
| `$EDITOR` | `vim` |

Set your preferred editor:
```bash
export EDITOR=code        # VS Code
export EDITOR=emacs       # Emacs
export EDITOR=nano        # Nano
```

## hotnote render

Render a note as HTML.

### Usage

```bash
hotnote render [title] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `title` | Title or slug of the note |

### Examples

```bash
# Render to HTML
hotnote render "My Research"

# Pipe to file
hotnote render "My Research" > output.html
```

### Supported Markdown

| Element | Example | Output |
|---------|---------|--------|
| Headings | `# H1` | `<h1>H1</h1>` |
| Bold | `**text**` | `<strong>text</strong>` |
| Italic | `*text*` | `<em>text</em>` |
| Links | `[text](url)` | `<a href="url">text</a>` |
| Code | `` `code` `` | `<code>code</code>` |
| Code blocks | ` ``` ` | `<pre><code>` |

### How It Works

Uses [Goldmark](https://github.com/yuin/goldmark) for parsing and conversion.

## hotnote workspace

Manage multiple note workspaces.

### Usage

```bash
hotnote workspace <subcommand> [flags]
```

### Subcommands

#### init

Initialize the default workspace.

```bash
hotnote workspace init
```

Creates:
- `~/.config/hotnote/config.yaml`
- `~/.local/share/hotnote/workspaces/default/`

#### list

List all workspaces.

```bash
hotnote workspace list
```

Output:
```
default  *
work     *
```

Current workspace marked with `*`.

#### new

Create a new workspace.

```bash
hotnote workspace new <name> [flags]
```

| Flag | Description |
|------|-------------|
| `--path` | Custom path for workspace |

```bash
# Create workspace in default location
hotnote workspace new work

# Create workspace in custom location
hotnote workspace new archive --path /mnt/notes/archive
```

#### use

Switch to a different workspace.

```bash
hotnote workspace use <name>
```

```bash
hotnote workspace use work
```

All subsequent `new`, `list`, `open`, and `render` commands will use the selected workspace.

## hotnote version

Display version information.

### Usage

```bash
hotnote --version
```

### Example

```bash
hotnote --version
# Output: hotnote version 0.1.0
```

## hotnote ai

AI-powered note operations (Phase 2 feature, stub implementation).

### Usage

```bash
hotnote ai <subcommand> [flags]
```

### Status

This command is currently a stub. Full AI integration is planned for Phase 2.

## Adding Commands

To add a new command, follow this pattern:

1. Create a new file: `cmd/<name>.go`
2. Define the command variable:

   ```go
   var <Name>Cmd = &cobra.Command{
       Use:   "<name>",
       Short: "Description of what it does",
       Run:   run<Name>,
   }

   func run<Name>(cmd *cobra.Command, args []string) {
       // Implementation
   }
   ```

3. Register in `init()` at the bottom of the file:

   ```go
   func init() {
       RootCmd.AddCommand(<Name>Cmd)
   }
   ```

### Conventions

| Item | Convention |
|------|------------|
| Filename | `<name>.go` (lowercase) |
| Variable | `<Name>Cmd` (PascalCase + "Cmd") |
| Use string | `<name>` (lowercase) |
| Init function | Bottom of file |

### Example

For a `delete` command:

```go
// cmd/delete.go
package cmd

var DeleteCmd = &cobra.Command{
    Use:   "delete <note>",
    Short: "Delete a note",
    Args:  cobra.ExactArgs(1),
    Run:   runDelete,
}

func runDelete(cmd *cobra.Command, args []string) {
    // Implementation
}

func init() {
    RootCmd.AddCommand(DeleteCmd)
}
```
