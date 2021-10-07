package commands

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/helpers"
	"github.com/ugent-library/biblio-backend/internal/routes"
	"github.com/ugent-library/go-graceful/server"
	"github.com/ugent-library/go-oidc/oidc"

	// "github.com/ugent-library/go-oidc/oidc"
	"github.com/ugent-library/go-web/mix"
	"github.com/ugent-library/go-web/urls"
	"github.com/unrolled/render"
)

func init() {
	serverStartCmd.Flags().String("base-url", "", "base url")
	serverStartCmd.Flags().String("host", defaultHost, "server host")
	serverStartCmd.Flags().Int("port", defaultPort, "server port")
	serverStartCmd.Flags().String("session-name", defaultSessionName, "session name")
	serverStartCmd.Flags().String("session-secret", "", "session secret")
	serverStartCmd.Flags().Int("session-max-age", defaultSessionMaxAge, "session lifetime")

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
		host := viper.GetString("host")
		port := viper.GetInt("port")
		baseURL := viper.GetString("base-url")

		if baseURL == "" {
			if host == "" {
				baseURL = "http://localhost"
			} else {
				baseURL = "http://" + host
			}
			if port != 80 {
				baseURL = fmt.Sprintf("%s:%d", baseURL, port)
			}
		}

		// engine
		e, err := engine.New(engine.Config{
			URL:      viper.GetString("librecat-url"),
			Username: viper.GetString("librecat-username"),
			Password: viper.GetString("librecat-password"),
		})
		if err != nil {
			log.Fatal(err)
		}

		// router
		router := mux.NewRouter()

		// renderer
		renderer := render.New(render.Options{
			Directory:                   "templates",
			Extensions:                  []string{".gohtml"},
			Layout:                      "layouts/layout",
			RenderPartialsWithoutPrefix: true,
			Funcs: []template.FuncMap{
				sprig.FuncMap(),
				mix.FuncMap(mix.Config{
					ManifestFile: "static/mix-manifest.json",
					PublicPath:   "/static/",
				}),
				urls.FuncMap(router),
				helpers.FuncMap(),
			},
		})

		// sessions & auth
		sessionName := viper.GetString("session-name")

		sessionStore := sessions.NewCookieStore([]byte(viper.GetString("session-secret")))
		sessionStore.MaxAge(viper.GetInt("session-max-age"))
		sessionStore.Options.Path = "/"
		sessionStore.Options.HttpOnly = true
		sessionStore.Options.Secure = strings.HasPrefix(baseURL, "https")

		oidcClient, err := oidc.New(oidc.Config{
			URL:          viper.GetString("oidc-url"),
			ClientID:     viper.GetString("oidc-client-id"),
			ClientSecret: viper.GetString("oidc-client-secret"),
			RedirectURL:  baseURL + "/auth/openid-connect/callback",
		})

		if err != nil {
			log.Fatal(err)
		}

		// add middleware
		router.Use(handlers.RecoveryHandler())

		// add routes
		routes.Register(
			baseURL,
			e,
			router,
			renderer,
			sessionName,
			sessionStore,
			oidcClient,
		)

		// logging
		handler := handlers.LoggingHandler(os.Stdout, router)

		// start server
		server.New(handler,
			server.WithHost(host),
			server.WithPort(port),
		).Start()
	},
}
