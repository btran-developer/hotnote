// Package frontmatter provides YAML frontmatter extraction and parsing.
package frontmatter

import (
	"bytes"
	"time"

	"gopkg.in/yaml.v3"
)

// Extract parses YAML frontmatter from raw markdown content.
// It returns the parsed frontmatter data and true if valid frontmatter
// was found, or nil and false otherwise.
func Extract(content []byte) (map[string]interface{}, bool) {
	prefix := []byte("---\n")
	if !bytes.HasPrefix(content, prefix) {
		return nil, false
	}

	// Find the closing delimiter
	rest := content[len(prefix):]
	end := bytes.Index(rest, []byte("\n---"))
	if end == -1 {
		return nil, false
	}

	frontmatterBytes := rest[:end]
	var data map[string]interface{}
	if err := yaml.Unmarshal(frontmatterBytes, &data); err != nil {
		return nil, false
	}

	return data, true
}

// ParseCreatedAt extracts the created_at field from frontmatter data.
// It handles both RFC3339 strings (written by the create command) and
// time.Time values (from YAML parsing).
// Returns the parsed time and true, or zero time and false if absent/invalid.
func ParseCreatedAt(data map[string]interface{}) (time.Time, bool) {
	if data == nil {
		return time.Time{}, false
	}
	raw, ok := data["created_at"]
	if !ok {
		return time.Time{}, false
	}

	switch v := raw.(type) {
	case time.Time:
		return v, true
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	default:
		return time.Time{}, false
	}
}
