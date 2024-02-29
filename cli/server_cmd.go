package cli

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonelquinteros/gotext"
	"github.com/nics/ich"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
	"github.com/thomaspoignant/go-feature-flag/retriever/githubretriever"
	"github.com/ugent-library/biblio-backoffice/api/v2"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/helpers"
	"github.com/ugent-library/biblio-backoffice/jobs"
	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/projects"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/routes"
	"github.com/ugent-library/biblio-backoffice/urls"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverStartCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Biblio backoffice HTTP server",
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := newServices()

		services.MediaTypeSearchService.IndexAll()
		// e.LicenseSearchService.IndexAll()

		// feature flags
		if config.FF.FilePath != "" {
			err := ffclient.Init(ffclient.Config{
				PollingInterval: 5 * time.Second,
				Context:         context.TODO(),
				Retriever: &fileretriever.Retriever{
					Path: config.FF.FilePath,
				},
			})
			if err != nil {
				return err
			}
			defer ffclient.Close()
		} else if config.FF.GitHubRepo != "" {
			err := ffclient.Init(ffclient.Config{
				PollingInterval: 5 * time.Second,
				Context:         context.TODO(),
				Retriever: &githubretriever.Retriever{
					GithubToken:    config.FF.GitHubToken,
					RepositorySlug: config.FF.GitHubRepo,
					Branch:         config.FF.GitHubBranch,
					FilePath:       config.FF.GitHubPath,
				},
			})
			if err != nil {
				return err
			}
			defer ffclient.Close()
		}

		// start jobs
		err := jobs.Start(context.TODO(), jobs.JobsConfig{
			PgxPool:       services.PgxPool,
			PeopleRepo:    services.PeopleRepo,
			PeopleIndex:   services.PeopleIndex,
			ProjectsRepo:  services.ProjectsRepo,
			ProjectsIndex: services.ProjectsIndex,
			Logger:        logger,
		})
		if err != nil {
			return err
		}

		// setup router
		router, err := buildRouter(services)
		if err != nil {
			return err
		}

		// setup server
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		server := graceful.WithDefaults(&http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  5 * time.Minute,
			WriteTimeout: 5 * time.Minute,
		})
		zapLogger.Infof("starting server at %s", addr)
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			return err
		}
		zapLogger.Info("gracefully stopped server")

		return nil
	},
}

func buildRouter(services *backends.Services) (*ich.Mux, error) {
	b := config.BaseURL
	if b == "" {
		if config.Host == "" {
			b = "http://localhost"
		} else {
			b = "http://" + config.Host
		}
		if config.Port != 80 {
			b = fmt.Sprintf("%s:%d", b, config.Port)
		}
	}
	baseURL, err := url.Parse(b)
	if err != nil {
		return nil, err
	}

	// router
	router := ich.New()

	// assets
	assets, err := mix.New(mix.Config{
		ManifestFile: "static/mix-manifest.json",
		PublicPath:   baseURL.Path + "/static/",
	})
	if err != nil {
		return nil, err
	}

	// renderer
	funcMaps := []template.FuncMap{
		sprig.FuncMap(),
		urls.FuncMap(router, baseURL.Scheme, baseURL.Host),
		helpers.FuncMap(),
		{
			"assetPath": assets.AssetPath,
			"appMode": func() string { // TODO eliminate need for this
				return config.Env
			},
			"vocabulary": func(k string) []string { // TODO eliminate need for this?
				return vocabularies.Map[k]
			},
		},
	}

	// init render
	render.AuthURL = "/login"

	for _, funcs := range funcMaps {
		render.Funcs(funcs)
	}
	render.MustParse()

	// init bind
	bind.PathValueFunc = chi.URLParam

	// locale
	loc := gotext.NewLocale("locales", "en")
	loc.AddDomain("default")

	// timezone
	timezone, err := time.LoadLocation(config.Timezone)
	if err != nil {
		return nil, err
	}

	// sessions & auth
	sessionSecret := []byte(config.Session.Secret)
	sessionName := config.Session.Name
	sessionStore := sessions.NewCookieStore(sessionSecret)
	sessionStore.MaxAge(config.Session.MaxAge)
	sessionStore.Options.Path = "/"
	if baseURL.Path != "" {
		sessionStore.Options.Path = baseURL.Path
	}
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = baseURL.Scheme == "https"

	oidcAuth, err := oidc.NewAuth(context.TODO(), oidc.Config{
		URL:          config.OIDC.URL,
		ClientID:     config.OIDC.ClientID,
		ClientSecret: config.OIDC.ClientSecret,
		RedirectURL:  baseURL.String() + "/auth/openid-connect/callback",
		CookieName:   config.Session.Name + ".state",
		CookieSecret: []byte(config.Session.Secret),
		Insecure:     config.Env != "local",
	})
	if err != nil {
		return nil, err
	}

	// api server
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.PgConn)
	if err != nil {
		return nil, err
	}

	peopleRepo, err := people.NewRepo(people.RepoConfig{
		Conn: pool,
	})
	if err != nil {
		return nil, err
	}

	projectsRepo, err := projects.NewRepo(projects.RepoConfig{
		Conn: pool,
	})
	if err != nil {
		return nil, err
	}

	apiService := api.NewService(peopleRepo, projectsRepo)
	apiServer, err := api.NewServer(apiService, &api.ApiSecurityHandler{APIKey: config.APIKey})
	if err != nil {
		return nil, err
	}

	// add routes
	routes.Register(routes.Config{
		Version: routes.Version{
			Branch: version.Branch,
			Commit: version.Commit,
			Image:  version.Image,
		},
		Env:              config.Env,
		Services:         services,
		BaseURL:          baseURL,
		Router:           router,
		Assets:           assets,
		SessionStore:     sessionStore,
		SessionName:      sessionName,
		Timezone:         timezone,
		Loc:              loc,
		Logger:           zapLogger,
		OIDCAuth:         oidcAuth,
		FrontendURL:      config.Frontend.URL,
		FrontendUsername: config.Frontend.Username,
		FrontendPassword: config.Frontend.Password,
		IPRanges:         config.IPRanges,
		MaxFileSize:      config.MaxFileSize,
		CSRFName:         config.CSRF.Name,
		CSRFSecret:       config.CSRF.Secret,
		ApiServer:        apiServer,
	})

	return router, nil
}
