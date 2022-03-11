package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/workers/dataset"
	"github.com/ugent-library/biblio-backend/internal/workers/orcid"
	"go.temporal.io/sdk/worker"
)

func init() {
	startWorkerCmd.AddCommand(startStoreDatasetWorkerCmd)
	startWorkerCmd.AddCommand(startORCIDWorkerCmd)
	workerCmd.AddCommand(startWorkerCmd)
	rootCmd.AddCommand(workerCmd)
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
			PublicationService:       e.PublicationService,
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
			DatasetService: e.DatasetService,
		}

		w := worker.New(e.Temporal, "store-dataset", worker.Options{})
		w.RegisterWorkflow(dataset.StoreDatasetWorkflow)
		w.RegisterActivity(a.StoreDatasetInRepository)

		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	},
}
