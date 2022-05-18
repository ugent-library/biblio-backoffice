package webapp

import (
	"context"
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
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/controllers"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/helpers"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/mix"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/routes"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/urls"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-oidc/oidc"
	"github.com/unrolled/render"

	_ "github.com/ugent-library/biblio-backend/services/webapp/internal/translations"
)

type service struct {
	server *http.Server
}

func AddCommands(cmd *cobra.Command, services *backends.Services) {
	cmd.AddCommand(&cobra.Command{
		Use:   "routes",
		Short: "print routes",
		Run: func(cmd *cobra.Command, args []string) {
			router := buildRouter(services)
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
	})
}

func New(e *backends.Services) (*service, error) {
	router := buildRouter(e)

	// logging
	handler := handlers.LoggingHandler(os.Stdout, router)

	// start server
	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  3 * time.Minute,
		WriteTimeout: 3 * time.Minute,
	}

	return &service{server}, nil
}

func (s *service) Name() string {
	return "biblio-backend-webapp"
}

func (s *service) Start() error {
	return s.server.ListenAndServe()
}

func (s *service) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func buildRouter(e *backends.Services) *mux.Router {
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
	r := render.New(render.Options{
		Directory:                   "services/webapp/templates",
		Extensions:                  []string{".gohtml"},
		Layout:                      "layouts/layout",
		RenderPartialsWithoutPrefix: true,
		Funcs: []template.FuncMap{
			sprig.FuncMap(),
			mix.FuncMap(mix.Config{
				ManifestFile: "services/webapp/static/mix-manifest.json",
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

	// controller base
	base := controllers.Base{
		Services:       e,
		Mode:         viper.GetString("mode"),
		BaseURL:      baseURL,
		Router:       router,
		Render:       r,
		Localizer:    localizer,
		SessionName:  sessionName,
		SessionStore: sessionStore,
		OIDC:         oidcClient,
	}

	// add routes
	routes.Register(base)

	return router
}
