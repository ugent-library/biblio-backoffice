package commands

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	indexDatasetCmd.AddCommand(indexDatasetCreateCmd)
	indexDatasetCmd.AddCommand(indexDatasetDeleteCmd)
	indexCmd.AddCommand(indexDatasetCmd)
	rootCmd.AddCommand(indexCmd)
}

var indexCmd = &cobra.Command{
	Use:   "index [command]",
	Short: "Index commands",
}

var indexDatasetCmd = &cobra.Command{
	Use:   "dataset [command]",
	Short: "Dataset index commands",
}

var indexDatasetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create dataset index",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newEs6Client().CreateDatasetIndex(); err != nil {
			log.Fatal(err)
		}
	},
}

var indexDatasetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete dataset index",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newEs6Client().DeleteDatasetIndex(); err != nil {
			log.Fatal(err)
		}
	},
}
