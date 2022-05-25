package commands

import (
	"log"

	"github.com/spf13/cobra"
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

		store := newStore()
		if err := store.PurgeAllPublications(); err != nil {
			log.Fatal(err)
		}
		if err := store.PurgeAllDatasets(); err != nil {
			log.Fatal(err)
		}

		if err := newPublicationSearchService().DeleteIndex(); err != nil {
			log.Fatal(err)
		}
		if err := newDatasetSearchService().DeleteIndex(); err != nil {
			log.Fatal(err)
		}
		if err := newPublicationSearchService().CreateIndex(); err != nil {
			log.Fatal(err)
		}
		if err := newDatasetSearchService().CreateIndex(); err != nil {
			log.Fatal(err)
		}

		if err := newFileStore().PurgeAll(); err != nil {
			log.Fatal(err)
		}
	},
}
