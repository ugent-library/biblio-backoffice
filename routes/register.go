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
	"github.com/nics/ich"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/handlers/authenticating"
	"github.com/ugent-library/biblio-backoffice/handlers/dashboard"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetcreating"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetediting"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetexporting"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetsearching"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetviewing"
	"github.com/ugent-library/biblio-backoffice/handlers/frontoffice"
	"github.com/ugent-library/biblio-backoffice/handlers/home"
	"github.com/ugent-library/biblio-backoffice/handlers/impersonating"
	"github.com/ugent-library/biblio-backoffice/handlers/mediatypes"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationbatch"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationcreating"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationediting"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationexporting"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationsearching"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationviewing"
	"github.com/ugent-library/biblio-backoffice/locale"
	mw "github.com/ugent-library/middleware"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
	"go.uber.org/zap"
)

type Config struct {
	Env              string
	Services         *backends.Services
	BaseURL          *url.URL
	Router           *ich.Mux
	SessionStore     sessions.Store
	SessionName      string
	Timezone         *time.Location
	Localizer        *locale.Localizer
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
	c.Router.Use(mw.MethodOverride( // TODO eliminate need for method override
		mw.MethodFromHeader(mw.MethodHeader),
		mw.MethodFromForm(mw.MethodParam),
	))
	c.Router.Use(zaphttp.SetLogger(c.Logger.Desugar(), zapchi.RequestID))
	c.Router.Use(middleware.RequestLogger(zapchi.LogFormatter()))
	c.Router.Use(middleware.Recoverer)
	c.Router.Use(middleware.StripSlashes)

	// static files
	c.Router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// status endpoint
	// TODO add checkers
	c.Router.Get("/status", health.NewHandler(health.NewChecker()))
	// TODO add /info endpoint

	// handlers
	baseHandler := handlers.BaseHandler{
		Logger:          c.Logger,
		Router:          c.Router,
		SessionStore:    c.SessionStore,
		SessionName:     c.SessionName,
		Timezone:        c.Timezone,
		Localizer:       c.Localizer,
		UserService:     c.Services.UserService,
		BaseURL:         c.BaseURL,
		FrontendBaseUrl: c.FrontendURL,
	}
	homeHandler := &home.Handler{
		BaseHandler: baseHandler,
	}
	authenticatingHandler := &authenticating.Handler{
		BaseHandler: baseHandler,
		OIDCAuth:    c.OIDCAuth,
	}
	impersonatingHandler := &impersonating.Handler{
		BaseHandler:       baseHandler,
		UserSearchService: c.Services.UserSearchService,
	}
	// tasksHandler := &tasks.Handler{
	// 	BaseHandler: baseHandler,
	// 	Tasks:       services.Tasks,
	// }
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

	frontofficeRouter := c.Router.PathPrefix("").Subrouter()
	// frontoffice data exchange api
	frontofficeRouter.HandleFunc("/frontoffice/publication/{id}", frontofficeHandler.BasicAuth(frontofficeHandler.GetPublication)).
		Methods("GET")
	frontofficeRouter.HandleFunc("/frontoffice/publication", frontofficeHandler.BasicAuth(frontofficeHandler.GetAllPublications)).
		Methods("GET")
	frontofficeRouter.HandleFunc("/frontoffice/dataset/{id}", frontofficeHandler.BasicAuth(frontofficeHandler.GetDataset)).
		Methods("GET")
	frontofficeRouter.HandleFunc("/frontoffice/dataset", frontofficeHandler.BasicAuth(frontofficeHandler.GetAllDatasets)).
		Methods("GET")
	// frontoffice file download
	frontofficeRouter.HandleFunc("/download/{id}/{file_id}", frontofficeHandler.DownloadFile).
		Methods("GET", "HEAD")

	r := c.Router.PathPrefix("").Subrouter()
	r.Use(csrf.Protect(
		[]byte(c.CSRFSecret),
		csrf.CookieName(c.CSRFName),
		csrf.Path("/"),
		csrf.Secure(c.BaseURL.Scheme == "https"),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.FieldName("csrf-token"),
	))

	// home
	r.HandleFunc("/",
		homeHandler.Wrap(homeHandler.Home)).
		Methods("GET").
		Name("home")

	// authenticate user
	r.HandleFunc("/auth/openid-connect/callback",
		authenticatingHandler.Wrap(authenticatingHandler.Callback)).
		Methods("GET")
	r.HandleFunc("/login",
		authenticatingHandler.Wrap(authenticatingHandler.Login)).
		Methods("GET").
		Name("login")
	r.HandleFunc("/logout",
		authenticatingHandler.Wrap(authenticatingHandler.Logout)).
		Methods("GET").
		Name("logout")
	// change user role
	r.HandleFunc("/role/{role}",
		authenticatingHandler.Wrap(authenticatingHandler.UpdateRole)).
		Methods("PUT").
		Name("update_role")

	// impersonate user
	r.HandleFunc("/impersonation/add",
		impersonatingHandler.Wrap(impersonatingHandler.AddImpersonation)).
		Methods("GET").
		Name("add_impersonation")
	r.HandleFunc("/impersonation/suggestions",
		impersonatingHandler.Wrap(impersonatingHandler.AddImpersonationSuggest)).
		Methods("GET").
		Name("suggest_impersonations")
	r.HandleFunc("/impersonation",
		impersonatingHandler.Wrap(impersonatingHandler.CreateImpersonation)).
		Methods("POST").
		Name("create_impersonation")
	// TODO why doesn't a DELETE with methodoverride work here?
	r.HandleFunc("/delete-impersonation",
		impersonatingHandler.Wrap(impersonatingHandler.DeleteImpersonation)).
		Methods("POST").
		Name("delete_impersonation")

	// tasks
	// r.HandleFunc("/task/{id}/status", tasksHandler.Wrap(tasksHandler.Status)).
	// 	Methods("GET").
	// 	Name("task_status")

	// dashboard
	r.HandleFunc("/dashboard/publications/{type}", dashboardHandler.Wrap(dashboardHandler.Publications)).
		Methods("GET").
		Name("dashboard_publications")
	r.HandleFunc("/dashboard/datasets/{type}", dashboardHandler.Wrap(dashboardHandler.Datasets)).
		Methods("GET").
		Name("dashboard_datasets")

	// search datasets
	r.HandleFunc("/dataset",
		datasetSearchingHandler.Wrap(datasetSearchingHandler.Search)).
		Methods("GET").
		Name("datasets")

	// add dataset
	r.HandleFunc("/dataset/add",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.Add)).
		Methods("GET", "POST").
		Name("dataset_add")
	r.HandleFunc("/dataset/import",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddImport)).
		Methods("POST").
		Name("dataset_add_import")
	r.HandleFunc("/dataset/import/confirm",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.ConfirmImport)).
		Methods("POST").
		Name("dataset_confirm_import")
	r.HandleFunc("/dataset/{id}/add/description",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddDescription)).
		Methods("GET").
		Name("dataset_add_description")
	r.HandleFunc("/dataset/{id}/add/confirm",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddConfirm)).
		Methods("GET").
		Name("dataset_add_confirm")
	r.HandleFunc("/dataset/{id}/save",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddSaveDraft)).
		Methods("POST").
		Name("dataset_add_save_draft")
	r.HandleFunc("/dataset/{id}/add/publish",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddPublish)).
		Methods("POST").
		Name("dataset_add_publish")
	r.HandleFunc("/dataset/{id}/add/finish",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddFinish)).
		Methods("GET").
		Name("dataset_add_finish")

	// export datasets
	r.HandleFunc("/dataset.{format}",
		datasetExportingHandler.Wrap(datasetExportingHandler.ExportByCurationSearch)).
		Methods("GET").
		Name("export_datasets")

	// view dataset
	r.HandleFunc("/dataset/{id}",
		datasetViewingHandler.Wrap(datasetViewingHandler.Show)).
		Methods("GET").
		Name("dataset")
	r.HandleFunc("/dataset/{id}/description",
		datasetViewingHandler.Wrap(datasetViewingHandler.ShowDescription)).
		Methods("GET").
		Name("dataset_description")
	r.HandleFunc("/dataset/{id}/contributors",
		datasetViewingHandler.Wrap(datasetViewingHandler.ShowContributors)).
		Methods("GET").
		Name("dataset_contributors")
	r.HandleFunc("/dataset/{id}/publications",
		datasetViewingHandler.Wrap(datasetViewingHandler.ShowPublications)).
		Methods("GET").
		Name("dataset_publications")
	r.HandleFunc("/dataset/{id}/activity",
		datasetViewingHandler.Wrap(datasetViewingHandler.ShowActivity)).
		Methods("GET").
		Name("dataset_activity")

	// publish dataset
	r.HandleFunc("/dataset/{id}/publish/confirm",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmPublish)).
		Methods("GET").
		Name("dataset_confirm_publish")
	r.HandleFunc("/dataset/{id}/publish",
		datasetEditingHandler.Wrap(datasetEditingHandler.Publish)).
		Methods("POST").
		Name("dataset_publish")

	// withdraw dataset
	r.HandleFunc("/dataset/{id}/publish/withdraw",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmWithdraw)).
		Methods("GET").
		Name("dataset_confirm_withdraw")
	r.HandleFunc("/dataset/{id}/withdraw",
		datasetEditingHandler.Wrap(datasetEditingHandler.Withdraw)).
		Methods("POST").
		Name("dataset_withdraw")

	// re-publish dataset
	r.HandleFunc("/dataset/{id}/republish/confirm",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmRepublish)).
		Methods("GET").
		Name("dataset_confirm_republish")
	r.HandleFunc("/dataset/{id}/republish",
		datasetEditingHandler.Wrap(datasetEditingHandler.Republish)).
		Methods("POST").
		Name("dataset_republish")

	// lock dataset
	r.HandleFunc("/dataset/{id}/lock",
		datasetEditingHandler.Wrap(datasetEditingHandler.Lock)).
		Methods("POST").
		Name("dataset_lock")
	r.HandleFunc("/dataset/{id}/unlock",
		datasetEditingHandler.Wrap(datasetEditingHandler.Unlock)).
		Methods("POST").
		Name("dataset_unlock")

	// delete dataset
	r.HandleFunc("/dataset/{id}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDelete)).
		Methods("GET").
		Name("dataset_confirm_delete")
	r.HandleFunc("/dataset/{id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.Delete)).
		Methods("DELETE").
		Name("dataset_delete")

	// edit dataset activity
	r.HandleFunc("/dataset/{id}/message/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditMessage)).
		Methods("GET").
		Name("dataset_edit_message")
	r.HandleFunc("/dataset/{id}/message",
		datasetEditingHandler.Wrap(datasetEditingHandler.UpdateMessage)).
		Methods("PUT").
		Name("dataset_update_message")
	r.HandleFunc("/dataset/{id}/reviewer-tags/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditReviewerTags)).
		Methods("GET").
		Name("dataset_edit_reviewer_tags")
	r.HandleFunc("/dataset/{id}/reviewer-tags",
		datasetEditingHandler.Wrap(datasetEditingHandler.UpdateReviewerTags)).
		Methods("PUT").
		Name("dataset_update_reviewer_tags")
	r.HandleFunc("/dataset/{id}/reviewer-note/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditReviewerNote)).
		Methods("GET").
		Name("dataset_edit_reviewer_note")
	r.HandleFunc("/dataset/{id}/reviewer-note",
		datasetEditingHandler.Wrap(datasetEditingHandler.UpdateReviewerNote)).
		Methods("PUT").
		Name("dataset_update_reviewer_note")

	// edit dataset details
	r.HandleFunc("/dataset/{id}/details/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditDetails)).
		Methods("GET").
		Name("dataset_edit_details")
	// r.HandleFunc("/dataset/{id}/details/edit/access-level",
	// 	datasetEditingHandler.Wrap(datasetEditingHandler.EditDetailsAccessLevel)).
	// 	Methods("PUT").
	// 	Name("dataset_edit_details_access_level")
	r.HandleFunc("/dataset/{id}/details/edit/refresh-form",
		datasetEditingHandler.Wrap(datasetEditingHandler.RefreshEditFileForm)).
		Methods("PUT").
		Name("dataset_edit_file_refresh_form")
	r.HandleFunc("/dataset/{id}/details",
		datasetEditingHandler.Wrap(datasetEditingHandler.UpdateDetails)).
		Methods("PUT").
		Name("dataset_update_details")

	// edit dataset projects
	r.HandleFunc("/dataset/{id}/projects/add",
		datasetEditingHandler.Wrap(datasetEditingHandler.AddProject)).
		Methods("GET").
		Name("dataset_add_project")
	r.HandleFunc("/dataset/{id}/projects/suggestions",
		datasetEditingHandler.Wrap(datasetEditingHandler.SuggestProjects)).
		Methods("GET").
		Name("dataset_suggest_projects")
	r.HandleFunc("/dataset/{id}/projects",
		datasetEditingHandler.Wrap(datasetEditingHandler.CreateProject)).
		Methods("POST").
		Name("dataset_create_project")
	r.HandleFunc("/dataset/{id}/{snapshot_id}/projects/confirm-delete/{project_id:.+}",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteProject)).
		Methods("GET").
		Name("dataset_confirm_delete_project")
	r.HandleFunc("/dataset/{id}/projects/{project_id:.+}",
		datasetEditingHandler.Wrap(datasetEditingHandler.DeleteProject)).
		Methods("DELETE").
		Name("dataset_delete_project")

	// edit dataset links
	r.HandleFunc("/dataset/{id}/links/add",
		datasetEditingHandler.Wrap(datasetEditingHandler.AddLink)).
		Methods("GET").
		Name("dataset_add_link")
	r.HandleFunc("/dataset/{id}/links",
		datasetEditingHandler.Wrap(datasetEditingHandler.CreateLink)).
		Methods("POST").
		Name("dataset_create_link")
	r.HandleFunc("/dataset/{id}/links/{link_id}/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditLink)).
		Methods("GET").
		Name("dataset_edit_link")
	r.HandleFunc("/dataset/{id}/links/{link_id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.UpdateLink)).
		Methods("PUT").
		Name("dataset_update_link")
	r.HandleFunc("/dataset/{id}/{snapshot_id}/links/{link_id}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteLink)).
		Methods("GET").
		Name("dataset_confirm_delete_link")
	r.HandleFunc("/dataset/{id}/links/{link_id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.DeleteLink)).
		Methods("DELETE").
		Name("dataset_delete_link")

	// edit dataset departments
	r.HandleFunc("/dataset/{id}/departments/add",
		datasetEditingHandler.Wrap(datasetEditingHandler.AddDepartment)).
		Methods("GET").
		Name("dataset_add_department")
	r.HandleFunc("/dataset/{id}/departments/suggestions",
		datasetEditingHandler.Wrap(datasetEditingHandler.SuggestDepartments)).
		Methods("GET").
		Name("dataset_suggest_departments")
	r.HandleFunc("/dataset/{id}/departments",
		datasetEditingHandler.Wrap(datasetEditingHandler.CreateDepartment)).
		Methods("POST").
		Name("dataset_create_department")
	r.HandleFunc("/dataset/{id}/{snapshot_id}/departments/{department_id}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteDepartment)).
		Methods("GET").
		Name("dataset_confirm_delete_department")
	r.HandleFunc("/dataset/{id}/departments/{department_id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.DeleteDepartment)).
		Methods("DELETE").
		Name("dataset_delete_department")

	// edit dataset abstracts
	r.HandleFunc("/dataset/{id}/abstracts/add",
		datasetEditingHandler.Wrap(datasetEditingHandler.AddAbstract)).
		Methods("GET").
		Name("dataset_add_abstract")
	r.HandleFunc("/dataset/{id}/abstracts",
		datasetEditingHandler.Wrap(datasetEditingHandler.CreateAbstract)).
		Methods("POST").
		Name("dataset_create_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{abstract_id}/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditAbstract)).
		Methods("GET").
		Name("dataset_edit_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{abstract_id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.UpdateAbstract)).
		Methods("PUT").
		Name("dataset_update_abstract")
	r.HandleFunc("/dataset/{id}/{snapshot_id}/abstracts/{abstract_id}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteAbstract)).
		Methods("GET").
		Name("dataset_confirm_delete_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{abstract_id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.DeleteAbstract)).
		Methods("DELETE").
		Name("dataset_delete_abstract")

	// edit dataset publications
	r.HandleFunc("/dataset/{id}/publications/add",
		datasetEditingHandler.Wrap(datasetEditingHandler.AddPublication)).
		Methods("GET").
		Name("dataset_add_publication")
	r.HandleFunc("/dataset/{id}/publications/suggestions",
		datasetEditingHandler.Wrap(datasetEditingHandler.SuggestPublications)).
		Methods("GET").
		Name("dataset_suggest_publications")
	r.HandleFunc("/dataset/{id}/publications",
		datasetEditingHandler.Wrap(datasetEditingHandler.CreatePublication)).
		Methods("POST").
		Name("dataset_create_publication")
	r.HandleFunc("/dataset/{id}/{snapshot_id}/publications/{publication_id}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeletePublication)).
		Methods("GET").
		Name("dataset_confirm_delete_publication")
	r.HandleFunc("/dataset/{id}/publications/{publication_id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.DeletePublication)).
		Methods("DELETE").
		Name("dataset_delete_publication")

	// edit dataset contributors
	r.HandleFunc("/dataset/{id}/contributors/{role}/order",
		datasetEditingHandler.Wrap(datasetEditingHandler.OrderContributors)).
		Methods("POST").
		Name("dataset_order_contributors")
	r.HandleFunc("/dataset/{id}/contributors/{role}/add",
		datasetEditingHandler.Wrap(datasetEditingHandler.AddContributor)).
		Methods("GET").
		Name("dataset_add_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}/suggestions",
		datasetEditingHandler.Wrap(datasetEditingHandler.AddContributorSuggest)).
		Methods("GET").
		Name("dataset_add_contributor_suggest")
	r.HandleFunc("/dataset/{id}/contributors/{role}/confirm-create",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmCreateContributor)).
		Methods("GET").
		Name("dataset_confirm_create_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}",
		datasetEditingHandler.Wrap(datasetEditingHandler.CreateContributor)).
		Methods("POST").
		Name("dataset_create_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditContributor)).
		Methods("GET").
		Name("dataset_edit_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}/suggestions",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditContributorSuggest)).
		Methods("GET").
		Name("dataset_edit_contributor_suggest")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}/confirm-update",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmUpdateContributor)).
		Methods("GET").
		Name("dataset_confirm_update_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}",
		datasetEditingHandler.Wrap(datasetEditingHandler.UpdateContributor)).
		Methods("PUT").
		Name("dataset_update_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteContributor)).
		Methods("GET").
		Name("dataset_confirm_delete_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}",
		datasetEditingHandler.Wrap(datasetEditingHandler.DeleteContributor)).
		Methods("DELETE").
		Name("dataset_delete_contributor")

	// add publication
	r.HandleFunc("/publication/add",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.Add)).
		Methods("GET").
		Name("publication_add")
	r.HandleFunc("/publication/add-single/import",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleImport)).
		Methods("POST").
		Name("publication_add_single_import")
	r.HandleFunc("/publication/add-single/import/confirm",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleImportConfirm)).
		Methods("POST").
		Name("publication_add_single_import_confirm")
	r.HandleFunc("/publication/{id}/add/description",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleDescription)).
		Methods("GET").
		Name("publication_add_single_description")
	r.HandleFunc("/publication/{id}/add/confirm",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleConfirm)).
		Methods("GET").
		Name("publication_add_single_confirm")
	r.HandleFunc("/publication/{id}/add/publish",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSinglePublish)).
		Methods("POST").
		Name("publication_add_single_publish")
	r.HandleFunc("/publication/{id}/add/finish",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleFinish)).
		Methods("GET").
		Name("publication_add_single_finish")
	r.HandleFunc("/publication/add-multiple/import",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleImport)).
		Methods("POST").
		Name("publication_add_multiple_import")
	r.HandleFunc("/publication/add-multiple/{batch_id}/confirm",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleConfirm)).
		Methods("GET").
		Name("publication_add_multiple_confirm")
	r.HandleFunc("/publication/add-multiple/{batch_id}/publication/{id}",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleShow)).
		Methods("GET").
		Name("publication_add_multiple_show")
	r.HandleFunc("/publication/add-multiple/{batch_id}/save",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleSave)).
		Methods("POST").
		Name("publication_add_multiple_save_draft")
	r.HandleFunc("/publication/add-multiple/{batch_id}/publish",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultiplePublish)).
		Methods("POST").
		Name("publication_add_multiple_publish")
	r.HandleFunc("/publication/add-multiple/{batch_id}/finish",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleFinish)).
		Methods("GET").
		Name("publication_add_multiple_finish")

	// search publications
	r.HandleFunc("/publication",
		publicationSearchingHandler.Wrap(publicationSearchingHandler.Search)).
		Methods("GET").
		Name("publications")

	// export publications
	r.HandleFunc("/publication.{format}",
		publicationExportingHandler.Wrap(publicationExportingHandler.ExportByCurationSearch)).
		Methods("GET").
		Name("export_publications")

	// publication batch operations
	r.HandleFunc("/publication/batch",
		publicationBatchHandler.Wrap(publicationBatchHandler.Show)).
		Methods("GET").
		Name("publication_batch")
	r.HandleFunc("/publication/batch",
		publicationBatchHandler.Wrap(publicationBatchHandler.Process)).
		Methods("POST").
		Name("publication_process_batch")

	// view publication
	r.HandleFunc("/publication/{id}",
		publicationViewingHandler.Wrap(publicationViewingHandler.Show)).
		Methods("GET").
		Name("publication")
	r.HandleFunc("/publication/{id}/description",
		publicationViewingHandler.Wrap(publicationViewingHandler.ShowDescription)).
		Methods("GET").
		Name("publication_description")
	r.HandleFunc("/publication/{id}/files",
		publicationViewingHandler.Wrap(publicationViewingHandler.ShowFiles)).
		Methods("GET").
		Name("publication_files")
	r.HandleFunc("/publication/{id}/contributors",
		publicationViewingHandler.Wrap(publicationViewingHandler.ShowContributors)).
		Methods("GET").
		Name("publication_contributors")
	r.HandleFunc("/publication/{id}/datasets",
		publicationViewingHandler.Wrap(publicationViewingHandler.ShowDatasets)).
		Methods("GET").
		Name("publication_datasets")
	r.HandleFunc("/publication/{id}/activity",
		publicationViewingHandler.Wrap(publicationViewingHandler.ShowActivity)).
		Methods("GET").
		Name("publication_activity")
	r.HandleFunc("/publication/{id}/files/{file_id}",
		publicationViewingHandler.Wrap(publicationViewingHandler.DownloadFile)).
		Methods("GET").
		Name("publication_download_file")
	// r.HandleFunc("/publication/{id}/files/{file_id}/thumbnail",
	// 	publicationViewingHandler.Wrap(publicationViewingHandler.FileThumbnail)).
	// 	Methods("GET").
	// 	Name("publication_file_thumbnail")

	// publish publication
	r.HandleFunc("/publication/{id}/publish/confirm",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmPublish)).
		Methods("GET").
		Name("publication_confirm_publish")
	r.HandleFunc("/publication/{id}/publish",
		publicationEditingHandler.Wrap(publicationEditingHandler.Publish)).
		Methods("POST").
		Name("publication_publish")

	// withdraw publication
	r.HandleFunc("/publication/{id}/withdraw/confirm",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmWithdraw)).
		Methods("GET").
		Name("publication_confirm_withdraw")
	r.HandleFunc("/publication/{id}/withdraw",
		publicationEditingHandler.Wrap(publicationEditingHandler.Withdraw)).
		Methods("POST").
		Name("publication_withdraw")

	// re-publish publication
	r.HandleFunc("/publication/{id}/republish/confirm",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmRepublish)).
		Methods("GET").
		Name("publication_confirm_republish")
	r.HandleFunc("/publication/{id}/republish",
		publicationEditingHandler.Wrap(publicationEditingHandler.Republish)).
		Methods("POST").
		Name("publication_republish")

	// lock publication
	r.HandleFunc("/publication/{id}/lock",
		publicationEditingHandler.Wrap(publicationEditingHandler.Lock)).
		Methods("POST").
		Name("publication_lock")
	r.HandleFunc("/publication/{id}/unlock",
		publicationEditingHandler.Wrap(publicationEditingHandler.Unlock)).
		Methods("POST").
		Name("publication_unlock")

	// delete publication
	r.HandleFunc("/publication/{id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDelete)).
		Methods("GET").
		Name("publication_confirm_delete")
	r.HandleFunc("/publication/{id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.Delete)).
		Methods("DELETE").
		Name("publication_delete")

	// edit publication activity
	r.HandleFunc("/publication/{id}/message/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditMessage)).
		Methods("GET").
		Name("publication_edit_message")
	r.HandleFunc("/publication/{id}/message",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateMessage)).
		Methods("PUT").
		Name("publication_update_message")
	r.HandleFunc("/publication/{id}/reviewer-tags/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditReviewerTags)).
		Methods("GET").
		Name("publication_edit_reviewer_tags")
	r.HandleFunc("/publication/{id}/reviewer-tags",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateReviewerTags)).
		Methods("PUT").
		Name("publication_update_reviewer_tags")
	r.HandleFunc("/publication/{id}/reviewer-note/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditReviewerNote)).
		Methods("GET").
		Name("publication_edit_reviewer_note")
	r.HandleFunc("/publication/{id}/reviewer-note",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateReviewerNote)).
		Methods("PUT").
		Name("publication_update_reviewer_note")

	// edit publication details
	r.HandleFunc("/publication/{id}/details/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditDetails)).
		Methods("GET").
		Name("publication_edit_details")
	r.HandleFunc("/publication/{id}/details",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateDetails)).
		Methods("PUT").
		Name("publication_update_details")

	// edit publication type
	r.HandleFunc("/publication/{id}/type/confirm",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmUpdateType)).
		Methods("GET").
		Name("publication_confirm_update_type")
	r.HandleFunc("/publication/{id}/type",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateType)).
		Methods("PUT").
		Name("publication_update_type")

	// edit publication conference
	r.HandleFunc("/publication/{id}/conference/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditConference)).
		Methods("GET").
		Name("publication_edit_conference")
	r.HandleFunc("/publication/{id}/conference",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateConference)).
		Methods("PUT").
		Name("publication_update_conference")

	// edit publication additional info
	r.HandleFunc("/publication/{id}/additional-info/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditAdditionalInfo)).
		Methods("GET").
		Name("publication_edit_additional_info")
	r.HandleFunc("/publication/{id}/additional-info",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateAdditionalInfo)).
		Methods("PUT").
		Name("publication_update_additional_info")

	// edit publication projects
	r.HandleFunc("/publication/{id}/projects/add",
		publicationEditingHandler.Wrap(publicationEditingHandler.AddProject)).
		Methods("GET").
		Name("publication_add_project")
	r.HandleFunc("/publication/{id}/projects/suggestions",
		publicationEditingHandler.Wrap(publicationEditingHandler.SuggestProjects)).
		Methods("GET").
		Name("publication_suggest_projects")
	r.HandleFunc("/publication/{id}/projects",
		publicationEditingHandler.Wrap(publicationEditingHandler.CreateProject)).
		Methods("POST").
		Name("publication_create_project")
	// project_id is last part of url because some id's contain slashes
	r.HandleFunc("/publication/{id}/{snapshot_id}/projects/confirm-delete/{project_id:.+}",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteProject)).
		Methods("GET").
		Name("publication_confirm_delete_project")
	// project_id is last part of url because some id's contain slashes
	r.HandleFunc("/publication/{id}/projects/{project_id:.+}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteProject)).
		Methods("DELETE").
		Name("publication_delete_project")

	// edit publication links
	r.HandleFunc("/publicaton/{id}/links/add",
		publicationEditingHandler.Wrap(publicationEditingHandler.AddLink)).
		Methods("GET").
		Name("publication_add_link")
	r.HandleFunc("/publication/{id}/links",
		publicationEditingHandler.Wrap(publicationEditingHandler.CreateLink)).
		Methods("POST").
		Name("publication_create_link")
	r.HandleFunc("/publication/{id}/links/{link_id}/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditLink)).
		Methods("GET").
		Name("publication_edit_link")
	r.HandleFunc("/publication/{id}/links/{link_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateLink)).
		Methods("PUT").
		Name("publication_update_link")
	r.HandleFunc("/publication/{id}/{snapshot_id}/links/{link_id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteLink)).
		Methods("GET").
		Name("publication_confirm_delete_link")
	r.HandleFunc("/publication/{id}/links/{link_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteLink)).
		Methods("DELETE").
		Name("publication_delete_link")

	// edit publication departments
	r.HandleFunc("/publication/{id}/departments/add",
		publicationEditingHandler.Wrap(publicationEditingHandler.AddDepartment)).
		Methods("GET").
		Name("publication_add_department")
	r.HandleFunc("/publication/{id}/departments/suggestions",
		publicationEditingHandler.Wrap(publicationEditingHandler.SuggestDepartments)).
		Methods("GET").
		Name("publication_suggest_departments")
	r.HandleFunc("/publication/{id}/departments",
		publicationEditingHandler.Wrap(publicationEditingHandler.CreateDepartment)).
		Methods("POST").
		Name("publication_create_department")
	r.HandleFunc("/publication/{id}/{snapshot_id}/departments/{department_id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteDepartment)).
		Methods("GET").
		Name("publication_confirm_delete_department")
	r.HandleFunc("/publication/{id}/departments/{department_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteDepartment)).
		Methods("DELETE").
		Name("publication_delete_department")

	// edit publication abstracts
	r.HandleFunc("/publication/{id}/abstracts/add",
		publicationEditingHandler.Wrap(publicationEditingHandler.AddAbstract)).
		Methods("GET").
		Name("publication_add_abstract")
	r.HandleFunc("/publication/{id}/abstracts",
		publicationEditingHandler.Wrap(publicationEditingHandler.CreateAbstract)).
		Methods("POST").
		Name("publication_create_abstract")
	r.HandleFunc("/publication/{id}/abstracts/{abstract_id}/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditAbstract)).
		Methods("GET").
		Name("publication_edit_abstract")
	r.HandleFunc("/publication/{id}/abstracts/{abstract_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateAbstract)).
		Methods("PUT").
		Name("publication_update_abstract")
	r.HandleFunc("/publication/{id}/{snapshot_id}/abstracts/{abstract_id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteAbstract)).
		Methods("GET").
		Name("publication_confirm_delete_abstract")
	r.HandleFunc("/publication/{id}/abstracts/{abstract_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteAbstract)).
		Methods("DELETE").
		Name("publication_delete_abstract")

	// edit publication lay summaries
	r.HandleFunc("/publication/{id}/lay_summaries/add",
		publicationEditingHandler.Wrap(publicationEditingHandler.AddLaySummary)).
		Methods("GET").
		Name("publication_add_lay_summary")
	r.HandleFunc("/publication/{id}/lay_summaries",
		publicationEditingHandler.Wrap(publicationEditingHandler.CreateLaySummary)).
		Methods("POST").
		Name("publication_create_lay_summary")
	r.HandleFunc("/publication/{id}/lay_summaries/{lay_summary_id}/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditLaySummary)).
		Methods("GET").
		Name("publication_edit_lay_summary")
	r.HandleFunc("/publication/{id}/lay_summaries/{lay_summary_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateLaySummary)).
		Methods("PUT").
		Name("publication_update_lay_summary")
	r.HandleFunc("/publication/{id}/{snapshot_id}/lay_summaries/{lay_summary_id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteLaySummary)).
		Methods("GET").
		Name("publication_confirm_delete_lay_summary")
	r.HandleFunc("/publication/{id}/lay_summaries/{lay_summary_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteLaySummary)).
		Methods("DELETE").
		Name("publication_delete_lay_summary")

	// edit publication datasets
	r.HandleFunc("/publication/{id}/datasets/add",
		publicationEditingHandler.Wrap(publicationEditingHandler.AddDataset)).
		Methods("GET").
		Name("publication_add_dataset")
	r.HandleFunc("/publication/{id}/datasets/suggestions",
		publicationEditingHandler.Wrap(publicationEditingHandler.SuggestDatasets)).
		Methods("GET").
		Name("publication_suggest_datasets")
	r.HandleFunc("/publication/{id}/datasets",
		publicationEditingHandler.Wrap(publicationEditingHandler.CreateDataset)).
		Methods("POST").
		Name("publication_create_dataset")
	r.HandleFunc("/publication/{id}/{snapshot_id}/datasets/{dataset_id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteDataset)).
		Methods("GET").
		Name("publication_confirm_delete_dataset")
	r.HandleFunc("/publication/{id}/datasets/{dataset_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteDataset)).
		Methods("DELETE").
		Name("publication_delete_dataset")

	// edit publication contributors
	r.HandleFunc("/publication/{id}/contributors/{role}/order",
		publicationEditingHandler.Wrap(publicationEditingHandler.OrderContributors)).
		Methods("POST").
		Name("publication_order_contributors")
	r.HandleFunc("/publication/{id}/contributors/{role}/add",
		publicationEditingHandler.Wrap(publicationEditingHandler.AddContributor)).
		Methods("GET").
		Name("publication_add_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}/suggestions",
		publicationEditingHandler.Wrap(publicationEditingHandler.AddContributorSuggest)).
		Methods("GET").
		Name("publication_add_contributor_suggest")
	r.HandleFunc("/publication/{id}/contributors/{role}/confirm-create",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmCreateContributor)).
		Methods("GET").
		Name("publication_confirm_create_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}",
		publicationEditingHandler.Wrap(publicationEditingHandler.CreateContributor)).
		Methods("POST").
		Name("publication_create_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditContributor)).
		Methods("GET").
		Name("publication_edit_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}/suggestions",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditContributorSuggest)).
		Methods("GET").
		Name("publication_edit_contributor_suggest")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}/confirm-update",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmUpdateContributor)).
		Methods("GET").
		Name("publication_confirm_update_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateContributor)).
		Methods("PUT").
		Name("publication_update_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteContributor)).
		Methods("GET").
		Name("publication_confirm_delete_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteContributor)).
		Methods("DELETE").
		Name("publication_delete_contributor")

	// edit publication files
	r.HandleFunc("/publication/{id}/files",
		publicationEditingHandler.Wrap(publicationEditingHandler.UploadFile)).
		Methods("POST").
		Name("publication_upload_file")
	r.HandleFunc("/publication/{id}/files/{file_id}/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditFile)).
		Methods("GET").
		Name("publication_edit_file")
	r.HandleFunc("/publication/{id}/refresh-files",
		publicationEditingHandler.Wrap(publicationEditingHandler.RefreshFiles)).
		Methods("GET").
		Name("publication_refresh_files")
	r.HandleFunc("/publication/{id}/files/{file_id}/refresh-form",
		publicationEditingHandler.Wrap(publicationEditingHandler.RefreshEditFileForm)).
		Methods("GET").
		Name("publication_edit_file_refresh_form")
	r.HandleFunc("/publication/{id}/files/{file_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateFile)).
		Methods("PUT").
		Name("publication_update_file")
	r.HandleFunc("/publication/{id}/{snapshot_id}/files/{file_id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteFile)).
		Methods("GET").
		Name("publication_confirm_delete_file")
	r.HandleFunc("/publication/{id}/files/{file_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteFile)).
		Methods("DELETE").
		Name("publication_delete_file")

	// media types
	r.HandleFunc("/media_type/suggestions",
		mediaTypesHandler.Wrap(mediaTypesHandler.Suggest)).
		Methods("GET").
		Name("suggest_media_types")

	r.NotFoundHandler = baseHandler.Wrap(baseHandler.NotFound)
}
