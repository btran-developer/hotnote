# Go Coding Guidelines for Hotnote

## Formatting
- `go fmt`; max 120 chars; tabs; blank lines between sections

## Naming
- Packages: lowercase, single word
- Interfaces: -er suffix (Reader, Writer)
- MixedCaps; exported: uppercase first letter

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
