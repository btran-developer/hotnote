package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"hotnotego/internal/storage"
	"os"
)

var md = goldmark.New()

var renderCmd = &cobra.Command{
	Use:   "render [title]",
	Short: "Render markdown to HTML",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		store := storage.NewStore(dataDir)
		path := store.Path(title)

		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading note: %v\n", err)
			os.Exit(1)
		}
		var buf bytes.Buffer
		err = md.Convert(content, &buf)
		if err != nil {
			fmt.Printf("Error rendering markdown: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(buf.String())
	},
}

func init() {
	RootCmd.AddCommand(renderCmd)
}
