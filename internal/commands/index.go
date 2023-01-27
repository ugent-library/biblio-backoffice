package commands

import (
	"log"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/models"
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
		es := newDatasetSearchService()
		store := newRepository()
		var indexWG sync.WaitGroup

		// indexing channel
		indexC := make(chan *models.Dataset)

		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			es.IndexMultiple(indexC)
		}()

		// send recs to indexer
		store.EachDataset(func(p *models.Dataset) bool {
			indexC <- p
			return true
		})

		close(indexC)

		// wait for indexing to finish
		indexWG.Wait()
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
	Short: "Reindex all publications",
	Run: func(cmd *cobra.Command, args []string) {
		es := newPublicationSearchService()
		store := newRepository()
		var indexWG sync.WaitGroup

		// indexing channel
		indexC := make(chan *models.Publication)

		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			es.IndexMultiple(indexC)
		}()

		// send recs to indexer
		store.EachPublication(func(p *models.Publication) bool {
			indexC <- p
			return true
		})

		close(indexC)

		// wait for indexing to finish
		indexWG.Wait()
	},
}
