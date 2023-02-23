package commands

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/server"
)

func init() {
	apiStartCmd.Flags().String("host", defaultAPIHost, "api server host")
	apiStartCmd.Flags().Int("port", defaultAPIPort, "api server port")
	apiStartCmd.Flags().String("username", "", "api server administrator username")
	apiStartCmd.Flags().String("password", "", "api server administrator password")
	apiStartCmd.Flags().Bool("api-tls-enabled", false, "api server enable TLS encryped connections")
	apiStartCmd.Flags().String("api-tls-servercert", "", "api server location of server certificate file")
	apiStartCmd.Flags().String("api-tls-serverkey", "", "api server location of server certificate key file")

	apiCmd.AddCommand(apiStartCmd)
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api [command]",
	Short: "api commands",
}

var apiStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the api server",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup logger
		logger := newLogger()

		// Setup services
		e := Services()

		e.Repository.AddPublicationListener(func(p *models.Publication) {
			if p.DateUntil == nil {
				if err := e.PublicationSearchService.Index(p); err != nil {
					logger.Errorf("error indexing publication %s: %w", p.ID, err)
				}
			}
		})
		e.Repository.AddDatasetListener(func(d *models.Dataset) {
			if d.DateUntil == nil {
				if err := e.DatasetSearchService.Index(d); err != nil {
					logger.Errorf("error indexing dataset %s: %w", d.ID, err)
				}
			}
		})

		srv := server.New(e, logger)
		addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
		log.Printf("Listening at %s", addr)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		if err := srv.Serve(listener); err != nil {
			log.Fatal(err)
		}
	},
}
