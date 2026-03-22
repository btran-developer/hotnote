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
