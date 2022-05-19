package commands

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	indexDatasetCmd.AddCommand(indexDatasetCreateCmd)
	indexDatasetCmd.AddCommand(indexDatasetDeleteCmd)
	indexDatasetCmd.AddCommand(indexDatasetAllCmd)
	indexCmd.AddCommand(indexDatasetCmd)
	indexPublicationCmd.AddCommand(indexPublicationCreateCmd)
	indexPublicationCmd.AddCommand(indexPublicationDeleteCmd)
	indexPublicationCmd.AddCommand(indexPublicationAllCmd)
	indexCmd.AddCommand(indexPublicationCmd)
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
		if err := newDatasetSearchService().CreateIndex(); err != nil {
			log.Fatal(err)
		}
	},
}

var indexDatasetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete dataset index",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newDatasetSearchService().DeleteIndex(); err != nil {
			log.Fatal(err)
		}
	},
}

var indexDatasetAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Reindex all datasets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Engine().IndexAllDatasets(); err != nil {
			log.Fatal(err)
		}
	},
}

var indexPublicationCmd = &cobra.Command{
	Use:   "publication [command]",
	Short: "Publication index commands",
}

var indexPublicationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create publication index",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newPublicationSearchService().CreateIndex(); err != nil {
			log.Fatal(err)
		}
	},
}

var indexPublicationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete publication index",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newPublicationSearchService().DeleteIndex(); err != nil {
			log.Fatal(err)
		}
	},
}

var indexPublicationAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Reindex all datasets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Engine().IndexAllPublications(); err != nil {
			log.Fatal(err)
		}
	},
}
