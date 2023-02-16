package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(DatasetCmd)
}

var DatasetCmd = &cobra.Command{
	Use:   "dataset [command]",
	Short: "Dataset commands",
}
