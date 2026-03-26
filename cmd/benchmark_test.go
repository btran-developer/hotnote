package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"hotnotego/internal/fsutil"
	slugifypkg "hotnotego/internal/slugify"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

// setupBenchmarkWorkspace creates a temporary workspace with n notes for benchmarking
func setupBenchmarkWorkspace(b *testing.B, n int) (workspacePath, configPath string, cleanup func()) {
	b.Helper()

	// Create temp directory for workspace
	tempDir, err := os.MkdirTemp("", "hotnote-bench-*")
	if err != nil {
		b.Fatalf("create temp dir: %v", err)
	}

	workspacePath = filepath.Join(tempDir, "workspace")
	// Config goes to tempdir/.config/hotnote/config.yaml
	configDir := filepath.Join(tempDir, ".config", "hotnote")
	configPath = filepath.Join(configDir, "config.yaml")

	// Create workspace directory
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		os.RemoveAll(tempDir)
		b.Fatalf("create workspace dir: %v", err)
	}

	// Create config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		b.Fatalf("create config dir: %v", err)
	}

	// Create config file
	configContent := fmt.Sprintf("current_workspace: default\nworkspaces:\n  default: %s\n", workspacePath)
	if err := fsutil.AtomicWrite(configPath, []byte(configContent), 0644); err != nil {
		os.RemoveAll(tempDir)
		b.Fatalf("write config: %v", err)
	}

	// Create n notes
	for i := 0; i < n; i++ {
		title := fmt.Sprintf("Note %d", i)
		slug := fmt.Sprintf("note-%d", i)
		noteID := uuid.New()
		createdAt := time.Now().UTC().Format(time.RFC3339)
		content := fmt.Sprintf(`---
id: %s
title: %s
created_at: %s
updated_at: %s
tags: []
---

# %s

This is the content of note %d.
`, noteID, title, createdAt, createdAt, title, i)

		notePath := filepath.Join(workspacePath, slug+".md")
		if err := fsutil.AtomicWrite(notePath, []byte(content), 0644); err != nil {
			os.RemoveAll(tempDir)
			b.Fatalf("create note %d: %v", i, err)
		}
	}

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return workspacePath, configPath, cleanup
}

// setBenchmarkConfig sets up environment for benchmark
func setBenchmarkConfig(b *testing.B, configPath string) func() {
	b.Helper()

	// Set HOME to config directory's parent so ~/.config/hotnote/config.yaml resolves correctly
	// configPath is: tempDir/.config/hotnote/config.yaml
	// So homeDir should be tempDir (parent of .config)
	configDir := filepath.Dir(configPath)            // tempDir/.config/hotnote
	homeDir := filepath.Dir(filepath.Dir(configDir)) // tempDir

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)

	return func() {
		os.Setenv("HOME", oldHome)
	}
}

// BenchmarkNewNote benchmarks creating a new note in a workspace with existing notes
func BenchmarkNewNote(b *testing.B) {
	_, configPath, cleanup := setupBenchmarkWorkspace(b, 500)
	defer cleanup()

	restoreConfig := setBenchmarkConfig(b, configPath)
	defer restoreConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a unique note title for each iteration
		title := fmt.Sprintf("Benchmark Note %d", i)

		// Simulate the new command
		wm, err := workspace.NewManager()
		if err != nil {
			b.Fatalf("create workspace manager: %v", err)
		}
		store := storage.NewStore(wm)

		slug := slugifypkg.Slugify(title)
		noteID := uuid.New()
		createdAt := time.Now().UTC().Format(time.RFC3339)
		content := fmt.Sprintf("---\nid: %s\ntitle: %s\ncreated_at: %s\nupdated_at: %s\ntags: []\n---\n\n# %s\n\n", noteID, title, createdAt, createdAt, title)

		// Use a unique slug to avoid conflicts
		uniqueSlug := fmt.Sprintf("%s-%d", slug, i)
		if err := store.Ensure(uniqueSlug, []byte(content)); err != nil {
			b.Fatalf("create note: %v", err)
		}
	}
}

// BenchmarkListNotes benchmarks listing notes with different note counts
func BenchmarkListNotes(b *testing.B) {
	noteCounts := []int{10, 100, 500}

	for _, count := range noteCounts {
		b.Run(fmt.Sprintf("notes_%d", count), func(b *testing.B) {
			_, configPath, cleanup := setupBenchmarkWorkspace(b, count)
			defer cleanup()

			restoreConfig := setBenchmarkConfig(b, configPath)
			defer restoreConfig()

			wm, err := workspace.NewManager()
			if err != nil {
				b.Fatalf("create workspace manager: %v", err)
			}

			_, wp, err := wm.Current()
			if err != nil {
				b.Fatalf("get current workspace: %v", err)
			}

			files, err := os.ReadDir(wp)
			if err != nil {
				b.Fatalf("read notes directory: %v", err)
			}

			for i := 0; i < 3; i++ {
				notes := make([]struct {
					slug    string
					modTime time.Time
				}, 0, len(files))

				for _, file := range files {
					if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
						info, err := file.Info()
						if err != nil {
							b.Fatalf("get file info: %v", err)
						}
						slug := file.Name()[:len(file.Name())-3]
						notes = append(notes, struct {
							slug    string
							modTime time.Time
						}{slug: slug, modTime: info.ModTime()})
					}
				}
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				notes := make([]struct {
					slug    string
					modTime time.Time
				}, 0, len(files))

				for _, file := range files {
					if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
						info, err := file.Info()
						if err != nil {
							b.Fatalf("get file info: %v", err)
						}
						slug := file.Name()[:len(file.Name())-3]
						notes = append(notes, struct {
							slug    string
							modTime time.Time
						}{slug: slug, modTime: info.ModTime()})
					}
				}
			}
		})
	}
}

