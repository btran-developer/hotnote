package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "note-cli",
	Short: "A brief description of your application",
	Long:  `A more detailed description of your application's purpose and functionality`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hi!")
	},
}