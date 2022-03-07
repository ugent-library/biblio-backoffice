package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/services"
	"github.com/ugent-library/biblio-backend/services/webapp"
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

	webapp.AddCommands(serverCmd, Engine())
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
		wa, err := webapp.New(Engine())
		if err != nil {
			log.Fatal(err)
		}

		if err = services.Start(wa); err != nil {
			log.Fatal(err)
		}
	},
}
