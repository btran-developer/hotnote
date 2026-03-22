package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
)

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
				md := goldmark.New()
		var buf []byte
		buf, err = md.Convert(content)
		if err != nil {
			fmt.Printf("Error rendering markdown: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println(string(buf))
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}