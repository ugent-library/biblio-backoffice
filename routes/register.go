package routes

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/jpillora/ipfilter"
	"github.com/leonelquinteros/gotext"
	"github.com/nics/ich"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/handlers/authenticating"
	"github.com/ugent-library/biblio-backoffice/handlers/candidaterecords"
	"github.com/ugent-library/biblio-backoffice/handlers/dashboard"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetcreating"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetediting"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetexporting"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetsearching"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetviewing"
	"github.com/ugent-library/biblio-backoffice/handlers/frontoffice"
	"github.com/ugent-library/biblio-backoffice/handlers/impersonating"
	"github.com/ugent-library/biblio-backoffice/handlers/mediatypes"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationbatch"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationcreating"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationediting"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationexporting"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationsearching"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationviewing"
	"github.com/ugent-library/httpx"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"
	"go.uber.org/zap"
)

type Version struct {
	Branch string
	Commit string
	Image  string
}

type Config struct {
	Version          Version
	Env              string
	Services         *backends.Services
	BaseURL          *url.URL
	Router           *ich.Mux
	Assets           mix.Manifest
	SessionStore     sessions.Store
	SessionName      string
	Timezone         *time.Location
	Loc              *gotext.Locale
	Logger           *zap.SugaredLogger
	OIDCAuth         *oidc.Auth
	FrontendURL      string
	FrontendUsername string
	FrontendPassword string
	IPRanges         string
	MaxFileSize      int
	CSRFName         string
	CSRFSecret       string
}

