package cli

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/backends"
)

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().Bool("force", false, "force destructive reset of all data")
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Destructive reset",
	RunE: func(cmd *cobra.Command, args []string) error {
		if force, _ := cmd.Flags().GetBool("force"); !force {
			cmd.Println("The --force flag is required to perform a destructive reset.")
			return nil
		}

		ctx := context.Background()

		services := newServices()

		if err := services.Repo.PurgeAllPublications(); err != nil {
			return err
		}
		if err := services.Repo.PurgeAllDatasets(); err != nil {
			return err
		}

		publicationSwitcher, err := services.SearchService.NewPublicationIndexSwitcher(backends.BulkIndexerConfig{
			OnError: func(err error) {
				// TODO
				log.Fatal(err)
			},
			OnIndexError: func(id string, err error) {
				// TODO
				log.Fatal(err)
			},
		})
		if err != nil {
			return err
		}
		if err := publicationSwitcher.Switch(ctx); err != nil {
			return err
		}

		datasetSwitcher, err := services.SearchService.NewDatasetIndexSwitcher(backends.BulkIndexerConfig{
			OnError: func(err error) {
				// TODO
				log.Fatal(err)
			},
			OnIndexError: func(id string, err error) {
				// TODO
				log.Fatal(err)
			},
		})
		if err != nil {
			return err
		}
		if err := datasetSwitcher.Switch(ctx); err != nil {
			return err
		}

		if err := services.FileStore.DeleteAll(ctx); err != nil {
			return err
		}

		return nil
	},
}