// BenchmarkListNotesSorted benchmarks listing with different sort orders
func BenchmarkListNotesSorted(b *testing.B) {
	_, configPath, cleanup := setupBenchmarkWorkspace(b, 500)
	defer cleanup()

	restoreConfig := setBenchmarkConfig(b, configPath)
	defer restoreConfig()

	sortOrders := []string{"name", "updated", "created"}

	for _, sortOrder := range sortOrders {
		b.Run(fmt.Sprintf("sort_%s", sortOrder), func(b *testing.B) {
			wm, err := workspace.NewManager()
			if err != nil {
				b.Fatalf("create workspace manager: %v", err)
			}

			_, wp, err := wm.Current()
			if err != nil {
				b.Fatalf("get current workspace: %v", err)
			}

			files, err := os.ReadDir(wp)
			if err != nil {
				b.Fatalf("read notes directory: %v", err)
			}

			type noteData struct {
				slug    string
				modTime time.Time
				crTime  time.Time
			}
			var notes []noteData

			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
					info, err := file.Info()
					if err != nil {
						continue
					}
					slug := file.Name()[:len(file.Name())-3]
					notes = append(notes, noteData{
						slug:    slug,
						modTime: info.ModTime(),
						crTime:  info.ModTime(),
					})
				}
			}

			for i := 0; i < 3; i++ {
				switch sortOrder {
				case "updated":
					sort.Slice(notes, func(i, j int) bool {
						return notes[i].modTime.After(notes[j].modTime)
					})
				case "created":
					sort.Slice(notes, func(i, j int) bool {
						return notes[i].crTime.After(notes[j].crTime)
					})
				default:
					sort.Slice(notes, func(i, j int) bool {
						return notes[i].slug < notes[j].slug
					})
				}
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				switch sortOrder {
				case "updated":
					sort.Slice(notes, func(i, j int) bool {
						return notes[i].modTime.After(notes[j].modTime)
					})
				case "created":
					sort.Slice(notes, func(i, j int) bool {
						return notes[i].crTime.After(notes[j].crTime)
					})
				default:
					sort.Slice(notes, func(i, j int) bool {
						return notes[i].slug < notes[j].slug
					})
				}
			}
		})
	}
}

// BenchmarkOpenNote benchmarks opening a note (file lookup)
func BenchmarkOpenNote(b *testing.B) {
	_, configPath, cleanup := setupBenchmarkWorkspace(b, 500)
	defer cleanup()

	restoreConfig := setBenchmarkConfig(b, configPath)
	defer restoreConfig()

	slug := fmt.Sprintf("note-%d", 250)

	wm, err := workspace.NewManager()
	if err != nil {
		b.Fatalf("create workspace manager: %v", err)
	}

	store := storage.NewStore(wm)
	path, err := store.Path(slug)
	if err != nil {
		b.Fatalf("get note path: %v", err)
	}

	// Warmup
	for i := 0; i < 3; i++ {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			b.Fatalf("note not found: %s", path)
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			b.Fatalf("note not found: %s", path)
		}
	}
}

// BenchmarkRenderNote benchmarks rendering a note to HTML
func BenchmarkRenderNote(b *testing.B) {
	workspacePath, configPath, cleanup := setupBenchmarkWorkspace(b, 500)
	defer cleanup()

	restoreConfig := setBenchmarkConfig(b, configPath)
	defer restoreConfig()

	// Read a sample note for rendering
	samplePath := filepath.Join(workspacePath, "note-250.md")
	content, err := os.ReadFile(samplePath)
	if err != nil {
		b.Fatalf("read sample note: %v", err)
	}

	// Warmup
	for i := 0; i < 3; i++ {
		var buf bytes.Buffer
		if err := md.Convert(content, &buf); err != nil {
			b.Fatalf("render markdown: %v", err)
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err = md.Convert(content, &buf)
		if err != nil {
			b.Fatalf("render markdown: %v", err)
		}
	}
}

// BenchmarkWorkspaceManagerCreation benchmarks creating a new workspace manager
func BenchmarkWorkspaceManagerCreation(b *testing.B) {
	_, configPath, cleanup := setupBenchmarkWorkspace(b, 0)
	defer cleanup()

	restoreConfig := setBenchmarkConfig(b, configPath)
	defer restoreConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := workspace.NewManager()
		if err != nil {
			b.Fatalf("create workspace manager: %v", err)
		}
	}
}

// BenchmarkConfigLoading benchmarks loading the config file
func BenchmarkConfigLoading(b *testing.B) {
	_, configPath, cleanup := setupBenchmarkWorkspace(b, 0)
	defer cleanup()

	restoreConfig := setBenchmarkConfig(b, configPath)
	defer restoreConfig()

	// Create manager once to ensure config exists
	_, err := workspace.NewManager()
	if err != nil {
		b.Fatalf("create workspace manager: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a new manager which loads config
		wm, err := workspace.NewManager()
		if err != nil {
			b.Fatalf("create workspace manager: %v", err)
		}

		// Access current workspace to force config load
		_, _, err = wm.Current()
		if err != nil {
			b.Fatalf("get current workspace: %v", err)
		}
	}
}
