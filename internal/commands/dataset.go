package commands

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

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
