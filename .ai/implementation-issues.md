# HotNote Implementation Issues

This document breaks down the work needed to align the current implementation with the requirements specified in:
- ./.ai/hotnote_phase1_prd.md (Phase 1 PRD)
- ./.ai/hotnote_cli_spec.md (CLI Specification)

Each issue is designed to be manageable for individual AI agents to implement.

## Implementation Status

### Completed Issues (as of 2026-03-22)
- ✅ **Workspace functionality** - All workspace commands implemented (init, list, use, new)
- ✅ **Frontmatter implementation** - Notes include UUID, title, timestamps, and tags
- ✅ **Editor integration** - Basic $EDITOR integration with vim fallback
- ✅ **JSON output support** - List command supports --json flag
- ✅ **Render JSON output** - Render command supports --json flag
- ✅ **Version flag** - Global --version flag implemented
- ✅ **Workspace JSON output** - All workspace commands support --json flag
- ✅ **Error exit codes** - Proper exit codes (0-4) implemented per CLI Spec
- ✅ **Standardized error messages** - Error messages match CLI Spec formats (lowercase, no "Error:" prefix)

### Remaining Issues

### 5. Error Handling Improvement (Partially Missing)
Based on CLI Spec §17-23 and throughout

**Issue 5.1: Implement proper error codes**
- Define and use consistent error codes:
  - 0: Success
  - 1: General error
  - 2: Not found
  - 3: Invalid input
  - 4: Config error
- Ensure all commands return appropriate exit codes

**Issue 5.2: Standardize error messages** ✅ Complete
- Make error messages clear and consistent
- Match expected formats from CLI Spec where specified
- Include relevant context (e.g., note name, path)

**Issue 5.3: Improve error wrapping** ✅ Complete
- Continue using `fmt.Errorf("operation: %w", err)` pattern
- Add context to all errors from storage layer
- Add context to all errors from workspace layer
- Ensure errors are meaningful to users

### 7. Slugify Enhancement (Partially Missing)
Based on PRD §63 and CLI Spec §32

**Issue 7.1: Improve slug generation**
- Enhance slugify function to properly handle special characters
- Convert to lowercase
- Replace spaces with hyphens
- Remove special characters (keep only alphanumeric and hyphens)
- Handle edge cases like multiple consecutive hyphens
- Example: "My Research Plan!" → "my-research-plan"

### 8. Atomic File Writes (Missing)
Based on CLI Spec §191

**Issue 8.1: Implement atomic file writes**
- When writing note files, use temporary files and rename
- Prevents partial writes if process is interrupted
- Important for reliability (PRD §164)

### 9. Deterministic Output (Partially Missing)
Based on PRD §167 and CLI Spec §193

**Issue 9.1: Ensure deterministic output**
- Make sure command output is consistent and predictable
- For list command, ensure consistent ordering when no sort specified
- Format dates/times consistently
- Ensure JSON output is properly formatted and valid

### 10. Performance Optimization (Ongoing)
Based on PRD §163

**Issue 10.1: Ensure fast operation (<100ms)**
- Profile key operations (new, list, open, render)
- Optimize where necessary
- Cache goldmark renderer instance (already partially done)
- Optimize filesystem operations

## Suggested Implementation Order

1. **JSON output for render and workspace commands** (Issues 4.2-4.3) - Quick wins
2. **Error handling improvements** (Issues 5.1-5.3) - Important for robustness
3. **Slugify enhancement** (Issue 7.1) - Improves correctness
4. **Atomic file writes** (Issue 8.1) - Improves reliability
5. **Deterministic output** (Issue 9.1) - Ensures consistent behavior
6. **Performance optimization** (Issue 10.1) - Ongoing improvement

Each issue should be implemented as a separate commit with clear, descriptive commit messages following conventional commits specification.

## Verification Approach

After implementing each issue, AI agents should:
1. Run the verification test cases from verification-plan.md related to that issue
2. Ensure no regressions in existing functionality
3. Verify that the implementation matches the exact specifications
4. Test edge cases where applicable