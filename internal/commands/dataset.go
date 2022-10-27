package commands

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/ulid"
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
		e := Services()

		var indexWG sync.WaitGroup

		// indexing channel
		indexC := make(chan *models.Dataset)

		// start bulk indexer
		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			e.DatasetSearchService.IndexMultiple(indexC)
		}()

		dec := json.NewDecoder(os.Stdin)

		lineNo := 0
		for {
			lineNo += 1
			d := models.Dataset{
				ID:     ulid.MustGenerate(),
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
			if err := e.Repository.SaveDataset(&d, nil); err != nil {
				log.Fatalf("Unable to store dataset from line %d : %v", lineNo, err)
			}

			indexC <- &d
		}

		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
	},
}

var datasetImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import datasets",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()

		var indexWG sync.WaitGroup

		// indexing channel
		indexC := make(chan *models.Dataset)

		// start bulk indexer
		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			e.DatasetSearchService.IndexMultiple(indexC)
		}()

		dec := json.NewDecoder(os.Stdin)

		lineNo := 0
		for {
			lineNo += 1
			d := models.Dataset{}
			if err := dec.Decode(&d); errors.Is(err, io.EOF) {
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
			if err := e.Repository.ImportCurrentDataset(&d); err != nil {
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

			indexC <- &d
		}

		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
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
