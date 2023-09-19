package commands

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
)

func init() {
	resetCmd.Flags().Bool("confirm", false, "destructive reset of all data")
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Destructive reset",
	Run: func(cmd *cobra.Command, args []string) {
		if confirm, _ := cmd.Flags().GetBool("confirm"); !confirm {
			return
		}

		ctx := context.Background()

		services := Services()

		if err := services.Repo.PurgeAllPublications(); err != nil {
			log.Fatal(err)
		}
		if err := services.Repo.PurgeAllDatasets(); err != nil {
			log.Fatal(err)
		}

		publicationSwitcher, err := services.SearchService.NewPublicationIndexSwitcher(backends.BulkIndexerConfig{
			OnError: func(err error) {
				log.Fatal(err)
			},
			OnIndexError: func(id string, err error) {
				log.Fatal(err)
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		if err := publicationSwitcher.Switch(ctx); err != nil {
			log.Fatal(err)
		}

		datasetSwitcher, err := services.SearchService.NewDatasetIndexSwitcher(backends.BulkIndexerConfig{
			OnError: func(err error) {
				log.Fatal(err)
			},
			OnIndexError: func(id string, err error) {
				log.Fatal(err)
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		if err := datasetSwitcher.Switch(ctx); err != nil {
			log.Fatal(err)
		}

		if err := services.FileStore.DeleteAll(ctx); err != nil {
			log.Fatal(err)
		}
	},
}
