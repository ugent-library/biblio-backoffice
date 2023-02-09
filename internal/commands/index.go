package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func init() {
	removeOldIndexesDatasetCmd.PersistentFlags().IntVarP(
		&keepMaxIndexes,
		"keep",
		"k",
		0,
		"keep number of old indexes (default: 0)",
	)
	removeOldIndexesPublicationCmd.PersistentFlags().IntVarP(
		&keepMaxIndexes,
		"keep",
		"k",
		0,
		"keep number of old indexes (default: 0)",
	)
	indexDatasetCmd.AddCommand(indexDatasetCreateCmd)
	indexDatasetCmd.AddCommand(indexDatasetDeleteCmd)
	indexDatasetCmd.AddCommand(indexDatasetAllCmd)
	indexDatasetCmd.AddCommand(reindexDatasetCmd)
	indexDatasetCmd.AddCommand(initAliasDatasetCmd)
	indexDatasetCmd.AddCommand(removeOldIndexesDatasetCmd)
	indexDatasetCmd.AddCommand(listOldIndexesDatasetCmd)
	indexCmd.AddCommand(indexDatasetCmd)
	indexPublicationCmd.AddCommand(indexPublicationCreateCmd)
	indexPublicationCmd.AddCommand(indexPublicationDeleteCmd)
	indexPublicationCmd.AddCommand(indexPublicationAllCmd)
	indexPublicationCmd.AddCommand(reindexPublicationCmd)
	indexPublicationCmd.AddCommand(initAliasPublicationCmd)
	indexPublicationCmd.AddCommand(removeOldIndexesPublicationCmd)
	indexPublicationCmd.AddCommand(listOldIndexesPublicationCmd)
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
		ctx := context.Background()

		es := newDatasetSearchService()
		store := newRepository()

		bi, err := es.NewBulkIndexer(backends.BulkIndexerConfig{
			OnError: func(err error) {
				log.Printf("Indexing failed : %s", err)
			},
			OnIndexError: func(id string, err error) {
				log.Printf("Indexing failed for dataset [id: %s] : %s", id, err)
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		defer bi.Close(ctx)

		// send recs to indexer
		store.EachDataset(func(d *models.Dataset) bool {
			if err := bi.Index(ctx, d); err != nil {
				log.Printf("Indexing failed for dataset [id: %s] : %s", d.ID, err)
			}
			return true
		})
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
		ctx := context.Background()

		es := newPublicationSearchService()
		store := newRepository()

		bi, err := es.NewBulkIndexer(backends.BulkIndexerConfig{
			OnError: func(err error) {
				log.Printf("Indexing failed : %s", err)
			},
			OnIndexError: func(id string, err error) {
				log.Printf("Indexing failed for publication [id: %s] : %s", id, err)
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		defer bi.Close(ctx)

		// send recs to indexer
		store.EachPublication(func(p *models.Publication) bool {
			if err := bi.Index(ctx, p); err != nil {
				log.Printf("Indexing failed for publication [id: %s] : %s", p.ID, err)
			}
			return true
		})
	},
}

var initAliasPublicationCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize alias index publications",
	Run: func(cmd *cobra.Command, args []string) {
		reindexer := Services().PublicationSearchService.NewReindexer()
		if e := reindexer.InitAlias(); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}
	},
}

var initAliasDatasetCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize alias index datasets",
	Run: func(cmd *cobra.Command, args []string) {
		reindexer := Services().DatasetSearchService.NewReindexer()
		if e := reindexer.InitAlias(); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}
	},
}

var reindexPublicationCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex publications (and switch alias)",
	Run: func(cmd *cobra.Command, args []string) {

		searcher := Services().PublicationSearchService
		reindexer := searcher.NewReindexer()
		repo := Services().Repository

		startTime := time.Now()

		//reindex from source to index, and then set alias
		{
			indexC := make(chan *models.Publication)
			sql := "SELECT * FROM publications WHERE date_until IS NULL"

			var indexWG sync.WaitGroup
			indexWG.Add(1)

			go func() {
				defer indexWG.Done()
				reindexer.Reindex(indexC)
			}()

			added := 0
			err := repo.SelectPublications(
				sql,
				[]any{},
				func(publication *models.Publication) bool {
					indexC <- publication
					added++
					return true
				},
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			close(indexC)
			indexWG.Wait()

			fmt.Fprintf(
				os.Stderr,
				"added %d docs for sql '%s' to new index\n",
				added,
				sql,
			)

		}

		fmt.Fprintf(os.Stderr, "new index live now\n")

		endTime := time.Now()

		indexFunc := func(sql string, sqlArgs []any) int {
			added := 0

			indexC := make(chan *models.Publication)

			var indexWG sync.WaitGroup
			indexWG.Add(1)

			go func() {
				defer indexWG.Done()
				searcher.IndexMultiple(indexC)
			}()

			err := repo.SelectPublications(
				sql,
				sqlArgs,
				func(publication *models.Publication) bool {
					indexC <- publication
					added++
					return true
				},
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			close(indexC)
			indexWG.Wait()

			fmt.Fprintf(
				os.Stderr,
				"added %d docs for sql '%s' with args %v\n",
				added,
				sql,
				sqlArgs,
			)

			return added
		}

		/*
			index publications we've missed,
			but do not add documents that are now
			added to this live index (for which endTime is needed)
			as we may have outdated data (again)

			keep looping until no more additions are left
		*/
		const pgDateFormat = "2006-01-02 15:04:05-07"

		for {
			sql := `SELECT * FROM publications WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2`
			sqlArgs := []any{
				startTime.Format(pgDateFormat),
				endTime.Format(pgDateFormat),
			}

			startTime = time.Now()
			added := indexFunc(
				sql,
				sqlArgs,
			)
			endTime = time.Now()

			if added <= 0 {
				break
			}
		}

	},
}

var reindexDatasetCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex datasets (and switch alias)",
	Run: func(cmd *cobra.Command, args []string) {
		searcher := Services().DatasetSearchService
		reindexer := searcher.NewReindexer()
		repo := Services().Repository

		startTime := time.Now()

		//reindex from source to index, and then set alias
		{
			indexC := make(chan *models.Dataset)
			sql := "SELECT * FROM datasets WHERE date_until IS NULL"

			var indexWG sync.WaitGroup
			indexWG.Add(1)

			go func() {
				defer indexWG.Done()
				reindexer.Reindex(indexC)
			}()

			added := 0
			err := repo.SelectDatasets(
				sql,
				[]any{},
				func(dataset *models.Dataset) bool {
					indexC <- dataset
					added++
					return true
				},
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			close(indexC)
			indexWG.Wait()

			fmt.Fprintf(
				os.Stderr,
				"added %d docs for sql '%s' to new index\n",
				added,
				sql,
			)

		}

		fmt.Fprintf(os.Stderr, "new index live now\n")

		endTime := time.Now()

		indexFunc := func(sql string, sqlArgs []any) int {
			added := 0

			indexC := make(chan *models.Dataset)

			var indexWG sync.WaitGroup
			indexWG.Add(1)

			go func() {
				defer indexWG.Done()
				searcher.IndexMultiple(indexC)
			}()

			err := repo.SelectDatasets(
				sql,
				sqlArgs,
				func(dataset *models.Dataset) bool {
					indexC <- dataset
					added++
					return true
				},
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			close(indexC)
			indexWG.Wait()

			fmt.Fprintf(
				os.Stderr,
				"added %d docs for sql '%s' with args %v\n",
				added,
				sql,
				sqlArgs,
			)

			return added
		}

		/*
			index datasets we've missed,
			but do not add documents that are now
			added to this live index (for which endTime is needed)
			as we may have outdated data (again)

			keep looping until no more additions are left
		*/
		const pgDateFormat = "2006-01-02 15:04:05-07"

		for {
			sql := `SELECT * FROM datasets WHERE date_until IS NULL AND date_from >= $1 AND date_from <= $2`
			sqlArgs := []any{
				startTime.Format(pgDateFormat),
				endTime.Format(pgDateFormat),
			}

			startTime = time.Now()
			added := indexFunc(
				sql,
				sqlArgs,
			)
			endTime = time.Now()

			if added <= 0 {
				break
			}
		}
	},
}

var keepMaxIndexes int = 0
var removeOldIndexesDatasetCmd = &cobra.Command{
	Use:   "remove-old-indexes",
	Short: "Remove old dataset indexes",
	Run: func(cmd *cobra.Command, args []string) {
		reindexer := Services().DatasetSearchService.NewReindexer()
		reindexer.RemoveOldIndexes(keepMaxIndexes)
	},
}

var listOldIndexesDatasetCmd = &cobra.Command{
	Use:   "indexes",
	Short: "List dataset indexes",
	Run: func(cmd *cobra.Command, args []string) {
		reindexer := Services().DatasetSearchService.NewReindexer()
		indexes, e := reindexer.ListIndexes()
		if e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}
		for _, idx := range indexes {
			active := "inactive"
			if idx["active"] == "true" {
				active = "active"
			}
			fmt.Printf(
				"%s %s\n",
				active,
				idx["index"],
			)
		}
	},
}

var removeOldIndexesPublicationCmd = &cobra.Command{
	Use:   "remove-old-indexes",
	Short: "Removes old publication indexes",
	Run: func(cmd *cobra.Command, args []string) {
		reindexer := Services().PublicationSearchService.NewReindexer()
		reindexer.RemoveOldIndexes(keepMaxIndexes)
	},
}

var listOldIndexesPublicationCmd = &cobra.Command{
	Use:   "indexes",
	Short: "List publication indexes",
	Run: func(cmd *cobra.Command, args []string) {
		reindexer := Services().PublicationSearchService.NewReindexer()
		indexes, e := reindexer.ListIndexes()
		if e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}
		for _, idx := range indexes {
			active := "inactive"
			if idx["active"] == "true" {
				active = "active"
			}
			fmt.Printf(
				"%s %s\n",
				active,
				idx["index"],
			)
		}
	},
}
