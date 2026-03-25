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
- Package: file-level before package
- Exports: complete sentences
- Inline: sparse; explain why- TODO: `// TODO(owner): desc`

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
- Run `go test -v ./...` before committing
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
