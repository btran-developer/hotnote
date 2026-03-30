package ai

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"hotnotego/internal/frontmatter"
)

// NoteIndex holds indexed metadata about a note.
type NoteIndex struct {
	Slug      string
	Path      string
	Title     string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
	Excerpt   string
	WordCount int
}

// ContextBuilder assembles context strings from scored notes.
type ContextBuilder struct {
	MaxContextNotes int
	MaxTokens       int
	CharLimit       int
}

// NewContextBuilder creates a ContextBuilder with defaults.
func NewContextBuilder(maxNotes, maxTokens, charLimit int) *ContextBuilder {
	if maxNotes == 0 {
		maxNotes = 20
	}
	if maxTokens == 0 {
		maxTokens = 6000
	}
	if charLimit == 0 {
		charLimit = 2000
	}
	return &ContextBuilder{
		MaxContextNotes: maxNotes,
		MaxTokens:       maxTokens,
		CharLimit:       charLimit,
	}
}

// BuildNoteIndex scans a workspace and indexes all notes.
func BuildNoteIndex(workspacePath string) ([]NoteIndex, error) {
	var index []NoteIndex

	err := filepath.WalkDir(workspacePath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".md") {
			relPath, _ := filepath.Rel(workspacePath, path) // Should not fail for workspace children
			slug := strings.TrimSuffix(relPath, ".md")
			slug = filepath.ToSlash(slug)

			content, err := os.ReadFile(path)
			if err != nil {
				return nil // Skip unreadable files, continue to next
			}

			meta, _ := frontmatter.Extract(content) // Returns (data, ok); nil handled by fallbacks

			// Get FileInfo for timestamps
			fileInfo, err := info.Info()
			if err != nil {
				log.Printf("Warning: failed to get file info for %s: %v", path, err)
				return nil
			}

			noteIndex := NoteIndex{
				Slug:      slug,
				Path:      path,
				Title:     getString(meta, "title", slug),
				Tags:      getStringSlice(meta, "tags"),
				CreatedAt: getTime(meta, "created_at", fileInfo.ModTime()),
				UpdatedAt: fileInfo.ModTime(),
				WordCount: len(strings.Fields(string(content))),
			}

			noteIndex.Excerpt = extractExcerpt(string(content))

			index = append(index, noteIndex)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("build index: %w", err)
	}

	return index, nil
}

// ScoredNote pairs a NoteIndex with its relevance score.
type ScoredNote struct {
	NoteIndex
	Score float64
}

// ScoreNotes ranks notes by relevance to a query.
func ScoreNotes(index []NoteIndex, query string) []ScoredNote {
	var scored []ScoredNote
	queryLower := strings.ToLower(query)

	for _, note := range index {
		score := calculateRelevance(note, queryLower)
		scored = append(scored, ScoredNote{NoteIndex: note, Score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	return scored
}

func calculateRelevance(note NoteIndex, query string) float64 {
	score := 0.0

	if strings.Contains(strings.ToLower(note.Title), query) {
		score += 0.30
	}

	for _, tag := range note.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			score += 0.25
			break
		}
	}

	if strings.Contains(strings.ToLower(note.Excerpt), query) {
		score += 0.30
	}

	daysOld := time.Since(note.UpdatedAt).Hours() / 24
	recencyScore := math.Exp(-daysOld/30) * 0.10
	score += recencyScore

	if strings.Contains(strings.ToLower(note.Path), query) {
		score += 0.05
	}

	return score
}

// BuildContext assembles a context string from scored notes.
func (b *ContextBuilder) BuildContext(notes []ScoredNote) (string, []string, error) {
	var included []string
	var context strings.Builder
	totalTokens := 0

	for i, note := range notes {
		if i >= b.MaxContextNotes {
			break
		}

		content, err := os.ReadFile(note.Path)
		if err != nil {
			continue
		}

		noteContent := string(content)
		if len(noteContent) > b.CharLimit {
			noteContent = noteContent[:b.CharLimit] + "\n\n... [truncated]"
		}

		noteTokens := EstimateTokens(noteContent)
		if totalTokens+noteTokens > b.MaxTokens {
			break
		}

		formatted := formatNoteForContext(note.NoteIndex, noteContent, note.Score)
		context.WriteString(formatted)
		context.WriteString("\n\n")

		included = append(included, note.Slug)
		totalTokens += noteTokens
	}

	return context.String(), included, nil
}

func formatNoteForContext(note NoteIndex, content string, score float64) string {
	tags := strings.Join(note.Tags, ", ")
	if tags == "" {
		tags = "none"
	}

	return fmt.Sprintf(`<note slug="%s" tags="%s">
<created>%s</created>
<updated>%s</updated>
<relevance>%.2f</relevance>
<content>
%s
</content>
</note>`, note.Slug, tags, note.CreatedAt.Format("2006-01-02"), note.UpdatedAt.Format("2006-01-02"), score, content)
}

var frontmatterRegex = regexp.MustCompile(`(?s)^---\n.*?\n---\n`)

func extractExcerpt(content string) string {
	fm := frontmatterRegex.FindStringIndex(content)
	if fm != nil {
		content = content[fm[1]:]
	}

	content = strings.TrimSpace(content)
	if len(content) > 500 {
		return content[:500]
	}
	return content
}

func getString(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return fallback
}

func getStringSlice(m map[string]interface{}, key string) []string {
	if v, ok := m[key].([]interface{}); ok {
		var result []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

func getTime(m map[string]interface{}, key string, fallback time.Time) time.Time {
	if v, ok := m[key].(string); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
	}
	return fallback
}

// EstimateTokens estimates token count from text content.
// Uses a weighted heuristic: 1 token ≈ 4 bytes for English,
// but counts runes when non-ASCII content is significant (>25%).
func EstimateTokens(text string) int {
	runes := []rune(text)
	// Fast path for pure ASCII content
	if len(text) == len(runes) {
		return len(text) / 4
	}

	// Count non-ASCII runes
	nonASCII := 0
	for _, r := range runes {
		if r > 127 {
			nonASCII++
		}
	}

	// If >25% non-ASCII, use rune-based estimate; otherwise use byte-based
	if nonASCII > len(runes)/4 {
		return len(runes) / 2
	}
	return len(text) / 4
}
