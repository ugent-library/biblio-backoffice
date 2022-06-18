package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services"
	"github.com/ugent-library/biblio-backend/internal/services/webapp"
)

func init() {
	serverCmd.PersistentFlags().String("base-url", "", "base url")

	serverStartCmd.Flags().String("mode", defaultMode, "server mode (development, staging or production)")
	serverStartCmd.Flags().String("host", defaultHost, "server host")
	serverStartCmd.Flags().Int("port", defaultPort, "server port")
	serverStartCmd.Flags().String("session-name", defaultSessionName, "session name")
	serverStartCmd.Flags().String("session-secret", "", "session secret")
	serverStartCmd.Flags().Int("session-max-age", defaultSessionMaxAge, "session lifetime")
	serverStartCmd.Flags().String("csrf-name", "", "csrf cookie name")
	serverStartCmd.Flags().String("csrf-secret", "", "csrf cookie secret")

	webapp.AddCommands(serverCmd, Services())
	serverCmd.AddCommand(serverStartCmd)
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server [command]",
	Short: "The biblio-backend HTTP server",
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the http server",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()
		e.Repository.AddPublicationListener(func(p *models.Publication) {
			if err := e.PublicationSearchService.Index(p); err != nil {
				log.Println(fmt.Errorf("error indexing publication %s: %w", p.ID, err))
			}
		})
		e.Repository.AddDatasetListener(func(d *models.Dataset) {
			if err := e.DatasetSearchService.Index(d); err != nil {
				log.Println(fmt.Errorf("error indexing dataset %s: %w", d.ID, err))
			}
		})

		wa, err := webapp.New(e)
		if err != nil {
			log.Fatal(err)
		}

		if err = services.Start(wa); err != nil {
			log.Fatal(err)
		}
	},
}
