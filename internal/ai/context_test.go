package ai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewContextBuilder(t *testing.T) {
	// Test with custom values
	builder := NewContextBuilder(10, 1000, 500)
	if builder.MaxContextNotes != 10 {
		t.Errorf("Expected MaxContextNotes 10, got %d", builder.MaxContextNotes)
	}
	if builder.MaxTokens != 1000 {
		t.Errorf("Expected MaxTokens 1000, got %d", builder.MaxTokens)
	}
	if builder.CharLimit != 500 {
		t.Errorf("Expected CharLimit 500, got %d", builder.CharLimit)
	}

	// Test with defaults (zero values)
	builder = NewContextBuilder(0, 0, 0)
	if builder.MaxContextNotes != 20 {
		t.Errorf("Expected default MaxContextNotes 20, got %d", builder.MaxContextNotes)
	}
	if builder.MaxTokens != 6000 {
		t.Errorf("Expected default MaxTokens 6000, got %d", builder.MaxTokens)
	}
	if builder.CharLimit != 2000 {
		t.Errorf("Expected default CharLimit 2000, got %d", builder.CharLimit)
	}
}

func TestScoreNotes(t *testing.T) {
	notes := []NoteIndex{
		{
			Slug:    "test-note",
			Title:   "Test Note",
			Tags:    []string{"test", "example"},
			Excerpt: "This is a test note about testing",
			Path:    "/test/path.md",
		},
	}

	scored := ScoreNotes(notes, "test")
	if len(scored) != 1 {
		t.Fatalf("Expected 1 scored note, got %d", len(scored))
	}

	// Should have high score since "test" matches title, tags, and excerpt
	if scored[0].Score == 0 {
		t.Error("Expected non-zero score for matching note")
	}
}

func TestExtractExcerpt(t *testing.T) {
	// Test with frontmatter
	content := "---\ntitle: Test\n---\nThis is the actual content of the note."
	excerpt := extractExcerpt(content)
	if excerpt != "This is the actual content of the note." {
		t.Errorf("Expected content without frontmatter, got: %s", excerpt)
	}

	// Test without frontmatter
	content = "Just plain content here"
	excerpt = extractExcerpt(content)
	if excerpt != "Just plain content here" {
		t.Errorf("Expected plain content, got: %s", excerpt)
	}

	// Test long content
	longContent := make([]byte, 1000)
	for i := range longContent {
		longContent[i] = 'a'
	}
	excerpt = extractExcerpt(string(longContent))
	if len(excerpt) != 500 {
		t.Errorf("Expected excerpt length 500, got %d", len(excerpt))
	}
}

func TestBuildNoteIndex(t *testing.T) {
	// Create a temporary directory with test files
	tempDir := t.TempDir()

	// Create a test note
	testContent := "---\ntitle: Test Note\ntags:\n  - test\n---\n# Test\n\nThis is test content."
	testFile := filepath.Join(tempDir, "test-note.md")
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	index, err := BuildNoteIndex(tempDir)
	if err != nil {
		t.Fatalf("BuildNoteIndex failed: %v", err)
	}

	if len(index) != 1 {
		t.Errorf("Expected 1 note in index, got %d", len(index))
	}

	if len(index) > 0 {
		note := index[0]
		if note.Slug != "test-note" {
			t.Errorf("Expected slug 'test-note', got '%s'", note.Slug)
		}
		if note.Title != "Test Note" {
			t.Errorf("Expected title 'Test Note', got '%s'", note.Title)
		}
		if len(note.Tags) != 1 || note.Tags[0] != "test" {
			t.Errorf("Expected tags [test], got %v", note.Tags)
		}
	}
}
