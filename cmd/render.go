package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
	exitorrors "hotnotego/internal/errors"
	"hotnotego/internal/storage"
	"hotnotego/internal/workspace"
)

var (
	frontmatterFormat string
)

var md = goldmark.New(
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

var renderCmd = &cobra.Command{
	Use:     "render [title]",
	Short:   "Render markdown to HTML",
	Aliases: []string{"rdr"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		wm, err := workspace.NewManager()
		if err != nil {
			handleWorkspaceError(err)
		}

		store := storage.NewStore(wm)

		resolvedSlug, err := store.Resolve(title)
		if errors.Is(err, storage.ErrNoteNotFound) {
			if jsonFlag {
				outputJSONError(exitorrors.ErrNoteNotFound.Error())
			} else {
				fmt.Println(exitorrors.ErrNoteNotFound.Error())
			}
			os.Exit(exitorrors.ExitNotFound)
		}
		if errors.Is(err, storage.ErrMultipleMatches) {
			if jsonFlag {
				outputJSONError(exitorrors.ErrMultipleMatches.Error())
			} else {
				fmt.Println(exitorrors.ErrMultipleMatches.Error())
			}
			os.Exit(exitorrors.ExitInvalidInput)
		}
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrNoteResolve.Error())
			} else {
				fmt.Println(exitorrors.ErrNoteResolve.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		path, err := store.Path(resolvedSlug)
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrNotePath.Error())
			} else {
				fmt.Println(exitorrors.ErrNotePath.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			if jsonFlag {
				outputJSONError("note not found")
			} else {
				fmt.Println("note not found")
			}
			os.Exit(exitorrors.ExitNotFound)
		}

		// Parse and handle frontmatter
		markdownContent, frontmatterData := extractFrontmatter(string(content))

		var renderContent string
		if frontmatterFormat != "" && frontmatterData != nil {
			// Include frontmatter in the rendered output
			renderContent = formatFrontmatter(frontmatterData, frontmatterFormat) + "\n\n" + markdownContent
		} else {
			// Just render the markdown content (strip frontmatter)
			renderContent = markdownContent
		}

		var buf bytes.Buffer
		err = md.Convert([]byte(renderContent), &buf)
		if err != nil {
			if jsonFlag {
				outputJSONError(exitorrors.ErrNoteRender.Error())
			} else {
				fmt.Println(exitorrors.ErrNoteRender.Error())
			}
			os.Exit(exitorrors.ExitGeneral)
		}

		if jsonFlag {
			jsonResponse := map[string]string{"content": buf.String()}
			if err := outputJSON(jsonResponse); err != nil {
				outputJSONError(exitorrors.ErrMarshalJSON.Error())
				os.Exit(exitorrors.ExitGeneral)
			}
		} else {
			fmt.Println(buf.String())
		}
	},
}

// extractFrontmatter extracts YAML frontmatter from markdown content
// Returns the markdown content (without frontmatter) and parsed frontmatter data
func extractFrontmatter(content string) (string, map[string]interface{}) {
	lines := strings.Split(content, "\n")

	// Check if content starts with frontmatter delimiter
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		var frontmatterLines []string
		var markdownLines []string
		inFrontmatter := false
		frontmatterEnded := false

		for i, line := range lines {
			if i == 0 && strings.TrimSpace(line) == "---" {
				inFrontmatter = true
				continue
			}

			if inFrontmatter && !frontmatterEnded {
				if strings.TrimSpace(line) == "---" {
					frontmatterEnded = true
					inFrontmatter = false
				} else {
					frontmatterLines = append(frontmatterLines, line)
				}
			} else {
				markdownLines = append(markdownLines, line)
			}
		}

		if frontmatterEnded && len(frontmatterLines) > 0 {
			// Parse YAML frontmatter
			frontmatterYAML := strings.Join(frontmatterLines, "\n")
			var data map[string]interface{}
			if err := yaml.Unmarshal([]byte(frontmatterYAML), &data); err == nil {
				return strings.Join(markdownLines, "\n"), data
			}
		}
	}

	// No frontmatter found or parsing failed
	return content, nil
}

// formatFrontmatter formats frontmatter data as HTML based on the specified format
func formatFrontmatter(data map[string]interface{}, format string) string {
	switch format {
	case "table":
		return formatFrontmatterAsTable(data)
	case "dl":
		return formatFrontmatterAsDefinitionList(data)
	case "pre":
		fallthrough
	default:
		return formatFrontmatterAsPreformatted(data)
	}
}

// formatFrontmatterAsTable formats frontmatter as an HTML table
func formatFrontmatterAsTable(data map[string]interface{}) string {
	var buf strings.Builder
	buf.WriteString("<table class=\"frontmatter\">\n")
	buf.WriteString("  <thead>\n")
	buf.WriteString("    <tr><th>Key</th><th>Value</th></tr>\n")
	buf.WriteString("  </thead>\n")
	buf.WriteString("  <tbody>\n")

	for key, value := range data {
		buf.WriteString(fmt.Sprintf("    <tr><td>%s</td><td>%v</td></tr>\n", key, value))
	}

	buf.WriteString("  </tbody>\n")
	buf.WriteString("</table>\n")
	return buf.String()
}

// formatFrontmatterAsDefinitionList formats frontmatter as an HTML definition list
func formatFrontmatterAsDefinitionList(data map[string]interface{}) string {
	var buf strings.Builder
	buf.WriteString("<dl class=\"frontmatter\">\n")

	for key, value := range data {
		buf.WriteString(fmt.Sprintf("  <dt>%s</dt>\n", key))
		buf.WriteString(fmt.Sprintf("  <dd>%v</dd>\n", value))
	}

	buf.WriteString("</dl>\n")
	return buf.String()
}

// formatFrontmatterAsPreformatted formats frontmatter as preformatted text
func formatFrontmatterAsPreformatted(data map[string]interface{}) string {
	var buf strings.Builder
	buf.WriteString("```yaml\n")

	for key, value := range data {
		buf.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}

	buf.WriteString("```\n")
	return buf.String()
}

func init() {
	renderCmd.Flags().StringVar(&frontmatterFormat, "frontmatter-format", "", "Frontmatter format: table, dl, or pre (required when frontmatter desired)")
	RootCmd.AddCommand(renderCmd)
}
