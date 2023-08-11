package commands

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/oklog/ulid/v2"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backoffice/internal/app/helpers"
	"github.com/ugent-library/biblio-backoffice/internal/app/routes"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/locale"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/internal/urls"
	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
	"github.com/ugent-library/middleware"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"go.uber.org/zap"

	_ "github.com/ugent-library/biblio-backoffice/internal/translations"
)

func init() {
	serverCmd.PersistentFlags().String("base-url", "", "base url")

	serverStartCmd.Flags().String("mode", defaultMode, "server mode (development, staging or production)")
	serverStartCmd.Flags().String("host", "", "server host")
	serverStartCmd.Flags().Int("port", defaultPort, "server port")
	serverStartCmd.Flags().String("session-name", defaultSessionName, "session name")
	serverStartCmd.Flags().String("session-secret", "", "session secret")
	serverStartCmd.Flags().Int("session-max-age", defaultSessionMaxAge, "session lifetime")
	serverStartCmd.Flags().String("csrf-name", "", "csrf cookie name")
	serverStartCmd.Flags().String("csrf-secret", "", "csrf cookie secret")
	serverStartCmd.Flags().String("location", defaultTimezone, "location used for date and time display")

	serverCmd.AddCommand(serverRoutesCmd)
	serverCmd.AddCommand(serverStartCmd)
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server [command]",
	Short: "The biblio-backoffice HTTP server",
}

var serverRoutesCmd = &cobra.Command{
	Use:   "routes",
	Short: "print routes",
	Run: func(cmd *cobra.Command, args []string) {
		router := buildRouter(Services(), newLogger())
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
		production := viper.GetString("mode") == "production"

		// setup logger
		logger := newLogger()

		// setup services
		e := Services()

		e.MediaTypeSearchService.IndexAll()
		// e.LicenseSearchService.IndexAll()

		// setup router
		router := buildRouter(e, logger)

		// apply these before request reaches the router
		handler := middleware.Apply(router,
			middleware.Recover(func(err any) {
				if production {
					logger.With(zap.Stack("stack")).Error(err)
				} else {
					logger.Error(err)
				}
			}),
			middleware.SetRequestID(func() string {
				return ulid.Make().String()
			}),
			zaphttp.LogRequests(logger.Desugar()),
		)

		// setup server
		addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
		server := graceful.WithDefaults(&http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  5 * time.Minute,
			WriteTimeout: 5 * time.Minute,
		})
		logger.Infof("starting server at %s", addr)
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			logger.Fatal(err)
		}
		logger.Info("gracefully stopped server")
	},
}

func buildRouter(services *backends.Services, logger *zap.SugaredLogger) *mux.Router {
	mode := viper.GetString("mode")

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
		logger.Fatal(err)
	}

	// router
	router := mux.NewRouter()

	// asets
	assets, err := mix.New(mix.Config{
		ManifestFile: "static/mix-manifest.json",
		PublicPath:   baseURL.Path + "/static/",
	})
	if err != nil {
		logger.Fatal(err)
	}

	// renderer
	funcMaps := []template.FuncMap{
		sprig.FuncMap(),
		urls.FuncMap(router),
		helpers.FuncMap(),
		{
			"assetPath": assets.AssetPath,
			"appMode": func() string { // TODO eliminate need for this
				return mode
			},
			"vocabulary": func(k string) []string { // TODO eliminate need for this?
				return vocabularies.Map[k]
			},
		},
	}

	// init render
	render.AuthURL = baseURL.Path + "/login"

	for _, funcs := range funcMaps {
		render.Funcs(funcs)
	}
	render.MustParse()

	// init bind
	bind.PathValuesFunc = func(r *http.Request) url.Values {
		p := url.Values{}
		for k, v := range mux.Vars(r) {
			p.Set(k, v)
		}
		return p
	}

	// localizer
	localizer := locale.NewLocalizer("en")

	//
	timezone, err := time.LoadLocation(viper.GetString("timezone"))
	if err != nil {
		logger.Fatal(err)
	}

	// sessions & auth
	sessionSecret := []byte(viper.GetString("session-secret"))
	sessionName := viper.GetString("session-name")
	sessionStore := sessions.NewCookieStore(sessionSecret)
	sessionStore.MaxAge(viper.GetInt("session-max-age"))
	sessionStore.Options.Path = "/"
	if baseURL.Path != "" {
		sessionStore.Options.Path = baseURL.Path
	}
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = baseURL.Scheme == "https"

	oidcAuth, err := oidc.NewAuth(context.TODO(), oidc.Config{
		URL:          viper.GetString("oidc-url"),
		ClientID:     viper.GetString("oidc-client-id"),
		ClientSecret: viper.GetString("oidc-client-secret"),
		RedirectURL:  baseURL.String() + "/auth/openid-connect/callback",
		CookieName:   viper.GetString("session-name") + ".state",
		CookieSecret: []byte(viper.GetString("session-secret")),
		Insecure:     mode != "production",
	})
	if err != nil {
		logger.Fatal(err)
	}

	// add routes
	routes.Register(services, baseURL, router, sessionStore, sessionName, timezone, localizer, logger, oidcAuth)

	return router
}