func Register(c Config) {
	c.Router.Use(middleware.RequestID)
	if c.Env != "local" {
		c.Router.Use(middleware.RealIP)
	}
	c.Router.Use(zaphttp.SetLogger(c.Logger.Desugar(), zapchi.RequestID))
	c.Router.Use(middleware.RequestLogger(zapchi.LogFormatter()))
	c.Router.Use(middleware.Recoverer)
	c.Router.Use(middleware.StripSlashes)

	// static files
	c.Router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// mount health and info
	c.Router.Get("/status", health.NewHandler(health.NewChecker())) // TODO add checkers
	c.Router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		httpx.RenderJSON(w, http.StatusOK, c.Version)
	})

	// handlers
	baseHandler := handlers.BaseHandler{
		Logger:          c.Logger,
		Router:          c.Router,
		SessionStore:    c.SessionStore,
		SessionName:     c.SessionName,
		Timezone:        c.Timezone,
		Loc:             c.Loc,
		UserService:     c.Services.UserService,
		BaseURL:         c.BaseURL,
		FrontendBaseUrl: c.FrontendURL,
	}
	authenticatingHandler := &authenticating.Handler{
		BaseHandler: baseHandler,
		OIDCAuth:    c.OIDCAuth,
	}
	impersonatingHandler := &impersonating.Handler{
		BaseHandler:       baseHandler,
		UserSearchService: c.Services.UserSearchService,
	}
	dashboardHandler := &dashboard.Handler{
		BaseHandler:            baseHandler,
		DatasetSearchIndex:     c.Services.DatasetSearchIndex,
		PublicationSearchIndex: c.Services.PublicationSearchIndex,
	}
	frontofficeHandler := &frontoffice.Handler{
		BaseHandler:      baseHandler,
		Repo:             c.Services.Repo,
		FileStore:        c.Services.FileStore,
		FrontendUsername: c.FrontendUsername,
		FrontendPassword: c.FrontendPassword,
		IPRanges:         c.IPRanges,
		IPFilter: ipfilter.New(ipfilter.Options{
			AllowedIPs:     strings.Split(c.IPRanges, ","),
			BlockByDefault: true,
		}),
	}
	datasetSearchingHandler := &datasetsearching.Handler{
		BaseHandler:        baseHandler,
		DatasetSearchIndex: c.Services.DatasetSearchIndex,
	}
	datasetExportingHandler := &datasetexporting.Handler{
		BaseHandler:          baseHandler,
		DatasetListExporters: c.Services.DatasetListExporters,
		DatasetSearchIndex:   c.Services.DatasetSearchIndex,
	}
	datasetViewingHandler := &datasetviewing.Handler{
		BaseHandler: baseHandler,
		Repo:        c.Services.Repo,
	}
	datasetCreatingHandler := &datasetcreating.Handler{
		BaseHandler:         baseHandler,
		Repo:                c.Services.Repo,
		DatasetSearchIndex:  c.Services.DatasetSearchIndex,
		DatasetSources:      c.Services.DatasetSources,
		OrganizationService: c.Services.OrganizationService,
	}
	datasetEditingHandler := &datasetediting.Handler{
		BaseHandler:               baseHandler,
		Repo:                      c.Services.Repo,
		ProjectService:            c.Services.ProjectService,
		ProjectSearchService:      c.Services.ProjectSearchService,
		OrganizationSearchService: c.Services.OrganizationSearchService,
		OrganizationService:       c.Services.OrganizationService,
		PersonSearchService:       c.Services.PersonSearchService,
		PersonService:             c.Services.PersonService,
		PublicationSearchIndex:    c.Services.PublicationSearchIndex,
	}
	publicationSearchingHandler := &publicationsearching.Handler{
		BaseHandler:            baseHandler,
		PublicationSearchIndex: c.Services.PublicationSearchIndex,
		FileStore:              c.Services.FileStore,
	}
	publicationExportingHandler := &publicationexporting.Handler{
		BaseHandler:              baseHandler,
		PublicationListExporters: c.Services.PublicationListExporters,
		PublicationSearchIndex:   c.Services.PublicationSearchIndex,
	}
	publicationViewingHandler := &publicationviewing.Handler{
		BaseHandler: baseHandler,
		Repo:        c.Services.Repo,
		FileStore:   c.Services.FileStore,
		MaxFileSize: c.MaxFileSize,
	}
	publicationCreatingHandler := &publicationcreating.Handler{
		BaseHandler:            baseHandler,
		Repo:                   c.Services.Repo,
		PublicationSearchIndex: c.Services.PublicationSearchIndex,
		PublicationSources:     c.Services.PublicationSources,
		PublicationDecoders:    c.Services.PublicationDecoders,
		OrganizationService:    c.Services.OrganizationService,
	}
	publicationEditingHandler := &publicationediting.Handler{
		BaseHandler:               baseHandler,
		Repo:                      c.Services.Repo,
		ProjectService:            c.Services.ProjectService,
		ProjectSearchService:      c.Services.ProjectSearchService,
		OrganizationSearchService: c.Services.OrganizationSearchService,
		OrganizationService:       c.Services.OrganizationService,
		PersonSearchService:       c.Services.PersonSearchService,
		PersonService:             c.Services.PersonService,
		DatasetSearchIndex:        c.Services.DatasetSearchIndex,
		FileStore:                 c.Services.FileStore,
		MaxFileSize:               c.MaxFileSize,
	}
	publicationBatchHandler := &publicationbatch.Handler{
		BaseHandler: baseHandler,
		Repo:        c.Services.Repo,
	}

	mediaTypesHandler := &mediatypes.Handler{
		BaseHandler:            baseHandler,
		MediaTypeSearchService: c.Services.MediaTypeSearchService,
	}

	// frontoffice data exchange api
	c.Router.Get("/frontoffice/publication/{id}", frontofficeHandler.BasicAuth(frontofficeHandler.GetPublication))
	c.Router.Get("/frontoffice/publication", frontofficeHandler.BasicAuth(frontofficeHandler.GetAllPublications))
	c.Router.Get("/frontoffice/dataset/{id}", frontofficeHandler.BasicAuth(frontofficeHandler.GetDataset))
	c.Router.Get("/frontoffice/dataset", frontofficeHandler.BasicAuth(frontofficeHandler.GetAllDatasets))
	// frontoffice file download
	c.Router.Get("/download/{id}/{file_id}", frontofficeHandler.DownloadFile)
	c.Router.Head("/download/{id}/{file_id}", frontofficeHandler.DownloadFile)

	c.Router.Group(func(r *ich.Mux) {
		r.Use(httpx.MethodOverride) // TODO eliminate need for method override with htmx
		r.Use(csrf.Protect(
			[]byte(c.CSRFSecret),
			csrf.CookieName(c.CSRFName),
			csrf.Path("/"),
			csrf.Secure(c.BaseURL.Scheme == "https"),
			csrf.SameSite(csrf.SameSiteStrictMode),
			csrf.FieldName("csrf-token"),
		))
		r.Use(secure.New(secure.Options{
			IsDevelopment: c.Env == "local",
			ContentSecurityPolicy: (&cspbuilder.Builder{
				Directives: map[string][]string{
					cspbuilder.DefaultSrc: {"'self'"},
					cspbuilder.ScriptSrc:  {"'self'", "$NONCE"},
					// TODO: htmx injects style
					cspbuilder.StyleSrc: {"'self'", "'unsafe-inline'"},
					cspbuilder.ImgSrc:   {"'self'", "data:"},
				},
			}).MustBuild(),
		}).Handler)

		// BEGIN NEW STYLE HANDLERS
		r.Group(func(r *ich.Mux) {
			r.Use(ctx.Set(ctx.Config{
				Services: c.Services,
				Router:   c.Router,
				Assets:   c.Assets,
				Timezone: c.Timezone,
				Loc:      c.Loc,
				Env:      c.Env,
				ErrorHandlers: map[int]http.HandlerFunc{
					http.StatusNotFound:            handlers.NotFound,
					http.StatusInternalServerError: handlers.InternalServerError,
				},
				SessionName:  c.SessionName,
				SessionStore: c.SessionStore,
				BaseURL:      c.BaseURL,
				FrontendURL:  c.FrontendURL,
				CSRFName:     "csrf-token",
			}))

			r.NotFound(handlers.NotFound)

			// home
			r.Get("/", handlers.Home).Name("home")

			r.Group(func(r *ich.Mux) {
				r.Use(ctx.RequireUser)

				r.With(ctx.SetNav("dashboard")).Get("/dashboard", handlers.DashBoard).Name("dashboard")
				r.Get("/dashboard-icon", handlers.DashBoardIcon).Name("dashboard_icon")
				// dashboard action required component
				r.Get("/action-required", handlers.ActionRequired).Name("action_required")
				// dashboard drafts to complete component
				r.Get("/drafts-to-complete", handlers.DraftsToComplete).Name("drafts_to_complete")
				// dashboard recent activity component
				r.Get("/recent-activity", handlers.RecentActivity).Name("recent_activity")
				// all candidate records
				r.With(ctx.SetNav("candidate_records")).Get("/candidate-records", candidaterecords.CandidateRecords).Name("candidate_records")
				r.Get("/candidate-records-icon", candidaterecords.CandidateRecordsIcon).Name("candidate_records_icon")
				r.Get("/candidate-records/{id}/confirm-reject", candidaterecords.ConfirmRejectCandidateRecord).Name("confirm_reject_candidate_record")
				r.Put("/candidate-records/{id}/reject", candidaterecords.RejectCandidateRecord).Name("reject_candidate_record")
				r.Put("/candidate-records/{id}/import", candidaterecords.ImportCandidateRecord).Name("import_candidate_record")
			})
		})
		// END NEW STYLE HANDLERS

		// authenticate user
		r.Get("/auth/openid-connect/callback",
			authenticatingHandler.Wrap(authenticatingHandler.Callback))
		r.Get("/login",
			authenticatingHandler.Wrap(authenticatingHandler.Login)).
			Name("login")
		r.Get("/logout",
			authenticatingHandler.Wrap(authenticatingHandler.Logout)).
			Name("logout")
		// change user role
		r.Put("/role/{role}",
			authenticatingHandler.Wrap(authenticatingHandler.UpdateRole)).
			Name("update_role")

		// impersonate user
		r.Get("/impersonation/add",
			impersonatingHandler.Wrap(impersonatingHandler.AddImpersonation)).
			Name("add_impersonation")
		r.Get("/impersonation/suggestions",
			impersonatingHandler.Wrap(impersonatingHandler.AddImpersonationSuggest)).
			Name("suggest_impersonations")
		r.Post("/impersonation",
			impersonatingHandler.Wrap(impersonatingHandler.CreateImpersonation)).
			Name("create_impersonation")
		// TODO why doesn't a DELETE with methodoverride work here?
		r.Post("/delete-impersonation",
			impersonatingHandler.Wrap(impersonatingHandler.DeleteImpersonation)).
			Name("delete_impersonation")

		// dashboard
		r.Get("/dashboard/publications/{type}", dashboardHandler.Wrap(dashboardHandler.Publications)).
			Name("dashboard_publications")
		r.Get("/dashboard/datasets/{type}", dashboardHandler.Wrap(dashboardHandler.Datasets)).
			Name("dashboard_datasets")
		r.Post("/dashboard/refresh-apublications/{type}", dashboardHandler.Wrap(dashboardHandler.RefreshAPublications)).
			Name("dashboard_refresh_apublications")
		r.Post("/dashboard/refresh-upublications/{type}", dashboardHandler.Wrap(dashboardHandler.RefreshUPublications)).
			Name("dashboard_refresh_upublications")

		// search datasets
		r.Get("/dataset",
			datasetSearchingHandler.Wrap(datasetSearchingHandler.Search)).
			Name("datasets")

		// add dataset
		r.Get("/dataset/add",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.Add)).
			Name("dataset_add")
		r.Post("/dataset/add",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.Add))
		r.Post("/dataset/import",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddImport)).
			Name("dataset_add_import")
		r.Post("/dataset/import/confirm",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.ConfirmImport)).
			Name("dataset_confirm_import")
		r.Get("/dataset/{id}/add/description",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddDescription)).
			Name("dataset_add_description")
		r.Get("/dataset/{id}/add/confirm",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddConfirm)).
			Name("dataset_add_confirm")
		r.Post("/dataset/{id}/save",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddSaveDraft)).
			Name("dataset_add_save_draft")
		r.Post("/dataset/{id}/add/publish",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddPublish)).
			Name("dataset_add_publish")
		r.Get("/dataset/{id}/add/finish",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddFinish)).
			Name("dataset_add_finish")

		// export datasets
		r.Get("/dataset.{format}",
			datasetExportingHandler.Wrap(datasetExportingHandler.ExportByCurationSearch)).
			Name("export_datasets")

		// view dataset
		r.Get("/dataset/{id}",
			datasetViewingHandler.Wrap(datasetViewingHandler.Show)).
			Name("dataset")
		r.Get("/dataset/{id}/description",
			datasetViewingHandler.Wrap(datasetViewingHandler.ShowDescription)).
			Name("dataset_description")
		r.Get("/dataset/{id}/contributors",
			datasetViewingHandler.Wrap(datasetViewingHandler.ShowContributors)).
			Name("dataset_contributors")
		r.Get("/dataset/{id}/publications",
			datasetViewingHandler.Wrap(datasetViewingHandler.ShowPublications)).
			Name("dataset_publications")
		r.Get("/dataset/{id}/activity",
			datasetViewingHandler.Wrap(datasetViewingHandler.ShowActivity)).
			Name("dataset_activity")

		// publish dataset
		r.Get("/dataset/{id}/publish/confirm",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmPublish)).
			Name("dataset_confirm_publish")
		r.Post("/dataset/{id}/publish",
			datasetEditingHandler.Wrap(datasetEditingHandler.Publish)).
			Name("dataset_publish")

		// withdraw dataset
		r.Get("/dataset/{id}/publish/withdraw",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmWithdraw)).
			Name("dataset_confirm_withdraw")
		r.Post("/dataset/{id}/withdraw",
			datasetEditingHandler.Wrap(datasetEditingHandler.Withdraw)).
			Name("dataset_withdraw")

		// re-publish dataset
		r.Get("/dataset/{id}/republish/confirm",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmRepublish)).
			Name("dataset_confirm_republish")
		r.Post("/dataset/{id}/republish",
			datasetEditingHandler.Wrap(datasetEditingHandler.Republish)).
			Name("dataset_republish")

		// lock dataset
		r.Post("/dataset/{id}/lock",
			datasetEditingHandler.Wrap(datasetEditingHandler.Lock)).
			Name("dataset_lock")
		r.Post("/dataset/{id}/unlock",
			datasetEditingHandler.Wrap(datasetEditingHandler.Unlock)).
			Name("dataset_unlock")

		// delete dataset
		r.Get("/dataset/{id}/confirm-delete",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDelete)).
			Name("dataset_confirm_delete")
		r.Delete("/dataset/{id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.Delete)).
			Name("dataset_delete")

		// edit dataset activity
		r.Get("/dataset/{id}/message/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditMessage)).
			Name("dataset_edit_message")
		r.Put("/dataset/{id}/message",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateMessage)).
			Name("dataset_update_message")
		r.Get("/dataset/{id}/reviewer-tags/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditReviewerTags)).
			Name("dataset_edit_reviewer_tags")
		r.Put("/dataset/{id}/reviewer-tags",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateReviewerTags)).
			Name("dataset_update_reviewer_tags")
		r.Get("/dataset/{id}/reviewer-note/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditReviewerNote)).
			Name("dataset_edit_reviewer_note")
		r.Put("/dataset/{id}/reviewer-note",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateReviewerNote)).
			Name("dataset_update_reviewer_note")

		// edit dataset details
		r.Get("/dataset/{id}/details/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditDetails)).
			Name("dataset_edit_details")
		r.Put("/dataset/{id}/details/edit/refresh-form",
			datasetEditingHandler.Wrap(datasetEditingHandler.RefreshEditFileForm)).
			Name("dataset_edit_file_refresh_form")
		r.Put("/dataset/{id}/details",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateDetails)).
			Name("dataset_update_details")

		// edit dataset projects
		r.Get("/dataset/{id}/projects/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddProject)).
			Name("dataset_add_project")
		r.Get("/dataset/{id}/projects/suggestions",
			datasetEditingHandler.Wrap(datasetEditingHandler.SuggestProjects)).
			Name("dataset_suggest_projects")
		r.Post("/dataset/{id}/projects",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateProject)).
			Name("dataset_create_project")
		r.Get("/dataset/{id}/{snapshot_id}/projects/confirm-delete/{project_id:.+}",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteProject)).
			Name("dataset_confirm_delete_project")
		r.Delete("/dataset/{id}/projects/{project_id:.+}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteProject)).
			Name("dataset_delete_project")

		// edit dataset links
		r.Get("/dataset/{id}/links/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddLink)).
			Name("dataset_add_link")
		r.Post("/dataset/{id}/links",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateLink)).
			Name("dataset_create_link")
		r.Get("/dataset/{id}/links/{link_id}/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditLink)).
			Name("dataset_edit_link")
		r.Put("/dataset/{id}/links/{link_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateLink)).
			Name("dataset_update_link")
		r.Get("/dataset/{id}/{snapshot_id}/links/{link_id}/confirm-delete",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteLink)).
			Name("dataset_confirm_delete_link")
		r.Delete("/dataset/{id}/links/{link_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteLink)).
			Name("dataset_delete_link")

		// edit dataset departments
		r.Get("/dataset/{id}/departments/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddDepartment)).
			Name("dataset_add_department")
		r.Get("/dataset/{id}/departments/suggestions",
			datasetEditingHandler.Wrap(datasetEditingHandler.SuggestDepartments)).
			Name("dataset_suggest_departments")
		r.Post("/dataset/{id}/departments",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateDepartment)).
			Name("dataset_create_department")
		r.Get("/dataset/{id}/{snapshot_id}/departments/{department_id}/confirm-delete",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteDepartment)).
			Name("dataset_confirm_delete_department")
		r.Delete("/dataset/{id}/departments/{department_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteDepartment)).
			Name("dataset_delete_department")

		// edit dataset abstracts
		r.Get("/dataset/{id}/abstracts/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddAbstract)).
			Name("dataset_add_abstract")
		r.Post("/dataset/{id}/abstracts",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateAbstract)).
			Name("dataset_create_abstract")
		r.Get("/dataset/{id}/abstracts/{abstract_id}/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditAbstract)).
			Name("dataset_edit_abstract")
		r.Put("/dataset/{id}/abstracts/{abstract_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateAbstract)).
			Name("dataset_update_abstract")
		r.Get("/dataset/{id}/{snapshot_id}/abstracts/{abstract_id}/confirm-delete",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteAbstract)).
			Name("dataset_confirm_delete_abstract")
		r.Delete("/dataset/{id}/abstracts/{abstract_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteAbstract)).
			Name("dataset_delete_abstract")

		// edit dataset publications
		r.Get("/dataset/{id}/publications/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddPublication)).
			Name("dataset_add_publication")
		r.Get("/dataset/{id}/publications/suggestions",
			datasetEditingHandler.Wrap(datasetEditingHandler.SuggestPublications)).
			Name("dataset_suggest_publications")
		r.Post("/dataset/{id}/publications",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreatePublication)).
			Name("dataset_create_publication")
		r.Get("/dataset/{id}/{snapshot_id}/publications/{publication_id}/confirm-delete",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeletePublication)).
			Name("dataset_confirm_delete_publication")
		r.Delete("/dataset/{id}/publications/{publication_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeletePublication)).
			Name("dataset_delete_publication")

		// edit dataset contributors
		r.Post("/dataset/{id}/contributors/{role}/order",
			datasetEditingHandler.Wrap(datasetEditingHandler.OrderContributors)).
			Name("dataset_order_contributors")
		r.Get("/dataset/{id}/contributors/{role}/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddContributor)).
			Name("dataset_add_contributor")
		r.Get("/dataset/{id}/contributors/{role}/suggestions",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddContributorSuggest)).
			Name("dataset_add_contributor_suggest")
		r.Get("/dataset/{id}/contributors/{role}/confirm-create",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmCreateContributor)).
			Name("dataset_confirm_create_contributor")
		r.Post("/dataset/{id}/contributors/{role}",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateContributor)).
			Name("dataset_create_contributor")
		r.Get("/dataset/{id}/contributors/{role}/{position}/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditContributor)).
			Name("dataset_edit_contributor")
		r.Get("/dataset/{id}/contributors/{role}/{position}/suggestions",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditContributorSuggest)).
			Name("dataset_edit_contributor_suggest")
		r.Get("/dataset/{id}/contributors/{role}/{position}/confirm-update",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmUpdateContributor)).
			Name("dataset_confirm_update_contributor")
		r.Put("/dataset/{id}/contributors/{role}/{position}",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateContributor)).
			Name("dataset_update_contributor")
		r.Get("/dataset/{id}/contributors/{role}/{position}/confirm-delete",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteContributor)).
			Name("dataset_confirm_delete_contributor")
		r.Delete("/dataset/{id}/contributors/{role}/{position}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteContributor)).
			Name("dataset_delete_contributor")

		// add publication
		r.Get("/publication/add",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.Add)).
			Name("publication_add")
		r.Post("/publication/add-single/import",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleImport)).
			Name("publication_add_single_import")
		r.Post("/publication/add-single/import/confirm",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleImportConfirm)).
			Name("publication_add_single_import_confirm")
		r.Get("/publication/{id}/add/description",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleDescription)).
			Name("publication_add_single_description")
		r.Get("/publication/{id}/add/confirm",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleConfirm)).
			Name("publication_add_single_confirm")
		r.Post("/publication/{id}/add/publish",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSinglePublish)).
			Name("publication_add_single_publish")
		r.Get("/publication/{id}/add/finish",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleFinish)).
			Name("publication_add_single_finish")
		r.Post("/publication/add-multiple/import",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleImport)).
			Name("publication_add_multiple_import")
		r.Get("/publication/add-multiple/{batch_id}/confirm",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleConfirm)).
			Name("publication_add_multiple_confirm")
		r.Get("/publication/add-multiple/{batch_id}/publication/{id}",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleShow)).
			Name("publication_add_multiple_show")
		r.Post("/publication/add-multiple/{batch_id}/save",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleSave)).
			Name("publication_add_multiple_save_draft")
		r.Post("/publication/add-multiple/{batch_id}/publish",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultiplePublish)).
			Name("publication_add_multiple_publish")
		r.Get("/publication/add-multiple/{batch_id}/finish",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleFinish)).
			Name("publication_add_multiple_finish")

		// search publications
		r.Get("/publication",
			publicationSearchingHandler.Wrap(publicationSearchingHandler.Search)).
			Name("publications")

		// export publications
		r.Get("/publication.{format}",
			publicationExportingHandler.Wrap(publicationExportingHandler.ExportByCurationSearch)).
			Name("export_publications")

		// publication batch operations
		r.Get("/publication/batch",
			publicationBatchHandler.Wrap(publicationBatchHandler.Show)).
			Name("publication_batch")
		r.Post("/publication/batch",
			publicationBatchHandler.Wrap(publicationBatchHandler.Process)).
			Name("publication_process_batch")

		// view publication
		r.Get("/publication/{id}",
			publicationViewingHandler.Wrap(publicationViewingHandler.Show)).
			Name("publication")
		r.Get("/publication/{id}/description",
			publicationViewingHandler.Wrap(publicationViewingHandler.ShowDescription)).
			Name("publication_description")
		r.Get("/publication/{id}/files",
			publicationViewingHandler.Wrap(publicationViewingHandler.ShowFiles)).
			Name("publication_files")
		r.Get("/publication/{id}/contributors",
			publicationViewingHandler.Wrap(publicationViewingHandler.ShowContributors)).
			Name("publication_contributors")
		r.Get("/publication/{id}/datasets",
			publicationViewingHandler.Wrap(publicationViewingHandler.ShowDatasets)).
			Name("publication_datasets")
		r.Get("/publication/{id}/activity",
			publicationViewingHandler.Wrap(publicationViewingHandler.ShowActivity)).
			Name("publication_activity")
		r.Get("/publication/{id}/files/{file_id}",
			publicationViewingHandler.Wrap(publicationViewingHandler.DownloadFile)).
			Name("publication_download_file")

		// publish publication
		r.Get("/publication/{id}/publish/confirm",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmPublish)).
			Name("publication_confirm_publish")
		r.Post("/publication/{id}/publish",
			publicationEditingHandler.Wrap(publicationEditingHandler.Publish)).
			Name("publication_publish")

		// withdraw publication
		r.Get("/publication/{id}/withdraw/confirm",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmWithdraw)).
			Name("publication_confirm_withdraw")
		r.Post("/publication/{id}/withdraw",
			publicationEditingHandler.Wrap(publicationEditingHandler.Withdraw)).
			Name("publication_withdraw")

		// re-publish publication
		r.Get("/publication/{id}/republish/confirm",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmRepublish)).
			Name("publication_confirm_republish")
		r.Post("/publication/{id}/republish",
			publicationEditingHandler.Wrap(publicationEditingHandler.Republish)).
			Name("publication_republish")

		// lock publication
		r.Post("/publication/{id}/lock",
			publicationEditingHandler.Wrap(publicationEditingHandler.Lock)).
			Name("publication_lock")
		r.Post("/publication/{id}/unlock",
			publicationEditingHandler.Wrap(publicationEditingHandler.Unlock)).
			Name("publication_unlock")

		// delete publication
		r.Get("/publication/{id}/confirm-delete",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDelete)).
			Name("publication_confirm_delete")
		r.Delete("/publication/{id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.Delete)).
			Name("publication_delete")

		// edit publication activity
		r.Get("/publication/{id}/message/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditMessage)).
			Name("publication_edit_message")
		r.Put("/publication/{id}/message",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateMessage)).
			Name("publication_update_message")
		r.Get("/publication/{id}/reviewer-tags/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditReviewerTags)).
			Name("publication_edit_reviewer_tags")
		r.Put("/publication/{id}/reviewer-tags",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateReviewerTags)).
			Name("publication_update_reviewer_tags")
		r.Get("/publication/{id}/reviewer-note/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditReviewerNote)).
			Name("publication_edit_reviewer_note")
		r.Put("/publication/{id}/reviewer-note",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateReviewerNote)).
			Name("publication_update_reviewer_note")

		// edit publication details
		r.Get("/publication/{id}/details/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditDetails)).
			Name("publication_edit_details")
		r.Put("/publication/{id}/details",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateDetails)).
			Name("publication_update_details")

		// edit publication type
		r.Get("/publication/{id}/type/confirm",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmUpdateType)).
			Name("publication_confirm_update_type")
		r.Put("/publication/{id}/type",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateType)).
			Name("publication_update_type")

		// edit publication conference
		r.Get("/publication/{id}/conference/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditConference)).
			Name("publication_edit_conference")
		r.Put("/publication/{id}/conference",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateConference)).
			Name("publication_update_conference")

		// edit publication additional info
		r.Get("/publication/{id}/additional-info/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditAdditionalInfo)).
			Name("publication_edit_additional_info")
		r.Put("/publication/{id}/additional-info",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateAdditionalInfo)).
			Name("publication_update_additional_info")

		// edit publication projects
		r.Get("/publication/{id}/projects/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddProject)).
			Name("publication_add_project")
		r.Get("/publication/{id}/projects/suggestions",
			publicationEditingHandler.Wrap(publicationEditingHandler.SuggestProjects)).
			Name("publication_suggest_projects")
		r.Post("/publication/{id}/projects",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateProject)).
			Name("publication_create_project")
		// project_id is last part of url because some id's contain slashes
		r.Get("/publication/{id}/{snapshot_id}/projects/confirm-delete/{project_id:.+}",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteProject)).
			Name("publication_confirm_delete_project")
		// project_id is last part of url because some id's contain slashes
		r.Delete("/publication/{id}/projects/{project_id:.+}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteProject)).
			Name("publication_delete_project")

		// edit publication links
		r.Get("/publicaton/{id}/links/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddLink)).
			Name("publication_add_link")
		r.Post("/publication/{id}/links",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateLink)).
			Name("publication_create_link")
		r.Get("/publication/{id}/links/{link_id}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditLink)).
			Name("publication_edit_link")
		r.Put("/publication/{id}/links/{link_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateLink)).
			Name("publication_update_link")
		r.Get("/publication/{id}/{snapshot_id}/links/{link_id}/confirm-delete",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteLink)).
			Name("publication_confirm_delete_link")
		r.Delete("/publication/{id}/links/{link_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteLink)).
			Name("publication_delete_link")

		// edit publication departments
		r.Get("/publication/{id}/departments/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddDepartment)).
			Name("publication_add_department")
		r.Get("/publication/{id}/departments/suggestions",
			publicationEditingHandler.Wrap(publicationEditingHandler.SuggestDepartments)).
			Name("publication_suggest_departments")
		r.Post("/publication/{id}/departments",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateDepartment)).
			Name("publication_create_department")
		r.Get("/publication/{id}/{snapshot_id}/departments/{department_id}/confirm-delete",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteDepartment)).
			Name("publication_confirm_delete_department")
		r.Delete("/publication/{id}/departments/{department_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteDepartment)).
			Name("publication_delete_department")

		// edit publication abstracts
		r.Get("/publication/{id}/abstracts/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddAbstract)).
			Name("publication_add_abstract")
		r.Post("/publication/{id}/abstracts",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateAbstract)).
			Name("publication_create_abstract")
		r.Get("/publication/{id}/abstracts/{abstract_id}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditAbstract)).
			Name("publication_edit_abstract")
		r.Put("/publication/{id}/abstracts/{abstract_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateAbstract)).
			Name("publication_update_abstract")
		r.Get("/publication/{id}/{snapshot_id}/abstracts/{abstract_id}/confirm-delete",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteAbstract)).
			Name("publication_confirm_delete_abstract")
		r.Delete("/publication/{id}/abstracts/{abstract_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteAbstract)).
			Name("publication_delete_abstract")

		// edit publication lay summaries
		r.Get("/publication/{id}/lay_summaries/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddLaySummary)).
			Name("publication_add_lay_summary")
		r.Post("/publication/{id}/lay_summaries",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateLaySummary)).
			Name("publication_create_lay_summary")
		r.Get("/publication/{id}/lay_summaries/{lay_summary_id}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditLaySummary)).
			Name("publication_edit_lay_summary")
		r.Put("/publication/{id}/lay_summaries/{lay_summary_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateLaySummary)).
			Name("publication_update_lay_summary")
		r.Get("/publication/{id}/{snapshot_id}/lay_summaries/{lay_summary_id}/confirm-delete",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteLaySummary)).
			Name("publication_confirm_delete_lay_summary")
		r.Delete("/publication/{id}/lay_summaries/{lay_summary_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteLaySummary)).
			Name("publication_delete_lay_summary")

		// edit publication datasets
		r.Get("/publication/{id}/datasets/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddDataset)).
			Name("publication_add_dataset")
		r.Get("/publication/{id}/datasets/suggestions",
			publicationEditingHandler.Wrap(publicationEditingHandler.SuggestDatasets)).
			Name("publication_suggest_datasets")
		r.Post("/publication/{id}/datasets",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateDataset)).
			Name("publication_create_dataset")
		r.Get("/publication/{id}/{snapshot_id}/datasets/{dataset_id}/confirm-delete",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteDataset)).
			Name("publication_confirm_delete_dataset")
		r.Delete("/publication/{id}/datasets/{dataset_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteDataset)).
			Name("publication_delete_dataset")

		// edit publication contributors
		r.Post("/publication/{id}/contributors/{role}/order",
			publicationEditingHandler.Wrap(publicationEditingHandler.OrderContributors)).
			Name("publication_order_contributors")
		r.Get("/publication/{id}/contributors/{role}/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddContributor)).
			Name("publication_add_contributor")
		r.Get("/publication/{id}/contributors/{role}/suggestions",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddContributorSuggest)).
			Name("publication_add_contributor_suggest")
		r.Get("/publication/{id}/contributors/{role}/confirm-create",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmCreateContributor)).
			Name("publication_confirm_create_contributor")
		r.Post("/publication/{id}/contributors/{role}",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateContributor)).
			Name("publication_create_contributor")
		r.Get("/publication/{id}/contributors/{role}/{position}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditContributor)).
			Name("publication_edit_contributor")
		r.Get("/publication/{id}/contributors/{role}/{position}/suggestions",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditContributorSuggest)).
			Name("publication_edit_contributor_suggest")
		r.Get("/publication/{id}/contributors/{role}/{position}/confirm-update",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmUpdateContributor)).
			Name("publication_confirm_update_contributor")
		r.Put("/publication/{id}/contributors/{role}/{position}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateContributor)).
			Name("publication_update_contributor")
		r.Get("/publication/{id}/contributors/{role}/{position}/confirm-delete",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteContributor)).
			Name("publication_confirm_delete_contributor")
		r.Delete("/publication/{id}/contributors/{role}/{position}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteContributor)).
			Name("publication_delete_contributor")

		// edit publication files
		r.Post("/publication/{id}/files",
			publicationEditingHandler.Wrap(publicationEditingHandler.UploadFile)).
			Name("publication_upload_file")
		r.Get("/publication/{id}/files/{file_id}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditFile)).
			Name("publication_edit_file")
		r.Get("/publication/{id}/refresh-files",
			publicationEditingHandler.Wrap(publicationEditingHandler.RefreshFiles)).
			Name("publication_refresh_files")
		r.Get("/publication/{id}/files/{file_id}/refresh-form",
			publicationEditingHandler.Wrap(publicationEditingHandler.RefreshEditFileForm)).
			Name("publication_edit_file_refresh_form")
		r.Put("/publication/{id}/files/{file_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateFile)).
			Name("publication_update_file")
		r.Get("/publication/{id}/{snapshot_id}/files/{file_id}/confirm-delete",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteFile)).
			Name("publication_confirm_delete_file")
		r.Delete("/publication/{id}/files/{file_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteFile)).
			Name("publication_delete_file")

		// media types
		r.Get("/media_type/suggestions",
			mediaTypesHandler.Wrap(mediaTypesHandler.Suggest)).
			Name("suggest_media_types")
	})
}
