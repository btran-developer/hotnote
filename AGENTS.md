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

## Concurrency
- Channels over mutexes
- Context for cancellation
- No global state
- sync/atomic for simple counters

## Hotnote Specific
- UTF-8 markdown
- UUIDs: github.com/google/uuid
- Slugs: lowercase, hyphen, ASCII
- Frontmatter: YAML