package commands

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func init() {
	reindexPublicationCmd.PersistentFlags().BoolVar(&seedIndex, "seed", false, "fully seed from database")
	reindexDatasetCmd.PersistentFlags().BoolVar(&seedIndex, "seed", false, "fully seed from database")
	indexDatasetCmd.AddCommand(indexDatasetCreateCmd)
	indexDatasetCmd.AddCommand(indexDatasetDeleteCmd)
	indexDatasetCmd.AddCommand(indexDatasetAllCmd)
	indexDatasetCmd.AddCommand(reindexDatasetCmd)
	indexDatasetCmd.AddCommand(initAliasDatasetCmd)
	indexCmd.AddCommand(indexDatasetCmd)
	indexPublicationCmd.AddCommand(indexPublicationCreateCmd)
	indexPublicationCmd.AddCommand(indexPublicationDeleteCmd)
	indexPublicationCmd.AddCommand(indexPublicationAllCmd)
	indexPublicationCmd.AddCommand(reindexPublicationCmd)
	indexPublicationCmd.AddCommand(initAliasPublicationCmd)
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

var initAliasPublicationCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize alias index publications",
	Run: func(cmd *cobra.Command, args []string) {
		search := Services().PublicationSearchService
		if e := search.Init(); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}
	},
}

var initAliasDatasetCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize alias index datasets",
	Run: func(cmd *cobra.Command, args []string) {
		search := Services().DatasetSearchService
		if e := search.Init(); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}
	},
}

/*
	--seed: feed entire table back into elasticsearch
			instead of relying on /_reindex and recent
			updates only
*/
var seedIndex bool = false
var reindexPublicationCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex publications (and switch alias)",
	Run: func(cmd *cobra.Command, args []string) {

		var startTime time.Time
		var endTime time.Time
		search := Services().PublicationSearchService
		repo := Services().Repository

		if seedIndex {
			startTime = time.Time{}
		} else {
			hits, hitsErr := search.Search(&models.SearchArgs{
				Query:    "",
				Filters:  map[string][]string{},
				PageSize: 1,
				Page:     1,
				Sort:     []string{"date_from:desc"},
			})
			if hitsErr != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", hitsErr.Error())
				os.Exit(1)
			}
			if len(hits.Hits) == 0 {
				startTime = time.Time{}
			} else {
				startTime = (*hits.Hits[0].DateFrom).Add(-time.Hour * 1)
			}
		}

		//copy entire index to new one and set alias
		if e := search.Reindex(); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}

		endTime = time.Now()

		/*
			index publications we've missed,
			but do not add documents that are now
			added to this live index (for which endTime is needed)
			as we may have outdated data (again)
		*/
		var indexWG sync.WaitGroup
		indexC := make(chan *models.Publication)

		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			search.IndexMultiple(indexC)
		}()

		pgDateFormat := "2006-01-02 15:04:05-07"
		sql := `SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2`
		dateFrom := startTime.Format(pgDateFormat)
		dateTo := endTime.Format(pgDateFormat)
		var countAdded int = 0
		err := repo.SelectPublications(
			sql,
			[]any{dateFrom, dateTo},
			func(publication *models.Publication) bool {
				indexC <- publication
				countAdded++
				return true
			},
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		close(indexC)
		indexWG.Wait()
	},
}

var reindexDatasetCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex datasets (and switch alias)",
	Run: func(cmd *cobra.Command, args []string) {

		search := Services().DatasetSearchService
		repo := Services().Repository

		var startTime time.Time
		var endTime time.Time

		if seedIndex {
			startTime = time.Time{}
		} else {
			hits, hitsErr := search.Search(&models.SearchArgs{
				Query:    "",
				Filters:  map[string][]string{},
				PageSize: 1,
				Page:     1,
				Sort:     []string{"date_from:desc"},
			})
			if hitsErr != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", hitsErr.Error())
				os.Exit(1)
			}
			if len(hits.Hits) == 0 {
				startTime = time.Time{}
			} else {
				startTime = (*hits.Hits[0].DateFrom).Add(-time.Hour * 1)
			}
		}

		//copy entire index to new one and set alias
		if e := search.Reindex(); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}

		endTime = time.Now()

		//index publications we've missed
		var indexWG sync.WaitGroup
		indexC := make(chan *models.Dataset)

		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			search.IndexMultiple(indexC)
		}()

		pgDateFormat := "2006-01-02 15:04:05-07"
		sql := `SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2`
		dateFrom := startTime.Format(pgDateFormat)
		dateTo := endTime.Format(pgDateFormat)
		err := repo.SelectDatasets(
			sql,
			[]any{dateFrom, dateTo},
			func(dataset *models.Dataset) bool {
				indexC <- dataset
				return true
			},
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		close(indexC)
		indexWG.Wait()
	},
}
