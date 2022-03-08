package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/workers/orcid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func init() {
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
		c, err := client.NewClient(client.Options{
			HostPort: client.DefaultHostPort,
		})
		if err != nil {
			log.Fatalln("Unable to create client", err)
		}
		defer c.Close()

		w := worker.New(c, "orcid", worker.Options{})

		w.RegisterWorkflowWithOptions(orcid.AddPublicationsWorkflow(e), workflow.RegisterOptions{
			Name: "AddPublicationsToORCIDWorkflow",
		})
		w.RegisterActivityWithOptions(orcid.AddPublications(e, e.ORCIDSandbox), activity.RegisterOptions{
			Name: "AddPublicationsToORCID",
		})

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	},
}
