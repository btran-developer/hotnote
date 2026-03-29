package frontmatter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtract_ValidFrontmatter(t *testing.T) {
	content := []byte("---\nid: abc-123\ntitle: Test\ncreated_at: 2026-03-28T14:30:00Z\n---\n\n# Test\n")
	data, ok := Extract(content)
	require.True(t, ok)
	assert.Equal(t, "abc-123", data["id"])
	assert.Equal(t, "Test", data["title"])
}

func TestExtract_NoFrontmatter(t *testing.T) {
	content := []byte("# Just markdown\n")
	_, ok := Extract(content)
	assert.False(t, ok)
}

func TestExtract_IncompleteFrontmatter(t *testing.T) {
	content := []byte("---\nid: abc\n")
	_, ok := Extract(content)
	assert.False(t, ok)
}

func TestExtract_EmptyFrontmatter(t *testing.T) {
	content := []byte("---\n---\n\n# Body\n")
	_, ok := Extract(content)
	assert.False(t, ok)
}

func TestExtract_InvalidYAML(t *testing.T) {
	content := []byte("---\n{{invalid yaml\n---\n\n# Body\n")
	_, ok := Extract(content)
	assert.False(t, ok)
}

func TestExtract_MinimalFrontmatter(t *testing.T) {
	content := []byte("---\nid: abc-123\n---\n")
	data, ok := Extract(content)
	require.True(t, ok)
	assert.Equal(t, "abc-123", data["id"])
}

func TestParseCreatedAt_String(t *testing.T) {
	data := map[string]interface{}{
		"created_at": "2026-03-28T14:30:00Z",
	}
	got, ok := ParseCreatedAt(data)
	require.True(t, ok)
	expected := time.Date(2026, 3, 28, 14, 30, 0, 0, time.UTC)
	assert.Equal(t, expected, got)
}

func TestParseCreatedAt_TimeValue(t *testing.T) {
	expected := time.Date(2026, 3, 28, 14, 30, 0, 0, time.UTC)
	data := map[string]interface{}{
		"created_at": expected,
	}
	got, ok := ParseCreatedAt(data)
	require.True(t, ok)
	assert.Equal(t, expected, got)
}

func TestParseCreatedAt_Missing(t *testing.T) {
	data := map[string]interface{}{
		"id": "abc-123",
	}
	_, ok := ParseCreatedAt(data)
	assert.False(t, ok)
}

func TestParseCreatedAt_InvalidString(t *testing.T) {
	data := map[string]interface{}{
		"created_at": "not-a-timestamp",
	}
	_, ok := ParseCreatedAt(data)
	assert.False(t, ok)
}

func TestParseCreatedAt_WrongType(t *testing.T) {
	data := map[string]interface{}{
		"created_at": 42,
	}
	_, ok := ParseCreatedAt(data)
	assert.False(t, ok)
}

func TestParseCreatedAt_NilData(t *testing.T) {
	_, ok := ParseCreatedAt(nil)
	assert.False(t, ok)
}
