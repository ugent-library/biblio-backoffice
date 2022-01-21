package webapp

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/backends/ianamedia"
	"github.com/ugent-library/biblio-backend/internal/backends/librecat"
	"github.com/ugent-library/biblio-backend/internal/backends/spdxlicenses"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/controllers"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/helpers"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/routes"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-oidc/oidc"
	"github.com/ugent-library/go-orcid/orcid"
	"github.com/ugent-library/go-web/mix"
	"github.com/ugent-library/go-web/urls"
	"github.com/unrolled/render"

	_ "github.com/ugent-library/biblio-backend/services/webapp/internal/translations"
)

type service struct {
	server *http.Server
}

func New() (*service, error) {
	router := buildRouter()

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

func (s *service) Serve() error {
	return s.server.ListenAndServe()
}

func (s *service) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
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
		LicenseSearchService:      spdxlicenses.New(),
		MediaTypeSearchService:    ianamedia.New(),
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

	// controller config
	config := controllers.Context{
		Engine:       e,
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
	routes.Register(config)

	return router
}
