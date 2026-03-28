# Go Coding Guidelines for Hotnote

## Formatting
- `go fmt`; max 120 chars; tabs; blank lines between sections

## Naming
- Packages: lowercase, single word
- Interfaces: -er suffix (Reader, Writer)
- MixedCaps; exported: uppercase first letter

## Constants & Palette
- Prefer constants over magic strings/numbers
- Use typed constants (e.g., `type ViewMode string`) for enum-like values
- Group related constants in a `const` block
- Document allowed values when not self-explanatory
- Use `DefaultPalette` colors instead of inline color tags (e.g., `SetTextColor(palette.Red)` not `[red]`)
- Avoid inline tview color tags like `[red]`, `[blue]` — prefer `SetTextColor()` with palette

## Imports
- std, 3rd-party, local groups; blank between; alphabetical within

## Errors
- Check explicitly; wrap: `fmt.Errorf("op: %w", err)`
- Sentinel: `var ErrNotFound = errors.New("not found")`

## Comments

Comments are written **with** the code, not after. Treat doc comments like tests —
every new exported symbol, sentinel error, and interface requires a doc comment
before the change is complete.

### Format Reference

| Artifact | Format | Example |
|----------|--------|---------|
| Package | `// Package <name> provides <summary>.` | `// Package storage provides atomic note file storage.` |
| Exported type | `// <Name> <verb>s ...` | `// Store manages note files on disk.` |
| Exported func | `// <Func> <verb>s ...` | `// NewStore creates a Store using the given WorkspaceManager.` |
| Interface | `// <I> <verb>s ...` + method docs | `// WorkspaceManager provides access to the active workspace.` |
| Interface method | `// <Method> <verb>s ...` | `// Current returns the name and path of the active workspace.` |
| Sentinel error | `// Err<X> is returned when <condition>.` | `// ErrNotFound is returned when the note does not exist.` |
| Constant | `// <Const> <verb>s ...` | `// ViewModeTree indicates the file tree panel has focus.` |
| Struct field | `// <Field> <verb>s ...` | `// Slug is the URL-safe identifier derived from the filename.` |
| Unexported func | `// <func> <verb>s ...` | `// load reads the configuration from disk.` |
| Inline | sparse; explain *why*, not *what* | `// Retry once for transient lock conflicts.` |
| TODO | `// TODO(owner): description` | `// TODO(btran): add batch delete support` |

### Comment Checklist (per code change)

- [ ] New package has a `// Package` comment in its primary file
- [ ] All new exported types, functions, and methods have doc comments
- [ ] New sentinel errors document when they are returned
- [ ] New interfaces and their methods are documented
- [ ] Non-obvious exported constants have doc comments
- [ ] Complex unexported functions have doc comments
- [ ] Inline comments explain *why*, not *what*

## Performance
- Pre-allocate when size known
- `strings.Builder` for concat
- Profile before optimizing

## Testing
- Table-driven; `*_test.go`
- Mock interfaces- Benchmark only when needed
- Always add unit tests for new code changes
- Clean up artifacts using `t.Cleanup()` or `defer`
- Do not leave test data in user directories (~/.config, ~/.local)

## Verification
- Run `go vet ./...` to catch basic issues
- Run `go test -v ./...` before committing
- Review doc comments as part of code review — same rigor as test coverage
- For manual CLI verification, clean up after:
  - Config files: `rm -f ~/.config/hotnote/config.yaml`
  - Workspace directories: `rm -rf ~/.local/share/hotnote/workspaces/`
  - Built binaries: `rm -f ./hotnote` or `rm -f cmd/hotnote/hotnote`

## Concurrency
- Channels over mutexes
- Context for cancellation
- No global state
- sync/atomic for simple counters

## Hotnote Specific- UTF-8 markdown
- UUIDs: github.com/google/uuid
- Slugs: lowercase, hyphen, ASCII
- Frontmatter: YAML

## Architecture

### Package Structure
```
cmd/              # CLI entry points (Cobra commands)
internal/tui/    # TUI application (tview)
internal/<name>/ # Shared packages (storage, workspace, core, etc.)
docs/            # User-facing documentation
.ai/             # Design docs and planning (for agents)
```

### Package Guidelines
- New shared logic goes in `internal/<name>/` first
- CLI and TUI both import from `internal/` packages
- Avoid circular dependencies between packages
- Use interfaces in `internal/` for testability

### Separation of Concerns
- CLI (`cmd/`) only handles command parsing and output formatting
- TUI (`internal/tui/`) handles presentation and user interaction
- Business logic lives in `internal/` packages
- Reuse business logic across interfaces, don't duplicate

### Configuration
- CLI flags for command-specific options
- Config file (`~/.config/hotnote/config.yaml`) for persistent settings
- Environment variables as override (e.g., $HOTNOTE_DATA_DIR)

### Dependency Philosophy
- Prefer Go standard library
- Add external dependencies only when necessary (e.g., tview, cobra, goldmark)
- Document why each dependency is needed in go.mod comment if non-obvious

### Interface Usage
- Define interfaces in `internal/` packages for testability
- Use concrete types in `cmd/` and `internal/tui/`
- Example: storage.Store interface → CLI and TUI both use it

## Git Commit Conventions
- Follow Conventional Commits (conventionalcommits.org)
- Format: `<type>(<scope>): <description>`
- Use imperative mood: "add" not "added"
- Body (bullets): only when multiple changes

Single: feat(commands): add delete command

Multiple:
feat(commands): add delete command
- Implement runDelete function
- Add --force flag

Examples: feat, fix, docs, refactor, test, chore

## Documentation

Documentation must stay accurate with code changes. After modifying code, check for gaps.

### Update Triggers

| Change | Check |
|--------|-------|
| Command | `docs/features/commands.md` |
| Storage | `docs/features/storage.md` |
| Workspace | `docs/features/workspace.md` |
| New package | `docs/architecture/overview.md` |
| Errors | `docs/architecture/error-handling.md` |

### Create New Docs When

- New feature without docs → create `docs/features/<name>.md`
- New architectural component → add to `docs/architecture/`
- Significant behavior change → update `.ai/implementation-issues.md`

### Checklist

- [ ] Update existing docs for changed behavior
- [ ] Create new docs for undocumented features
- [ ] Update project structure in `docs/architecture/overview.md`
- [ ] Update links in `docs/index.md`
- [ ] Sync `.ai/` design docs

### Avoid

- Orphaned docs (feature not implemented)
- Stale code examples
- Broken links
