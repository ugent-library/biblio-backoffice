package commands

import (
	"html/template"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/controllers"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/helpers/assets"
	"github.com/ugent-library/biblio-backend/internal/routes"
	"github.com/ugent-library/go-graceful/server"
	"github.com/unrolled/render"
)

func init() {
	serverStartCmd.Flags().String("base-url", defaultBaseURL, "base url")
	serverStartCmd.Flags().String("host", defaultHost, "server host")
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
		// engine
		e, err := engine.New(engine.Config{
			URL:      viper.GetString("librecat-url"),
			Username: viper.GetString("librecat-username"),
			Password: viper.GetString("librecat-password"),
		})
		if err != nil {
			log.Fatal(err)
		}

		// renderer
		r := render.New(render.Options{
			Directory:                   "templates",
			Extensions:                  []string{".gohtml"},
			Layout:                      "layout",
			RenderPartialsWithoutPrefix: true,
			Funcs: []template.FuncMap{
				assets.FuncMap(),
			},
		})

		// router
		router := chi.NewRouter()
		routes.Register(
			router,
			controllers.NewPublication(e, r),
		)

		// server
		s := server.New(router,
			server.WithHost(viper.GetString("host")),
			server.WithPort(viper.GetInt("port")),
		)
		s.Start()
	},
}
