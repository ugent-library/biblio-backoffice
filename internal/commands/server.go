package commands

import (
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/routes"
	"github.com/ugent-library/go-graceful/server"
)

func init() {
	serverStartCmd.Flags().String("base-url", defaultBaseURL, "base url")
	serverStartCmd.Flags().Int("port", defaultPort, "server port")

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
		r := chi.NewRouter()
		routes.Register(r)
		s := server.New(r,
			server.WithPort(viper.GetInt("port")),
		)
		s.Start()
	},
}
