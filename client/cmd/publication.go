package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(PublicationCmd)
}

var PublicationCmd = &cobra.Command{
	Use:   "publication [command]",
	Short: "Publication commands",
}
