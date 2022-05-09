package commands

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/workers/dataset"
	"github.com/ugent-library/biblio-backend/internal/workers/orcid"
	"go.temporal.io/sdk/worker"
)

func init() {
	startWorkerCmd.AddCommand(startStoreDatasetWorkerCmd)
	startWorkerCmd.AddCommand(startORCIDWorkerCmd)
	workerCmd.AddCommand(startWorkerCmd)
	rootCmd.AddCommand(workerCmd)
	rootCmd.AddCommand(snapstoreCmd)
}

var snapstoreCmd = &cobra.Command{
	Use:   "snapstore",
	Short: "start biblio-backend orcid Temporal worker",
	Run: func(cmd *cobra.Command, args []string) {
		dsn := "postgres://nsteenla:@localhost:5432/biblio_snapstore?sslmode=disable"
		db, err := pgxpool.Connect(context.Background(), dsn)
		if err != nil {
			log.Fatal(err)
		}
		c := snapstore.New(db)
		publicationStore := c.Store("publication")
		datasetStore := c.Store("dataset")

		pID := "publication-1"
		dID := "dataset-1"
		p := &models.Publication{Title: "Test publication"}
		d := &models.Dataset{Title: "Test dataset"}

		err = c.Transaction(context.Background(), func(o snapstore.Options) error {
			if err := publicationStore.AddVersion("nsteenla", pID, p, o); err != nil {
				return err
			}
			if err := datasetStore.AddVersion("nsteenla", dID, d, o); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		if err := publicationStore.AddSnapshot("nsteenla", pID, snapstore.StrategyMine, snapstore.Options{}); err != nil {
			log.Fatal(err)
		}
	},
}

var workerCmd = &cobra.Command{
	Use:   "worker [command]",
	Short: "biblio-backend Temporal workers",
}

var startWorkerCmd = &cobra.Command{
	Use:   "start [command]",
	Short: "start biblio-backend Temporal worker",
}

var startORCIDWorkerCmd = &cobra.Command{
	Use:   "orcid",
	Short: "start biblio-backend orcid Temporal worker",
	Run: func(cmd *cobra.Command, args []string) {
		e := Engine()
		defer e.Temporal.Close()

		a := &orcid.Activities{
			Store:                    e.Store,
			PublicationSearchService: e.PublicationSearchService,
			OrcidSandbox:             e.ORCIDSandbox,
		}

		w := worker.New(e.Temporal, "orcid", worker.Options{})
		w.RegisterWorkflow(orcid.SendPublicationsToORCIDWorkflow)
		w.RegisterActivity(a.SendPublicationsToORCID)

		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	},
}

var startStoreDatasetWorkerCmd = &cobra.Command{
	Use:   "store-dataset",
	Short: "start biblio-backend store-dataset Temporal worker",
	Run: func(cmd *cobra.Command, args []string) {
		e := Engine()
		defer e.Temporal.Close()

		a := &dataset.Activities{
			DatasetService: e.Store,
		}

		w := worker.New(e.Temporal, "store-dataset", worker.Options{})
		w.RegisterWorkflow(dataset.StoreDatasetWorkflow)
		w.RegisterActivity(a.StoreDatasetInRepository)

		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	},
}
