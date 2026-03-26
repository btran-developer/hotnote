package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/workspace"
)

var (
	sortFlag string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Run: func(cmd *cobra.Command, args []string) {
		wm, err := workspace.NewManager()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("create workspace manager: %v", err))
			} else {
				fmt.Printf("create workspace manager: %v\n", err)
			}
			os.Exit(exitorrors.ExitConfigError)
		}

		_, workspacePath, err := wm.Current()
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("get current workspace: %v", err))
			} else {
				fmt.Printf("get current workspace: %v\n", err)
			}
			os.Exit(exitorrors.ExitConfigError)
		}

		files, err := os.ReadDir(workspacePath)
		if err != nil {
			if jsonFlag {
				outputJSONError(fmt.Sprintf("read notes directory: %v", err))
			} else {
				fmt.Printf("read notes directory: %v\n", err)
			}
			os.Exit(exitorrors.ExitConfigError)
		}

		// Collect note information
		type noteInfo struct {
			name    string
			slug    string
			modTime time.Time
			crTime  time.Time
			path    string
		}

		var notes []noteInfo
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
				slug := file.Name()[:len(file.Name())-3] // Remove .md extension
				info, err := file.Info()
				if err != nil {
					continue
				}

				// Get creation time from filesystem
				stat, err := os.Stat(filepath.Join(workspacePath, file.Name()))
				var crTime time.Time
				if err == nil {
					// Try to get creation time (BirthTime on some systems, CTime on others)
					// Use ModTime as fallback
					crTime = getCreationTime(stat)
				} else {
					crTime = info.ModTime()
				}

				notes = append(notes, noteInfo{
					name:    file.Name(),
					slug:    slug,
					modTime: info.ModTime(),
					crTime:  crTime,
					path:    filepath.Join(workspacePath, file.Name()),
				})
			}
		}

		// Sort notes based on sortFlag
		switch sortFlag {
		case "updated":
			// Sort by modification time, newest first
			sort.Slice(notes, func(i, j int) bool {
				return notes[i].modTime.After(notes[j].modTime)
			})
		case "created":
			// Sort by creation time, newest first
			sort.Slice(notes, func(i, j int) bool {
				return notes[i].crTime.After(notes[j].crTime)
			})
		default: // "name" or empty
			// Sort alphabetically by slug (default)
			sort.Slice(notes, func(i, j int) bool {
				return notes[i].slug < notes[j].slug
			})
		}

		if jsonFlag {
			var jsonNotes []map[string]string
			for _, note := range notes {
				jsonNotes = append(jsonNotes, map[string]string{
					"slug":       note.slug,
					"path":       note.path,
					"updated_at": note.modTime.UTC().Format(time.RFC3339),
				})
			}
			if err := outputJSON(jsonNotes); err != nil {
				outputJSONError(fmt.Sprintf("marshal JSON: %v", err))
			}
		} else {
			// Human-readable output with date in format: 2006-01-02 15:04
			for _, note := range notes {
				dateStr := note.modTime.Format("2006-01-02 15:04")
				fmt.Printf("%s\t%s\n", note.slug, dateStr)
			}
		}
	},
}

// getCreationTime attempts to get the file creation time
// Falls back to modification time if creation time is not available
func getCreationTime(stat os.FileInfo) time.Time {
	// On Unix systems, we can try to get birth time
	// This is system-dependent, so we fall back to ModTime
	// TODO: Use syscall.Stat_t on Unix to get birth time if available
	return stat.ModTime()
}

func init() {
	listCmd.Flags().StringVar(&sortFlag, "sort", "name", "Sort order: name (default), updated, or created")
	RootCmd.AddCommand(listCmd)
}
