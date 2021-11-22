package commands

import (
	"crypto/tls"
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/controllers"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/helpers"
	"github.com/ugent-library/biblio-backend/internal/routes"
	"github.com/ugent-library/go-graceful/server"
	"github.com/ugent-library/go-locale/locale"
	"gopkg.in/cas.v2"

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
	e, err := engine.New(engine.Config{
		LibreCatURL:       viper.GetString("librecat-url"),
		LibreCatUsername:  viper.GetString("librecat-username"),
		LibreCatPassword:  viper.GetString("librecat-password"),
		ORCIDClientID:     viper.GetString("orcid-client-id"),
		ORCIDClientSecret: viper.GetString("orcid-client-secret"),
		ORCIDSandbox:      viper.GetBool("orcid-sandbox"),
	})
	if err != nil {
		log.Fatal(err)
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
	sessionStore.Options.Path = baseURL.Path
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = baseURL.Scheme == "https"

	// oidcClient, err := oidc.New(oidc.Config{
	// 	URL:          viper.GetString("oidc-url"),
	// 	ClientID:     viper.GetString("oidc-client-id"),
	// 	ClientSecret: viper.GetString("oidc-client-secret"),
	// 	RedirectURL:  baseURL.String() + "/auth/openid-connect/callback",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// controller config
	config := controllers.Context{
		Engine:       e,
		BaseURL:      baseURL,
		Router:       router,
		Render:       r,
		Localizer:    localizer,
		SessionName:  sessionName,
		SessionStore: sessionStore,
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

		// cas auth
		casURL, _ := url.Parse(viper.GetString("cas-url"))
		casOpts := &cas.Options{
			URL: casURL,
		}
		if viper.GetBool("cas-skip-verify-tls") {
			casOpts.Client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		}

		handler = cas.NewClient(casOpts).Handle(handler)

		// start server
		server.New(handler,
			server.WithHost(viper.GetString("host")),
			server.WithPort(viper.GetInt("port")),
			server.WithWriteTimeOut(3*time.Minute),
			server.WithReadTimeOut(3*time.Minute),
		).Start()
	},
}
