package commands

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/app/helpers"
	"github.com/ugent-library/biblio-backend/internal/app/routes"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/mix"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/controllers"
	"github.com/ugent-library/biblio-backend/internal/urls"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
	"github.com/ugent-library/go-oidc/oidc"
	unrolledrender "github.com/unrolled/render"

	_ "github.com/ugent-library/biblio-backend/internal/translations"
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

	serverCmd.AddCommand(serverRoutesCmd)
	serverCmd.AddCommand(serverStartCmd)
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server [command]",
	Short: "The biblio-backend HTTP server",
}

var serverRoutesCmd = &cobra.Command{
	Use:   "routes",
	Short: "print routes",
	Run: func(cmd *cobra.Command, args []string) {
		router := buildRouter(Services())
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
		// setup services
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

		// setup router
		router := buildRouter(e)

		// setup logging
		handler := handlers.LoggingHandler(os.Stdout, router)

		// setup server
		addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))

		server := &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  3 * time.Minute,
			WriteTimeout: 3 * time.Minute,
		}

		// start server
		ctx, stop := signal.NotifyContext(context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		errC := make(chan error)

		// listen for shutdown signal
		go func() {
			<-ctx.Done()

			log.Println("Stopping gracefully...")

			timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			defer func() {
				stop()
				cancel()
				close(errC)
			}()

			// disable keep-alive on shutdown
			server.SetKeepAlivesEnabled(false)

			if err := server.Shutdown(timeoutCtx); err != nil {
				errC <- err
			}

			log.Println("Stopped")
		}()

		go func() {
			log.Printf("Listening at %s", addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errC <- err
			}
		}()

		if err := <-errC; err != nil {
			log.Fatalf("Error while running: %s", err)
		}
	},
}

func buildRouter(services *backends.Services) *mux.Router {
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

	// router
	router := mux.NewRouter()

	// renderer
	funcMaps := []template.FuncMap{
		sprig.FuncMap(),
		mix.FuncMap(mix.Config{
			ManifestFile: "static/mix-manifest.json",
			PublicPath:   baseURL.Path + "/static/",
		}),
		urls.FuncMap(router),
		helpers.FuncMap(),
		{
			"appMode": func() string { // TODO eliminate need for this
				return viper.GetString("mode")
			},
			"vocabulary": func(k string) []string { // TODO eliminate need for this?
				return vocabularies.Map[k]
			},
		},
	}

	r := unrolledrender.New(unrolledrender.Options{
		Directory:                   "internal/services/webapp/templates",
		Extensions:                  []string{".gohtml"},
		Layout:                      "layouts/layout",
		RenderPartialsWithoutPrefix: true,
		Funcs:                       funcMaps,
	})

	// init render
	render.FuncMaps = funcMaps
	render.AuthURL = baseURL.Path + "/login"

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

	// controller base
	base := controllers.Base{
		Mode:         viper.GetString("mode"),
		BaseURL:      baseURL,
		Router:       router,
		Render:       r,
		Localizer:    localizer,
		SessionName:  sessionName,
		SessionStore: sessionStore,
	}

	// add routes
	routes.Register(services, base, oidcClient)

	return router
}
