# Data Flow

This document traces how data moves through the system for each operation.

## Creating a Note

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Store
    participant Workspace
    participant FS as File System

    User->>CLI: hotnote new "My Note"
    CLI->>CLI: slugify("My Note") → "my-note"
    CLI->>Workspace: NewManager()
    CLI->>Store: NewStore(workspace)
    CLI->>Store: Path("my-note")
    Store->>Workspace: Current()
    Workspace-->>Store: (name, path)
    Store-->>CLI: /path/to/my-note.md
    CLI->>FS: Create file with YAML frontmatter
    Note over FS: my-note.md created
    CLI-->>User: "Created note: my-note"
```

### Step-by-Step Details

1. **Slugify**: Convert title to URL-safe string
   ```go
   slug := slugify("My Note")  // "my-note"
   ```

2. **NewManager()**: Create workspace manager
   - Reads config from `~/.config/hotnote/config.yaml`
   - Returns `*Manager` instance

3. **NewStore(wm)**: Create storage with workspace dependency
   - Store holds reference to WorkspaceManager
   - Uses current workspace path for file operations

4. **Ensure()**: Create file if not exists
   - Uses `os.OpenFile` with `O_CREATE|O_EXCL` flags
   - Returns error if file already exists

5. **Write frontmatter**: YAML header with metadata
   ```yaml
   ---
   id: <uuid>
   title: My Note
   created_at: <RFC3339>
   updated_at: <RFC3339>
   tags: []
   ---
   ```

6. **Sync**: Flush to disk

## Listing Notes

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Workspace
    participant FS as File System

    User->>CLI: hotnote list
    CLI->>Workspace: NewManager()
    CLI->>Workspace: Current()
    Workspace-->>CLI: (name, path)
    CLI->>FS: os.ReadDir(path)
    FS-->>CLI: [file list]
    CLI->>CLI: Filter *.md files
    CLI-->>User: [note1.md, note2.md, ...]
```

## Opening a Note

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Store
    participant Workspace
    participant FS as File System
    participant Editor

    User->>CLI: hotnote open "my-note"
    CLI->>Workspace: NewManager()
    CLI->>Store: NewStore(workspace)
    CLI->>Store: Path("my-note")
    Store->>Workspace: Current()
    Workspace-->>Store: (name, path)
    Store-->>CLI: /path/to/my-note.md
    CLI->>FS: Check file exists
    FS-->>CLI: exists
    CLI->>Editor: exec.Command($EDITOR, path).Run()
    Editor->>FS: Read file
    Note over Editor: User edits note
    Editor-->>User: Editor closed
```

## Rendering Markdown

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Store
    participant Workspace
    participant FS as File System
    participant Goldmark

    User->>CLI: hotnote render "my-note"
    CLI->>Workspace: NewManager()
    CLI->>Store: NewStore(workspace)
    CLI->>Store: Path("my-note")
    Store->>Workspace: Current()
    Workspace-->>Store: (name, path)
    Store-->>CLI: /path/to/my-note.md
    CLI->>FS: Read file content
    FS-->>CLI: markdown content
    CLI->>Goldmark: Parse markdown
    Note over Goldmark: Convert to HTML AST
    Goldmark-->>CLI: HTML output
    CLI-->>User: <h1>My Note</h1>...
```

### Step-by-Step Details

1. **Read file**: Load markdown content from workspace
2. **Parse**: Goldmark converts markdown to AST
3. **Render**: Output HTML to stdout

## Workspace Switching

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Manager
    participant Config as config.yaml

    User->>CLI: hotnote workspace use notes
    CLI->>Manager: NewManager()
    Manager->>Config: Load config
    Config-->>Manager: {current: "default", ...}
    CLI->>Manager: Use("notes")
    Manager->>Config: Update current="notes"
    Config-->>Manager: Saved
    Manager-->>User: Workspace switched
```
