# Version

Hotnote uses semantic versioning (SemVer) for releases.

## Current Version

```
hotnote version 0.1.0
```

## Version Format

`<major>.<minor>.<patch>`

| Part | Meaning |
|------|---------|
| Major | Breaking changes |
| Minor | New features (backward compatible) |
| Patch | Bug fixes |

## Version Lifecycle

| Phase | Version | Status |
|-------|---------|--------|
| MVP | 0.1.x | Current |
| Phase 1 | 0.x.x | In progress |
| Phase 2 | 1.x.x | Planned |

## Viewing the Version

### CLI

```bash
hotnote --version
```

Output:
```
hotnote version 0.1.0
```

### Programmatic Access

The version is defined in `cmd/root.go`:

```go
var RootCmd = &cobra.Command{
    Use:     "hotnote",
    Short:   "A terminal-first markdown note system",
    Version: "0.1.0",
}
```

## Versioning Strategy

### Phase 1 (Current)

Focus on core functionality:
- File-based storage
- Basic CRUD operations
- Workspace management
- Editor integration

### Phase 2 (Planned)

Add AI and advanced features:
- AI-powered search
- Note summarization
- Smart organization

### Phase 3 (Future)

Complete feature set:
- TUI interface
- Plugin system
- Sync and collaboration

## Release Process

1. Update version in `cmd/root.go`
2. Update `CHANGELOG.md`
3. Create git tag: `v0.1.0`
4. Build and test
5. Create GitHub release

## Build with Version

```bash
go build -ldflags="-X main.version=0.2.0" ./cmd/hotnote
```
