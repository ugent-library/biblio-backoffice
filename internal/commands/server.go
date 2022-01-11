package commands

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/backends/librecat"
	"github.com/ugent-library/biblio-backend/internal/backends/spdx"
	"github.com/ugent-library/biblio-backend/internal/controllers"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/helpers"
	"github.com/ugent-library/biblio-backend/internal/routes"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-oidc/oidc"
	"github.com/ugent-library/go-orcid/orcid"

	"github.com/ugent-library/go-web/mix"
	"github.com/ugent-library/go-web/urls"
	"github.com/unrolled/render"
)

func init() {
	serverCmd.PersistentFlags().String("base-url", "", "base url")

	serverStartCmd.Flags().String("host", defaultHost, "server host")
	serverStartCmd.Flags().Int("port", defaultPort, "server port")
	serverStartCmd.Flags().String("session-name", defaultSessionName, "session name")
	serverStartCmd.Flags().String("session-secret", "", "session secret")
	serverStartCmd.Flags().Int("session-max-age", defaultSessionMaxAge, "session lifetime")
	serverStartCmd.Flags().String("csrf-name", "", "csrf cookie name")
	serverStartCmd.Flags().String("csrf-secret", "", "csrf cookie secret")

	serverCmd.AddCommand(serverRoutesCmd)
	serverCmd.AddCommand(serverStartCmd)
	rootCmd.AddCommand(serverCmd)
}

func buildRouter() *mux.Router {
	host := viper.GetString("host")
	port := viper.GetInt("port")

	b := viper.GetString("base-url")
	if b == "" {
		if host == "" {
			b = "http://localhost"
		} else {
			b = "http://" + host
		}
		if port != 80 {
			b = fmt.Sprintf("%s:%d", b, port)
		}
	}
	baseURL, err := url.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	// engine
	librecatClient := librecat.New(librecat.Config{
		URL:      viper.GetString("librecat-url"),
		Username: viper.GetString("librecat-username"),
		Password: viper.GetString("librecat-password"),
	})
	orcidConfig := orcid.Config{
		ClientID:     viper.GetString("orcid-client-id"),
		ClientSecret: viper.GetString("orcid-client-secret"),
		Sandbox:      viper.GetBool("orcid-sandbox"),
	}
	orcidClient := orcid.NewMemberClient(orcidConfig)
	e := &engine.Engine{
		ORCIDSandbox:              orcidConfig.Sandbox,
		ORCIDClient:               orcidClient,
		DatasetService:            librecatClient,
		DatasetSearchService:      librecatClient,
		PublicationService:        librecatClient,
		PublicationSearchService:  librecatClient,
		PersonService:             librecatClient,
		ProjectService:            librecatClient,
		UserService:               librecatClient,
		OrganizationSearchService: librecatClient,
		PersonSearchService:       librecatClient,
		ProjectSearchService:      librecatClient,
		LicenseSearchService:      spdx.New(),
	}

	// router
	router := mux.NewRouter()

	// renderer
	r := render.New(render.Options{
		Directory:                   "templates",
		Extensions:                  []string{".gohtml"},
		Layout:                      "layouts/layout",
		RenderPartialsWithoutPrefix: true,
		Funcs: []template.FuncMap{
			sprig.FuncMap(),
			mix.FuncMap(mix.Config{
				ManifestFile: "static/mix-manifest.json",
				PublicPath:   baseURL.Path + "/static/",
			}),
			urls.FuncMap(router),
			helpers.FuncMap(),
		},
	})

	// localizer
	localizer := locale.NewLocalizer("en")

	// sessions & auth
	sessionSecret := []byte(viper.GetString("session-secret"))
	sessionName := viper.GetString("session-name")
	sessionStore := sessions.NewCookieStore(sessionSecret)
	sessionStore.MaxAge(viper.GetInt("session-max-age"))
	if baseURL.Path != "" {
		sessionStore.Options.Path = baseURL.Path
	}
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = baseURL.Scheme == "https"

	oidcClient, err := oidc.New(oidc.Config{
		URL:          viper.GetString("oidc-url"),
		ClientID:     viper.GetString("oidc-client-id"),
		ClientSecret: viper.GetString("oidc-client-secret"),
		RedirectURL:  baseURL.String() + "/auth/openid-connect/callback",
	})
	if err != nil {
		log.Fatal(err)
	}

	// controller config
	config := controllers.Context{
		Engine:       e,
		BaseURL:      baseURL,
		Router:       router,
		Render:       r,
		Localizer:    localizer,
		SessionName:  sessionName,
		SessionStore: sessionStore,
		OIDC:         oidcClient,
	}

	// add routes
	routes.Register(config)

	return router
}

var serverCmd = &cobra.Command{
	Use:   "server [command]",
	Short: "The biblio-backend HTTP server",
}

var serverRoutesCmd = &cobra.Command{
	Use:   "routes",
	Short: "print routes",
	Run: func(cmd *cobra.Command, args []string) {
		router := buildRouter()
		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			hostTemplate, err := route.GetHostTemplate()
			if err == nil {
				fmt.Println("HOST:", hostTemplate)
			}
			pathTemplate, err := route.GetPathTemplate()
			if err == nil {
				fmt.Println("ROUTE:", pathTemplate)
			}
			pathRegexp, err := route.GetPathRegexp()
			if err == nil {
				fmt.Println("Path regexp:", pathRegexp)
			}
			queriesTemplates, err := route.GetQueriesTemplates()
			if err == nil {
				fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
			}
			queriesRegexps, err := route.GetQueriesRegexp()
			if err == nil {
				fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
			}
			methods, err := route.GetMethods()
			if err == nil {
				fmt.Println("Methods:", strings.Join(methods, ","))
			}
			fmt.Println()
			return nil
		})
	},
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the http server",
	Run: func(cmd *cobra.Command, args []string) {
		router := buildRouter()

		// logging
		handler := handlers.LoggingHandler(os.Stdout, router)

		// start server
		addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))

		server := graceful.WithDefaults(&http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  3 * time.Minute,
			WriteTimeout: 3 * time.Minute,
		})

		log.Println("Starting the server at " + addr)
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			log.Fatal(err)
		}
		log.Println("Server was shutdown gracefully")
	},
}
