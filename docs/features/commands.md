# Commands

Hotnote CLI commands for managing notes.

## Global Flags

These flags work with all commands:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--data-dir` | `-d` | `notes` | Data directory for notes |
| `--json` | | `false` | Output in JSON format |

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

### Examples

```bash
# List all notes
hotnote list

# List in JSON format
hotnote list --json
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

### Workspace Storage

| Purpose | Location |
|---------|----------|
| Config | `~/.config/hotnote/config.yaml` |
| Notes | `~/.local/share/hotnote/workspaces/<name>/` |
