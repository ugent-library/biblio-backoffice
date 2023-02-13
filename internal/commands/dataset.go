package commands

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func init() {
	datasetCmd.AddCommand(datasetGetCmd)
	datasetCmd.AddCommand(datasetAllCmd)
	datasetCmd.AddCommand(datasetAddCmd)
	datasetCmd.AddCommand(datasetImportCmd)
	datasetCmd.AddCommand(oldDatasetImportCmd)
	datasetCmd.AddCommand(datasetReindexCmd)
	rootCmd.AddCommand(datasetCmd)
}

var datasetCmd = &cobra.Command{
	Use:   "dataset [command]",
	Short: "Dataset commands",
}

var datasetGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get datasets by id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := newRepository()
		e := json.NewEncoder(os.Stdout)
		for _, id := range args {
			d, err := s.GetDataset(id)
			if err != nil {
				log.Fatal(err)
			}
			e.Encode(d)
		}
	},
}

var datasetAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all datasets",
	Run: func(cmd *cobra.Command, args []string) {
		s := newRepository()
		e := json.NewEncoder(os.Stdout)
		s.EachDataset(func(d *models.Dataset) bool {
			e.Encode(d)
			return true
		})
	},
}

var datasetAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add datasets",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		e := Services()

		dec := json.NewDecoder(os.Stdin)

		bi, err := e.DatasetSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
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

		lineNo := 0
		for {
			lineNo += 1
			d := &models.Dataset{
				ID:     ulid.Make().String(),
				Status: "private",
			}
			if err := dec.Decode(&d); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				log.Fatalf("Unable to decode dataset at line %d : %v", lineNo, err)
			}
			if err := d.Validate(); err != nil {
				log.Printf("Validation failed for dataset at line %d : %v", lineNo, err)
				continue
			}
			if err := e.Repository.SaveDataset(d, nil); err != nil {
				log.Fatalf("Unable to store dataset from line %d : %v", lineNo, err)
			}

			if err := bi.Index(ctx, d); err != nil {
				log.Printf("Indexing failed for dataset [id: %s] at line %d : %s", d.ID, lineNo, err)
			}
		}
	},
}

var datasetImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import datasets",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		e := Services()

		dec := json.NewDecoder(os.Stdin)

		bi, err := e.DatasetSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
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

		lineNo := 0
		for {
			lineNo += 1
			d := &models.Dataset{}
			if err := dec.Decode(d); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				log.Fatalf("Unable to decode dataset at line %d : %v", lineNo, err)
			}
			if err := d.Validate(); err != nil {
				log.Printf(
					"Validation failed for dataset[snapshot_id: %s, id: %s] at line %d : %v",
					d.SnapshotID,
					d.ID,
					lineNo,
					err,
				)
				continue
			}
			if err := e.Repository.ImportCurrentDataset(d); err != nil {
				log.Printf(
					"Unable to store dataset[snapshot_id: %s, id: %s] from line %d : %v",
					d.SnapshotID,
					d.ID,
					lineNo,
					err,
				)
				continue
			}
			log.Printf(
				"Added dataset[snapshot_id: %s, id: %s]",
				d.SnapshotID,
				d.ID,
			)

			if err := bi.Index(ctx, d); err != nil {
				log.Printf("Indexing failed for dataset [id: %s] : %s", d.ID, err)
			}
		}
	},
}

var oldDatasetImportCmd = &cobra.Command{
	Use:   "import-version",
	Short: "Import old datasets",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()

		dec := json.NewDecoder(os.Stdin)

		lineNo := 0
		for {
			lineNo += 1
			d := models.Dataset{}
			if err := dec.Decode(&d); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				log.Fatalf("Unable to decode old dataset at line %d : %v", lineNo, err)
			}
			if err := d.Validate(); err != nil {
				log.Printf("Validation failed for old dataset[snapshot_id: %s, id: %s] at line %d : %v",
					d.SnapshotID,
					d.ID,
					lineNo,
					err,
				)
				continue
			}
			if err := e.Repository.ImportOldDataset(&d); err != nil {
				log.Printf(
					"Unable to store old dataset[snapshot_id: %s, id: %s] from line %d : %v",
					d.SnapshotID,
					d.ID,
					lineNo,
					err,
				)
				continue
			}
			log.Printf(
				"Added old dataset[snapshot_id: %s, id: %s]",
				d.SnapshotID,
				d.ID,
			)
		}
	},
}

var datasetReindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex into a new search index",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		services := Services()

		startTime := time.Now()

		indexed := 0

		log.Println("Indexing to new index...")

		switcher, err := services.DatasetSearchService.NewIndexSwitcher(backends.BulkIndexerConfig{
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
		services.Repository.EachDataset(func(p *models.Dataset) bool {
			if err := switcher.Index(ctx, p); err != nil {
				log.Printf("Indexing failed for dataset [id: %s] : %s", p.ID, err)
			}
			indexed++
			return true
		})

		log.Printf("Indexed %d datasets...", indexed)

		log.Println("Switching to new index...")

		if err := switcher.Switch(ctx); err != nil {
			log.Fatal(err)
		}

		endTime := time.Now()

		log.Println("Indexing changes since start of reindex...")

		for {
			indexed = 0

			bi, err := services.DatasetSearchService.NewBulkIndexer(backends.BulkIndexerConfig{
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

			err = services.Repository.DatasetsBetween(startTime, endTime, func(p *models.Dataset) bool {
				if err := bi.Index(ctx, p); err != nil {
					log.Printf("Indexing failed for dataset [id: %s] : %s", p.ID, err)
				}
				indexed++
				return true
			})
			if err != nil {
				log.Fatal(err)
			}

			if err = bi.Close(ctx); err != nil {
				log.Fatal(err)
			}

			if indexed == 0 {
				break
			}

			log.Printf("Indexed %d datasets...", indexed)

			startTime = endTime
			endTime = time.Now()
		}

		log.Println("Done.")
	},
}
