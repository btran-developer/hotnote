# HotNote CLI Specification (Phase 1 — Final)

## 1. Global Behavior

### Command Format
hotnote <command> [args] [flags]

### Global Flags
--help        Show help
--version     Show version
--json        Output in JSON format (when applicable)

Rules:
- --json overrides human-readable output
- No extra logs in JSON mode

### Exit Codes
0 Success
1 General error
2 Not found
3 Invalid input
4 Config error

---

## 2. hotnote new

Usage:
hotnote new <title>

Behavior:
- Generate slug
- Auto-resolve collision (default)
- Create note under <workspace>/notes/
- Open in $EDITOR

Flags:
--no-open
--path <dir>
--slug <slug>
--no-auto-slug

Output:
Created note: project-a/my-idea.md

JSON:
{
  "status": "created",
  "slug": "my-idea",
  "path": "project-a/my-idea.md"
}

Errors:
- title is required
- path does not exist
- slug exists (if --no-auto-slug)

---

## 3. hotnote list

Usage:
hotnote list

Flags:
--flat
--sort updated
--sort created
--json

Output:
project-a/my-idea        2026-03-21
project-a/design         2026-03-20

JSON:
[
  {
    "slug": "my-idea",
    "path": "project-a/my-idea.md",
    "updated_at": "2026-03-21T10:00:00Z"
  }
]

Errors:
- no workspace configured

---

## 4. hotnote open

Usage:
hotnote open <slug>

Output:
(opens editor)

JSON:
{
  "status": "opened",
  "path": "project-a/my-idea.md"
}

Errors:
- note not found
- multiple matches

---

## 5. hotnote render

Usage:
hotnote render <slug>

Output:
(rendered markdown)

JSON:
{
  "content": "# My Idea\n\nThis is a note."
}

Errors:
- note not found

---

## 6. Workspace

Model:
<workspace>/
  notes/

Commands:
hotnote workspace init
hotnote workspace list
hotnote workspace use <name>
hotnote workspace new <name> [--path <path>]

---

## workspace init

Output:
Initialized workspace: default

Error:
workspace already initialized

---

## workspace list

Output:
* default    ~/.local/share/hotnote/workspaces/default

JSON:
{
  "current": "default",
  "workspaces": {}
}

---

## workspace use

Usage:
hotnote workspace use proj-a

Output:
Switched to workspace: proj-a

Error:
workspace not found

---

## workspace new

Usage:
hotnote workspace new proj-a

Output:
Created workspace: proj-a

Error:
workspace already exists

---

## 7. Rules

- Atomic file writes
- Deterministic output
- Human + JSON modes supported

---

## 8. hotnote mkdir

Usage:
hotnote mkdir <folder>

Behavior:
- Create folder in current workspace
- Support nested paths (e.g., `projects/2024`)
- Auto-create parent folders if needed

Flags:
--json
--pretty

Output:
Created folder: projects

JSON:
{
  "status": "created",
  "folder": "projects",
  "path": "~/.local/share/hotnote/workspaces/default/projects"
}

Errors:
- folder already exists
- workspace not initialized

---

## 9. hotnote rmdir

Usage:
hotnote rmdir <folder>

Behavior:
- Delete folder from current workspace
- Prompt for confirmation if folder not empty
- Skip prompt with --force

Flags:
--force
--json
--pretty

Output:
Delete folder 'projects' and all contents? [y/N]: y
Deleted folder: projects

JSON:
{
  "status": "deleted",
  "folder": "projects"
}

Errors:
- folder not found
- cannot delete workspace root

---

## 10. hotnote delete

Usage:
hotnote delete <slug>

Behavior:
- Delete note from current workspace
- Hybrid path resolution (see Section 12)
- Prompt for confirmation
- Skip prompt with --force

Flags:
--force
--json
--pretty

Output:
Delete note 'my-idea'? [y/N]: y
Deleted note: my-idea

JSON:
{
  "status": "deleted",
  "slug": "my-idea"
}

Errors:
- note not found
- multiple matches (show list)

---

## 11. workspace delete

Usage:
hotnote workspace delete <name>

Behavior:
- Delete workspace and all contents
- Recursively remove directory
- Prompt for confirmation
- Skip prompt with --force

Flags:
--force
--json
--pretty

Output:
Delete workspace 'work' and all contents? [y/N]: y
Deleted workspace: work

JSON:
{
  "status": "deleted",
  "workspace": "work"
}

Errors:
- workspace not found
- cannot delete current workspace (switch first)
- cannot delete last workspace

---

## 12. Hybrid Path Resolution

For commands that accept a slug/path argument (open, render, delete):

Resolution Logic:
1. If input contains `/`, treat as direct relative path to workspace root
2. If input does not contain `/`, search recursively by slug
3. If multiple matches found, show list and prompt user to specify full path

Examples:
hotnote open my-idea           # Searches recursively
hotnote open projects/my-idea  # Direct path lookup
hotnote render my-idea         # Searches recursively
hotnote render projects/my-idea # Direct path lookup
