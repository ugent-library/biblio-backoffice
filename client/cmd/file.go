package cmd

import "github.com/spf13/cobra"

// Set file buffer size
var fileBufSize = 524288

func init() {
	rootCmd.AddCommand(FileCmd)
}

var FileCmd = &cobra.Command{
	Use:   "file [command]",
	Short: "File commands",
}
