package routes

import (
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	mw "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/authenticating"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/datasetcreating"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/datasetediting"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/datasetsearching"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/datasetviewing"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/home"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/impersonating"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/mediatypes"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/orcid"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/publicationcreating"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/publicationediting"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/publicationsearching"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/publicationviewing"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/tasks"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/go-oidc/oidc"
)

func Register(services *backends.Services, baseURL *url.URL, router *mux.Router,
	sessionStore sessions.Store, sessionName string, localizer *locale.Localizer, oidcClient *oidc.Client) {
	basePath := baseURL.Path

	router.StrictSlash(true)
	router.UseEncodedPath()
	router.Use(mw.RecoveryHandler(mw.PrintRecoveryStack(true)))

	// static files
	router.PathPrefix(basePath + "/static/").Handler(http.StripPrefix(basePath+"/static/", http.FileServer(http.Dir("./static"))))

	baseHandler := handlers.BaseHandler{
		Router:       router,
		SessionStore: sessionStore,
		SessionName:  sessionName,
		Localizer:    localizer,
		UserService:  services.UserService,
	}
	homeHandler := &home.Handler{
		BaseHandler: baseHandler,
	}
	authenticatingHandler := &authenticating.Handler{
		BaseHandler: baseHandler,
		OIDCClient:  oidcClient,
	}
	impersonatingHandler := &impersonating.Handler{
		BaseHandler: baseHandler,
	}
	tasksHandler := &tasks.Handler{
		BaseHandler: baseHandler,
		Tasks:       services.Tasks,
	}
	datasetSearchingHandler := &datasetsearching.Handler{
		BaseHandler:          baseHandler,
		DatasetSearchService: services.DatasetSearchService,
	}
	datasetViewingHandler := &datasetviewing.Handler{
		BaseHandler: baseHandler,
		Repository:  services.Repository,
	}
	datasetCreatingHandler := &datasetcreating.Handler{
		BaseHandler:          baseHandler,
		Repository:           services.Repository,
		DatasetSearchService: services.DatasetSearchService,
		DatasetSources:       services.DatasetSources,
	}
	datasetEditingHandler := &datasetediting.Handler{
		BaseHandler:               baseHandler,
		Repository:                services.Repository,
		ProjectService:            services.ProjectService,
		ProjectSearchService:      services.ProjectSearchService,
		OrganizationSearchService: services.OrganizationSearchService,
		OrganizationService:       services.OrganizationService,
		PersonSearchService:       services.PersonSearchService,
		PersonService:             services.PersonService,
		PublicationSearchService:  services.PublicationSearchService,
	}
	publicationSearchingHandler := &publicationsearching.Handler{
		BaseHandler:              baseHandler,
		PublicationSearchService: services.PublicationSearchService,
		FileStore:                services.FileStore,
	}
	publicationViewingHandler := &publicationviewing.Handler{
		BaseHandler: baseHandler,
		Repository:  services.Repository,
	}
	publicationCreatingHandler := &publicationcreating.Handler{
		BaseHandler:              baseHandler,
		Repository:               services.Repository,
		PublicationSearchService: services.PublicationSearchService,
		PublicationSources:       services.PublicationSources,
		PublicationDecoders:      services.PublicationDecoders,
	}
	publicationEditingHandler := &publicationediting.Handler{
		BaseHandler:               baseHandler,
		Repository:                services.Repository,
		ProjectService:            services.ProjectService,
		ProjectSearchService:      services.ProjectSearchService,
		OrganizationSearchService: services.OrganizationSearchService,
		OrganizationService:       services.OrganizationService,
		PersonSearchService:       services.PersonSearchService,
		PersonService:             services.PersonService,
		DatasetSearchService:      services.DatasetSearchService,
		FileStore:                 services.FileStore,
	}
	orcidHandler := &orcid.Handler{
		BaseHandler:              baseHandler,
		Tasks:                    services.Tasks,
		Repository:               services.Repository,
		PublicationSearchService: services.PublicationSearchService,
		Sandbox:                  services.ORCIDSandbox,
	}
	mediaTypesHandler := &mediatypes.Handler{
		BaseHandler:            baseHandler,
		MediaTypeSearchService: services.MediaTypeSearchService,
	}

	// TODO fix absolute url generation
	// var schemes []string
	// if u.Scheme == "http" {
	// 	schemes = []string{"http", "https"}
	// } else {
	// 	schemes = []string{"https", "http"}
	// }
	// r = r.Schemes(schemes...).Host(u.Host).PathPrefix(u.Path).Subrouter()

	r := router.PathPrefix(basePath).Subrouter()
	r.Use(csrf.Protect(
		[]byte(viper.GetString("csrf-secret")),        // TODO pass as argument
		csrf.CookieName(viper.GetString("csrf-name")), // TODO pass as argument
		csrf.Path(basePath),
		csrf.Secure(baseURL.Scheme == "https"),
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

	// impersonate user
	r.HandleFunc("/impersonation/add",
		impersonatingHandler.Wrap(impersonatingHandler.AddImpersonation)).
		Methods("GET").
		Name("add_impersonation")
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
	r.HandleFunc("/task/{id}/status", tasksHandler.Wrap(tasksHandler.Status)).
		Methods("GET").
		Name("task_status")

	// search datasets
	r.HandleFunc("/curation/dataset",
		datasetSearchingHandler.Wrap(datasetSearchingHandler.CurationSearch)).
		Methods("GET").
		Name("cureation_datasets")
	r.HandleFunc("/dataset",
		datasetSearchingHandler.Wrap(datasetSearchingHandler.Search)).
		Methods("GET").
		Name("datasets")

	// add dataset
	r.HandleFunc("/dataset/add",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.Add)).
		Methods("GET").
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
	r.HandleFunc("/dataset/{id}/add/publish",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddPublish)).
		Methods("POST").
		Name("dataset_add_publish")
	r.HandleFunc("/dataset/{id}/add/finish",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddFinish)).
		Methods("GET").
		Name("dataset_add_finish")

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

	// publish dataset
	r.HandleFunc("/dataset/{id}/publish/confirm",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmPublish)).
		Methods("GET").
		Name("dataset_confirm_publish")
	r.HandleFunc("/dataset/{id}/publish",
		datasetEditingHandler.Wrap(datasetEditingHandler.Publish)).
		Methods("POST").
		Name("dataset_publish")

	// delete dataset
	r.HandleFunc("/dataset/{id}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDelete)).
		Methods("GET").
		Name("dataset_confirm_delete")
	r.HandleFunc("/dataset/{id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.Delete)).
		Methods("DELETE").
		Name("dataset_delete")

	// edit dataset details
	r.HandleFunc("/dataset/{id}/details/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditDetails)).
		Methods("GET").
		Name("dataset_edit_details")
	r.HandleFunc("/dataset/{id}/details/edit/access-level",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditDetailsAccessLevel)).
		Methods("PUT").
		Name("dataset_edit_details_access_level")
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
	// TODO: weird project_id's not problematic in urls?
	r.HandleFunc("/dataset/{id}/projects/{project_id}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteProject)).
		Methods("GET").
		Name("dataset_confirm_delete_project")
	r.HandleFunc("/dataset/{id}/projects/{project_id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.DeleteProject)).
		Methods("DELETE").
		Name("dataset_delete_project")

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
	r.HandleFunc("/dataset/{id}/departments/{department_id}/confirm-delete",
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
	r.HandleFunc("/dataset/{id}/abstracts/{abstract_id}/confirm-delete",
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
	r.HandleFunc("/dataset/{id}/publications/{publication_id}/confirm-delete",
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
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}/suggestions",
		datasetEditingHandler.Wrap(datasetEditingHandler.SuggestContributors)).
		Methods("GET").
		Name("dataset_suggest_contributors")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}/confirm",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmContributor)).
		Methods("POST").
		Name("dataset_confirm_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}/unconfirm",
		datasetEditingHandler.Wrap(datasetEditingHandler.UnconfirmContributor)).
		Methods("POST").
		Name("dataset_unconfirm_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}",
		datasetEditingHandler.Wrap(datasetEditingHandler.CreateContributor)).
		Methods("POST").
		Name("dataset_create_contributor")
	r.HandleFunc("/dataset/{id}/contributors/{role}/{position}/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditContributor)).
		Methods("GET").
		Name("dataset_edit_contributor")
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
	r.HandleFunc("/publication/add-multiple/{batch_id}/description",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleDescription)).
		Methods("GET").
		Name("publication_add_multiple_description")
	r.HandleFunc("/publication/add-multiple/{batch_id}/confirm",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleConfirm)).
		Methods("GET").
		Name("publication_add_multiple_confirm")
	r.HandleFunc("/publication/add-multiple/{batch_id}/publication/{id}",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleShow)).
		Methods("GET").
		Name("publication_add_multiple_show")
	r.HandleFunc("/publication/add-multiple/{batch_id}/publish",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultiplePublish)).
		Methods("POST").
		Name("publication_add_multiple_publish")
	r.HandleFunc("/publication/add-multiple/{batch_id}/finish",
		publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleFinish)).
		Methods("GET").
		Name("publication_add_multiple_finish")

	// search publications
	r.HandleFunc("/curation/publication",
		publicationSearchingHandler.Wrap(publicationSearchingHandler.CurationSearch)).
		Methods("GET").
		Name("cureation_publications")
	r.HandleFunc("/publication",
		publicationSearchingHandler.Wrap(publicationSearchingHandler.Search)).
		Methods("GET").
		Name("publications")

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

	// publish publication
	r.HandleFunc("/publication/{id}/publish/confirm",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmPublish)).
		Methods("GET").
		Name("publication_confirm_publish")
	r.HandleFunc("/publication/{id}/publish",
		publicationEditingHandler.Wrap(publicationEditingHandler.Publish)).
		Methods("POST").
		Name("publication_publish")

	// delete publication
	r.HandleFunc("/publication/{id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDelete)).
		Methods("GET").
		Name("publication_confirm_delete")
	r.HandleFunc("/publication/{id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.Delete)).
		Methods("DELETE").
		Name("publication_delete")

	// edit publication details
	r.HandleFunc("/publication/{id}/details/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditDetails)).
		Methods("GET").
		Name("publication_edit_details")
	r.HandleFunc("/publication/{id}/details",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateDetails)).
		Methods("PUT").
		Name("publication_update_details")

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
	r.HandleFunc("/publication/{id}/projects/{project_id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteProject)).
		Methods("GET").
		Name("publication_confirm_delete_project")
	r.HandleFunc("/publication/{id}/projects/{project_id}",
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
	r.HandleFunc("/publication/{id}/links/{link_id}/confirm-delete",
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
	r.HandleFunc("/publication/{id}/departments/{department_id}/confirm-delete",
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
	r.HandleFunc("/publication/{id}/abstracts/{abstract_id}/confirm-delete",
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
	r.HandleFunc("/publication/{id}/lay_summaries/{lay_summary_id}/confirm-delete",
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
	r.HandleFunc("/publication/{id}/datasets/{dataset_id}/confirm-delete",
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
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}/suggestions",
		publicationEditingHandler.Wrap(publicationEditingHandler.SuggestContributors)).
		Methods("GET").
		Name("publication_suggest_contributors")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}/confirm",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmContributor)).
		Methods("POST").
		Name("publication_confirm_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}/unconfirm",
		publicationEditingHandler.Wrap(publicationEditingHandler.UnconfirmContributor)).
		Methods("POST").
		Name("publication_unconfirm_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}",
		publicationEditingHandler.Wrap(publicationEditingHandler.CreateContributor)).
		Methods("POST").
		Name("publication_create_contributor")
	r.HandleFunc("/publication/{id}/contributors/{role}/{position}/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditContributor)).
		Methods("GET").
		Name("publication_edit_contributor")
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
	r.HandleFunc("/publication/{id}/files/{file_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DownloadFile)).
		Methods("GET").
		Name("publication_download_file")
	r.HandleFunc("/publication/{id}/files/{file_id}/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditFile)).
		Methods("GET").
		Name("publication_edit_file")
	r.HandleFunc("/publication/{id}/files/{file_id}/edit/license",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditFileLicense)).
		Methods("PUT").
		Name("publication_edit_file_license")
	r.HandleFunc("/publication/{id}/files/{file_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateFile)).
		Methods("PUT").
		Name("publication_update_file")
	r.HandleFunc("/publication/{id}/files/{file_id}/confirm-delete",
		publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmDeleteFile)).
		Methods("GET").
		Name("publication_confirm_delete_file")
	r.HandleFunc("/publication/{id}/files/{file_id}",
		publicationEditingHandler.Wrap(publicationEditingHandler.DeleteFile)).
		Methods("DELETE").
		Name("publication_delete_file")

	// orcid
	r.HandleFunc("/publication/orcid",
		orcidHandler.Wrap(orcidHandler.AddAll)).
		Methods("POST").
		Name("publication_orcid_add_all")
	r.HandleFunc("/publication/{id}/orcid",
		orcidHandler.Wrap(orcidHandler.Add)).
		Methods("POST").
		Name("publication_orcid_add")

	// media types
	r.HandleFunc("/media_type/suggestions",
		mediaTypesHandler.Wrap(mediaTypesHandler.Suggest)).
		Methods("GET").
		Name("suggest_media_types")
}
