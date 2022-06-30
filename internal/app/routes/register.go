package routes

import (
	"net/http"

	"github.com/gorilla/csrf"
	mw "github.com/gorilla/handlers"
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
	"github.com/ugent-library/biblio-backend/internal/app/handlers/publicationediting"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/publicationviewing"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/tasks"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/controllers"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/middleware"
	"github.com/ugent-library/go-oidc/oidc"
)

func Register(services *backends.Services, oldBase controllers.Base, oidcClient *oidc.Client) {
	router := oldBase.Router
	basePath := oldBase.BaseURL.Path

	router.StrictSlash(true)
	router.UseEncodedPath()
	router.Use(mw.RecoveryHandler(mw.PrintRecoveryStack(true)))

	// static files
	router.PathPrefix(basePath + "/static/").Handler(http.StripPrefix(basePath+"/static/", http.FileServer(http.Dir("./static"))))

	requireUser := middleware.RequireUser(oldBase.BaseURL.Path + "/login")
	setUser := middleware.SetUser(services.UserService, oldBase.SessionName, oldBase.SessionStore)

	publicationsController := controllers.NewPublications(
		oldBase,
		services.Repository,
		services.FileStore,
		services.PublicationSearchService,
		services.PublicationDecoders,
		services.PublicationSources,
		services.Tasks,
		services.ORCIDSandbox,
	)
	publicationFilesController := controllers.NewPublicationFiles(oldBase, services.Repository, services.FileStore)
	publicationDetailsController := controllers.NewPublicationDetails(oldBase, services.Repository)
	publicationConferenceController := controllers.NewPublicationConference(oldBase, services.Repository)
	publicationProjectsController := controllers.NewPublicationProjects(oldBase, services.Repository, services.ProjectSearchService, services.ProjectService)
	publicationDepartmentsController := controllers.NewPublicationDepartments(oldBase, services.Repository, services.OrganizationSearchService, services.OrganizationService)
	publicationAbstractsController := controllers.NewPublicationAbstracts(oldBase, services.Repository)
	publicationLinksController := controllers.NewPublicationLinks(oldBase, services.Repository)
	publicationContributorsController := controllers.NewPublicationContributors(oldBase, services.Repository, services.PersonSearchService, services.PersonService)
	publicationDatasetsController := controllers.NewPublicationDatasets(oldBase, services.Repository, services.DatasetSearchService)
	publicationAdditionalInfoController := controllers.NewPublicationAdditionalInfo(oldBase, services.Repository)
	publicationLaySummariesController := controllers.NewPublicationLaySummaries(oldBase, services.Repository)

	// NEW HANDLERS
	baseHandler := handlers.BaseHandler{
		Router:       oldBase.Router,
		SessionStore: oldBase.SessionStore,
		SessionName:  oldBase.SessionName,
		Localizer:    oldBase.Localizer,
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
	publicationViewingHandler := &publicationviewing.Handler{
		BaseHandler: baseHandler,
		Repository:  services.Repository,
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
		PublicationSearchService:  services.PublicationSearchService,
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
		[]byte(viper.GetString("csrf-secret")),
		csrf.CookieName(viper.GetString("csrf-name")),
		csrf.Path(basePath),
		csrf.Secure(oldBase.BaseURL.Scheme == "https"),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.FieldName("csrf-token"),
	))

	// NEW ROUTES
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

	// Add a dataset
	r.HandleFunc("/dataset/add",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.Add)).
		Methods("GET").
		Name("dataset_add")
	r.HandleFunc("/dataset/import",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.AddImport)).
		Methods("POST").
		Name("dataset_add_import")
	r.HandleFunc("/dataset/import/confirm",
		datasetCreatingHandler.Wrap(datasetCreatingHandler.ConfirmImportDataset)).
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

	// tasks
	r.HandleFunc("/task/{id}/status", tasksHandler.Wrap(tasksHandler.Status)).
		Methods("GET").
		Name("task_status")

	// search datasets
	r.HandleFunc("/dataset",
		datasetSearchingHandler.Wrap(datasetSearchingHandler.Search)).
		Methods("GET").
		Name("datasets")

	// dataset delete
	r.HandleFunc("/dataset/{id}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDelete)).
		Methods("GET").
		Name("dataset_confirm_delete")
	r.HandleFunc("/dataset/{id}",
		datasetEditingHandler.Wrap(datasetEditingHandler.Delete)).
		Methods("DELETE").
		Name("dataset_delete")

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
	r.HandleFunc("/dataset/{id}/projects/{position}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteProject)).
		Methods("GET").
		Name("dataset_confirm_delete_project")
	r.HandleFunc("/dataset/{id}/projects/{position}",
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
	r.HandleFunc("/dataset/{id}/departments/{position}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteDepartment)).
		Methods("GET").
		Name("dataset_confirm_delete_department")
	r.HandleFunc("/dataset/{id}/departments/{position}",
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
	r.HandleFunc("/dataset/{id}/abstracts/{position}/edit",
		datasetEditingHandler.Wrap(datasetEditingHandler.EditAbstract)).
		Methods("GET").
		Name("dataset_edit_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{position}",
		datasetEditingHandler.Wrap(datasetEditingHandler.UpdateAbstract)).
		Methods("PUT").
		Name("dataset_update_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{position}/confirm-delete",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteAbstract)).
		Methods("GET").
		Name("dataset_confirm_delete_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{position}",
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

	// publish dataset
	r.HandleFunc("/dataset/{id}/publish/confirm",
		datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmPublish)).
		Methods("GET").
		Name("dataset_confirm_publish")
	r.HandleFunc("/dataset/{id}/publish",
		datasetEditingHandler.Wrap(datasetEditingHandler.Publish)).
		Methods("POST").
		Name("dataset_publish")

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

	// edit publication details
	r.HandleFunc("/publication/{id}/details/edit",
		publicationEditingHandler.Wrap(publicationEditingHandler.EditDetails)).
		Methods("GET").
		Name("publication_edit_details")
	r.HandleFunc("/publication/{id}/details",
		publicationEditingHandler.Wrap(publicationEditingHandler.UpdateDetails)).
		Methods("PUT").
		Name("publication_update_details")

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

	// r.Use(handlers.HTTPMethodOverrideHandler)
	r.Use(locale.Detect(oldBase.Localizer))

	r.Use(setUser)

	// publications
	pubsRouter := r.PathPrefix("/publication").Subrouter()
	pubsRouter.Use(middleware.SetActiveMenu("publications"))
	pubsRouter.Use(requireUser)
	pubsRouter.HandleFunc("", publicationsController.List).
		Methods("GET").
		Name("publications")
	pubsRouter.HandleFunc("/add", publicationsController.Add).
		Methods("GET").
		Name("publication_add")
	pubsRouter.HandleFunc("/add", publicationsController.AddSelectMethod).
		Methods("POST").
		Name("publication_add_select_method")
	pubsRouter.HandleFunc("/add-single/import/confirm", publicationsController.AddSingleImportConfirm).
		Methods("POST").
		Name("publication_add_single_import_confirm")
	pubsRouter.HandleFunc("/add-single/import", publicationsController.AddSingleImport).
		Methods("POST").
		Name("publication_add_single_import")
	pubsRouter.HandleFunc("/add-multiple/import", publicationsController.AddMultipleImport).
		Methods("POST").
		Name("publication_add_multiple_import")
	pubsRouter.HandleFunc("/add-multiple/{batch_id}/description", publicationsController.AddMultipleDescription).
		Methods("GET").
		Name("publication_add_multiple_description")
	pubsRouter.HandleFunc("/add-multiple/{batch_id}/confirm", publicationsController.AddMultipleConfirm).
		Methods("GET").
		Name("publication_add_multiple_confirm")
	pubsRouter.HandleFunc("/add-multiple/{batch_id}/publish", publicationsController.AddMultiplePublish).
		Methods("POST").
		Name("publication_add_multiple_publish")

	pubRouter := pubsRouter.PathPrefix("/{id}").Subrouter()
	pubRouter.Use(middleware.SetPublication(services.Repository))
	pubRouter.Use(middleware.RequireCanViewPublication)
	pubEditRouter := pubRouter.PathPrefix("").Subrouter()
	pubEditRouter.Use(middleware.RequireCanEditPublication)
	pubPublishRouter := pubRouter.PathPrefix("").Subrouter()
	pubPublishRouter.Use(middleware.RequireCanPublishPublication)
	pubDeleteRouter := pubRouter.PathPrefix("").Subrouter()
	pubDeleteRouter.Use(middleware.RequireCanDeletePublication)
	pubRouter.HandleFunc("/delete", publicationsController.ConfirmDelete).
		Methods("GET").
		Name("publication_confirm_delete")
	// TODO why doesn't a DELETE with methodoverride work with CAS?
	pubDeleteRouter.HandleFunc("/delete", publicationsController.Delete).
		Methods("POST").
		Name("publication_delete")
	pubPublishRouter.HandleFunc("/publish", publicationsController.Publish).
		Methods("POST").
		Name("publication_publish")
	pubEditRouter.HandleFunc("/add-single/description", publicationsController.AddSingleDescription).
		Methods("GET").
		Name("publication_add_single_description")
	pubEditRouter.HandleFunc("/add-single/confirm", publicationsController.AddSingleConfirm).
		Methods("GET").
		Name("publication_add_single_confirm")
	pubEditRouter.HandleFunc("/add-single/publish", publicationsController.AddSinglePublish).
		Methods("POST").
		Name("publication_add_single_publish")
	pubRouter.HandleFunc("/add-multiple/{batch_id}", publicationsController.AddMultipleShow).
		Methods("GET").
		Name("publication_add_multiple_show")
	pubRouter.HandleFunc("/add-multiple/{batch_id}/confirm", publicationsController.AddMultipleConfirmShow).
		Methods("GET").
		Name("publication_add_multiple_confirm_show")
	// Publication files
	pubRouter.HandleFunc("/file/{file_id}", publicationFilesController.Download).
		Methods("GET").
		Name("publication_file")
	pubEditRouter.HandleFunc("/htmx/file", publicationFilesController.Upload).
		Methods("POST").
		Name("upload_publication_file")
	pubEditRouter.HandleFunc("/htmx/file/{file_id}/edit", publicationFilesController.Edit).
		Methods("GET").
		Name("publication_file_edit")
	pubEditRouter.HandleFunc("/htmx/file/{file_id}/license", publicationFilesController.License).
		Methods("PUT").
		Name("publication_file_license")
	pubEditRouter.HandleFunc("/htmx/file/{file_id}", publicationFilesController.Update).
		Methods("PUT").
		Name("publication_file_update")
	pubEditRouter.HandleFunc("/htmx/file/{file_id}/remove", publicationFilesController.ConfirmRemove).
		Methods("GET").
		Name("publication_file_confirm_remove")
	pubEditRouter.HandleFunc("/htmx/file/{file_id}/remove", publicationFilesController.Remove).
		Methods("PATCH").
		Name("publication_file_remove")
	// Publication HTMX fragments
	pubEditRouter.HandleFunc("/htmx/summary", publicationsController.Summary).
		Methods("GET").
		Name("publication_summary")
	// Publication details HTMX fragments
	pubEditRouter.HandleFunc("/htmx", publicationDetailsController.Show).
		Methods("GET").
		Name("publication_details")
	pubEditRouter.HandleFunc("/htmx/edit", publicationDetailsController.Edit).
		Methods("GET").
		Name("publication_details_edit_form")
	pubEditRouter.HandleFunc("/htmx/edit", publicationDetailsController.Update).
		Methods("PATCH").
		Name("publication_details_save_form")
	// Publication conference HTMX fragments
	pubEditRouter.HandleFunc("/htmx/conference", publicationConferenceController.Show).
		Methods("GET").
		Name("publication_conference")
	pubEditRouter.HandleFunc("/htmx/conference/edit", publicationConferenceController.Edit).
		Methods("GET").
		Name("publication_conference_edit_form")
	pubEditRouter.HandleFunc("/htmx/conference/edit", publicationConferenceController.Update).
		Methods("PATCH").
		Name("publication_conference_save_form")
	// Publication additional info HTMX fragments
	pubEditRouter.HandleFunc("/htmx/additional_info", publicationAdditionalInfoController.Show).
		Methods("GET").
		Name("publication_additional_info")
	pubEditRouter.HandleFunc("/htmx/additional_info/edit", publicationAdditionalInfoController.Edit).
		Methods("GET").
		Name("publication_additional_info_edit_form")
	pubEditRouter.HandleFunc("/htmx/additional_info/edit", publicationAdditionalInfoController.Update).
		Methods("PATCH").
		Name("publication_additional_info_save_form")
	// Publication projects HTMX fragments
	pubEditRouter.HandleFunc("/htmx/projects/list", publicationProjectsController.List).
		Methods("GET").
		Name("publication_add_project")
	pubEditRouter.HandleFunc("/htmx/projects/list/activesearch", publicationProjectsController.ActiveSearch).
		Methods("POST").
		Name("publication_projects_activesearch")
	pubEditRouter.HandleFunc("/htmx/projects/add/{project_id:[a-zA-Z0-9].*}", publicationProjectsController.Add).
		Methods("PATCH").
		Name("publication_projects_add_to_publication")
	pubEditRouter.HandleFunc("/htmx/projects/remove/{project_id:[a-zA-Z0-9].*}", publicationProjectsController.ConfirmRemove).
		Methods("GET").
		Name("publication_projects_confirm_remove_from_publication")
	pubEditRouter.HandleFunc("/htmx/projects/remove/{project_id:[a-zA-Z0-9].*}", publicationProjectsController.Remove).
		Methods("PATCH").
		Name("publication_projects_remove_from_publication")
	// Publication departments HTMX fragments
	pubEditRouter.HandleFunc("/htmx/departments/list", publicationDepartmentsController.List).
		Methods("GET").
		Name("publication_add_department")
	pubEditRouter.HandleFunc("/htmx/departments/list/activesearch", publicationDepartmentsController.ActiveSearch).
		Methods("POST").
		Name("publicationDepartments_activesearch")
	pubEditRouter.HandleFunc("/htmx/departments/add/{department_id}", publicationDepartmentsController.Add).
		Methods("PATCH").
		Name("publicationDepartments_add_to_publication")
	pubEditRouter.HandleFunc("/htmx/departments/remove/{department_id}", publicationDepartmentsController.ConfirmRemove).
		Methods("GET").
		Name("publicationDepartments_confirm_remove_from_publication")
	pubEditRouter.HandleFunc("/htmx/departments/remove/{department_id}", publicationDepartmentsController.Remove).
		Methods("PATCH").
		Name("publicationDepartments_remove_from_publication")
	// Publication abstracts HTMX fragments
	pubEditRouter.HandleFunc("/htmx/abstracts/add", publicationAbstractsController.Add).
		Methods("GET").
		Name("publication_abstracts_add_abstract")
	pubEditRouter.HandleFunc("/htmx/abstracts/create", publicationAbstractsController.Create).
		Methods("POST").
		Name("publication_abstracts_create_abstract")
	pubEditRouter.HandleFunc("/htmx/abstracts/edit/{delta}", publicationAbstractsController.Edit).
		Methods("GET").
		Name("publication_abstracts_edit_abstract")
	pubEditRouter.HandleFunc("/htmx/abstracts/update/{delta}", publicationAbstractsController.Update).
		Methods("PUT").
		Name("publication_abstracts_update_abstract")
	pubEditRouter.HandleFunc("/htmx/abstracts/remove/{delta}", publicationAbstractsController.ConfirmRemove).
		Methods("GET").
		Name("publication_abstracts_confirm_remove_from_publication")
	pubEditRouter.HandleFunc("/htmx/abstracts/remove/{delta}", publicationAbstractsController.Remove).
		Methods("DELETE").
		Name("publication_abstracts_remove_abstract")

	// Publication lay summaries HTMX fragments
	pubEditRouter.HandleFunc("/htmx/lay_summaries/add", publicationLaySummariesController.Add).
		Methods("GET").
		Name("publication_lay_summaries_add_lay_summary")
	pubEditRouter.HandleFunc("/htmx/lay_summaries/create", publicationLaySummariesController.Create).
		Methods("POST").
		Name("publication_lay_summaries_create_lay_summary")
	pubEditRouter.HandleFunc("/htmx/lay_summaries/edit/{delta}", publicationLaySummariesController.Edit).
		Methods("GET").
		Name("publication_lay_summaries_edit_lay_summary")
	pubEditRouter.HandleFunc("/htmx/lay_summaries/update/{delta}", publicationLaySummariesController.Update).
		Methods("PUT").
		Name("publication_lay_summaries_update_lay_summary")
	pubEditRouter.HandleFunc("/htmx/lay_summaries/remove/{delta}", publicationLaySummariesController.ConfirmRemove).
		Methods("GET").
		Name("publication_lay_summaries_confirm_remove_from_publication")
	pubEditRouter.HandleFunc("/htmx/lay_summaries/remove/{delta}", publicationLaySummariesController.Remove).
		Methods("DELETE").
		Name("publication_lay_summaries_remove_lay_summary")

	// Publication links HTMX fragments
	pubEditRouter.HandleFunc("/htmx/links/add", publicationLinksController.Add).
		Methods("GET").
		Name("publication_links_add_link")
	pubEditRouter.HandleFunc("/htmx/links/create", publicationLinksController.Create).
		Methods("POST").
		Name("publication_links_create_link")
	pubEditRouter.HandleFunc("/htmx/links/edit/{delta}", publicationLinksController.Edit).
		Methods("GET").
		Name("publication_links_edit_link")
	pubEditRouter.HandleFunc("/htmx/links/update/{delta}", publicationLinksController.Update).
		Methods("PUT").
		Name("publication_links_update_link")
	pubEditRouter.HandleFunc("/htmx/links/remove/{delta}", publicationLinksController.ConfirmRemove).
		Methods("GET").
		Name("publication_links_confirm_remove_from_publication")
	pubEditRouter.HandleFunc("/htmx/links/remove/{delta}", publicationLinksController.Remove).
		Methods("DELETE").
		Name("publication_links_remove_link")
	// Publication contributors HTMX fragments
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/add", publicationContributorsController.Add).
		Methods("GET").
		Name("publication_contributors_add")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}", publicationContributorsController.Create).
		Methods("POST").
		Name("publication_contributors_create")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/order", publicationContributorsController.Order).
		Methods("POST").
		Name("publication_contributors_order")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/remove", publicationContributorsController.ConfirmRemove).
		Methods("GET").
		Name("publication_contributors_confirm_remove")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/{position}", publicationContributorsController.Remove).
		Methods("DELETE").
		Name("publication_contributors_remove")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/edit", publicationContributorsController.Edit).
		Methods("GET").
		Name("publication_contributors_edit")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/choose", publicationContributorsController.Choose).
		Methods("PUT").
		Name("publication_contributors_choose")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/demote", publicationContributorsController.Demote).
		Methods("PUT").
		Name("publication_contributors_demote")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/promote", publicationContributorsController.Promote).
		Methods("PUT").
		Name("publication_contributors_promote")
	pubEditRouter.HandleFunc("/htmx/contributors/{role}/{position}", publicationContributorsController.Update).
		Methods("PUT").
		Name("publication_contributors_update")
	// Publication datasets HTMX fragments
	pubEditRouter.HandleFunc("/htmx/datasets/choose", publicationDatasetsController.Choose).
		Methods("GET").
		Name("publication_datasets_choose")
	pubEditRouter.HandleFunc("/htmx/datasets/activesearch", publicationDatasetsController.ActiveSearch).
		Methods("POST").
		Name("publication_datasets_activesearch")
	pubEditRouter.HandleFunc("/htmx/datasets/add/{dataset_id}", publicationDatasetsController.Add).
		Methods("PATCH").
		Name("publication_datasets_add")
	pubEditRouter.HandleFunc("/htmx/datasets/remove/{dataset_id}", publicationDatasetsController.ConfirmRemove).
		Methods("GET").
		Name("publication_datasets_confirm_remove")
	pubEditRouter.HandleFunc("/htmx/datasets/remove/{dataset_id}", publicationDatasetsController.Remove).
		Methods("PATCH").
		Name("publication_datasets_remove")
}
