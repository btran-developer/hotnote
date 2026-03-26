package slugify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugify_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase",
			input:    "hello world",
			expected: "hello-world",
		},
		{
			name:     "with uppercase",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "with special characters",
			input:    "Test!@#Note",
			expected: "testnote",
		},
		{
			name:     "with numbers",
			input:    "Note 123 Test",
			expected: "note-123-test",
		},
		{
			name:     "already slug",
			input:    "already-slug",
			expected: "already-slug",
		},
		{
			name:     "multiple spaces",
			input:    "A  B   C",
			expected: "a-b-c",
		},
		{
			name:     "consecutive hyphens from input",
			input:    "Note---Test",
			expected: "note-test",
		},
		{
			name:     "with punctuation",
			input:    "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "with underscores",
			input:    "hello_world",
			expected: "helloworld",
		},
		{
			name:     "with parentheses",
			input:    "Note (1)",
			expected: "note-1",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "numbers only",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "unicode characters",
			input:    "café",
			expected: "caf",
		},
		{
			name:     "leading and trailing spaces",
			input:    "  leading and trailing  ",
			expected: "leading-and-trailing",
		},
		{
			name:     "mixed case with numbers",
			input:    "Go 1.0 Release",
			expected: "go-10-release",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slugify(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
