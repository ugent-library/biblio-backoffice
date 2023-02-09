package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(FileCmd)
}

var FileCmd = &cobra.Command{
	Use:   "file [command]",
	Short: "File commands",
}
