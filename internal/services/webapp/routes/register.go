package routes

import (
	"net/http"

	"github.com/gorilla/csrf"
	mw "github.com/gorilla/handlers"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/datasetediting"
	"github.com/ugent-library/biblio-backend/internal/app/handlers/datasetviewing"
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
	router.Use(mw.RecoveryHandler())

	// static files
	router.PathPrefix(basePath + "/static/").Handler(http.StripPrefix(basePath+"/static/", http.FileServer(http.Dir("./internal/services/webapp/static"))))

	requireUser := middleware.RequireUser(oldBase.BaseURL.Path + "/login")
	setUser := middleware.SetUser(services.UserService, oldBase.SessionName, oldBase.SessionStore)

	homeController := controllers.NewHome(oldBase)
	authController := controllers.NewAuth(oldBase, oidcClient, services.UserService)
	usersController := controllers.NewUsers(oldBase, services.UserService)
	tasksController := controllers.NewTasks(oldBase, services.Tasks)

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

	datasetsController := controllers.NewDatasets(oldBase, services.Repository, services.DatasetSearchService, services.DatasetSources)
	datasetDetailsController := controllers.NewDatasetDetails(oldBase, services.Repository)
	datasetDepartmentsController := controllers.NewDatasetDepartments(oldBase, services.Repository, services.OrganizationSearchService, services.OrganizationService)
	datasetEditingHandlerontributorsController := controllers.NewDatasetContributors(oldBase, services.Repository, services.PersonSearchService, services.PersonService)
	datasetPublicationsController := controllers.NewDatasetPublications(oldBase, services.Repository, services.PublicationSearchService)

	licensesController := controllers.NewLicenses(oldBase, services.LicenseSearchService)
	mediaTypesController := controllers.NewMediaTypes(oldBase, services.MediaTypeSearchService)

	// NEW HANDLERS
	baseHandler := handlers.BaseHandler{
		SessionStore: oldBase.SessionStore,
		SessionName:  oldBase.SessionName,
		Localizer:    oldBase.Localizer,
		UserService:  services.UserService,
	}
	datasetViewingHandler := &datasetviewing.Handler{
		BaseHandler: baseHandler,
		Repository:  services.Repository,
	}
	datasetEditingHandler := &datasetediting.Handler{
		BaseHandler:          baseHandler,
		Repository:           services.Repository,
		ProjectService:       services.ProjectService,
		ProjectSearchService: services.ProjectSearchService,
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

	csrfMiddleware := csrf.Protect(
		[]byte(viper.GetString("csrf-secret")),
		csrf.CookieName(viper.GetString("csrf-name")),
		csrf.Path(basePath),
		csrf.Secure(oldBase.BaseURL.Scheme == "https"),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.FieldName("csrf-token"),
	)
	r.Use(csrfMiddleware)

	// NEW ROUTES
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
	// edit dataset abstracts
	r.HandleFunc("/dataset/{id}/abstracts/add", datasetEditingHandler.Wrap(datasetEditingHandler.AddAbstract)).
		Methods("GET").
		Name("dataset_add_abstract")
	r.HandleFunc("/dataset/{id}/abstracts", datasetEditingHandler.Wrap(datasetEditingHandler.CreateAbstract)).
		Methods("POST").
		Name("dataset_create_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{position}/edit", datasetEditingHandler.Wrap(datasetEditingHandler.EditAbstract)).
		Methods("GET").
		Name("dataset_edit_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{position}", datasetEditingHandler.Wrap(datasetEditingHandler.UpdateAbstract)).
		Methods("PUT").
		Name("dataset_update_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{position}/confirm-delete", datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteAbstract)).
		Methods("GET").
		Name("dataset_confirm_delete_abstract")
	r.HandleFunc("/dataset/{id}/abstracts/{position}", datasetEditingHandler.Wrap(datasetEditingHandler.DeleteAbstract)).
		Methods("DELETE").
		Name("dataset_delete_abstract")
	// edit dataset projects
	r.HandleFunc("/dataset/{id}/projects/add", datasetEditingHandler.Wrap(datasetEditingHandler.AddProject)).
		Methods("GET").
		Name("dataset_add_project")
	r.HandleFunc("/dataset/{id}/projects/suggestions", datasetEditingHandler.Wrap(datasetEditingHandler.ProjectSuggestions)).
		Methods("GET").
		Name("dataset_project_suggestions")
	r.HandleFunc("/dataset/{id}/projects", datasetEditingHandler.Wrap(datasetEditingHandler.CreateProject)).
		Methods("POST").
		Name("dataset_create_project")
	r.HandleFunc("/dataset/{id}/projects/{position}/confirm-delete", datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmDeleteProject)).
		Methods("GET").
		Name("dataset_confirm_delete_project")
	r.HandleFunc("/dataset/{id}/projects/{position}", datasetEditingHandler.Wrap(datasetEditingHandler.DeleteProject)).
		Methods("DELETE").
		Name("dataset_delete_project")

	// r.Use(handlers.HTTPMethodOverrideHandler)
	r.Use(locale.Detect(oldBase.Localizer))

	r.Use(setUser)

	// home
	r.HandleFunc("/", homeController.Home).Methods("GET").Name("home")

	// auth
	r.HandleFunc("/login", authController.Login).
		Methods("GET").
		Name("login")
	r.HandleFunc("/auth/openid-connect/callback", authController.Callback).
		Methods("GET")
	r.HandleFunc("/logout", authController.Logout).
		Methods("GET").
		Name("logout")

	// tasks
	taskRouter := r.PathPrefix("/task").Subrouter()
	taskRouter.Use(requireUser)
	taskRouter.HandleFunc("/{id}/status", tasksController.Status).
		Methods("GET").
		Name("task_status")

	// users
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(requireUser)
	userRouter.HandleFunc("/htmx/impersonate/choose", usersController.ImpersonateChoose).
		Methods("GET").
		Name("user_impersonate_choose")
	userRouter.HandleFunc("/impersonate", usersController.Impersonate).
		Methods("POST").
		Name("user_impersonate")
	// TODO why doesn't a DELETE with methodoverride work with CAS?
	userRouter.HandleFunc("/impersonate/remove", usersController.ImpersonateRemove).
		Methods("POST").
		Name("user_impersonate_remove")

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
	pubsRouter.HandleFunc("/orcid", publicationsController.ORCIDAddAll).
		Methods("POST").
		Name("publication_orcid_add_all")

	pubRouter := pubsRouter.PathPrefix("/{id}").Subrouter()
	pubRouter.Use(middleware.SetPublication(services.Repository))
	pubRouter.Use(middleware.RequireCanViewPublication)
	pubEditRouter := pubRouter.PathPrefix("").Subrouter()
	pubEditRouter.Use(middleware.RequireCanEditPublication)
	pubPublishRouter := pubRouter.PathPrefix("").Subrouter()
	pubPublishRouter.Use(middleware.RequireCanPublishPublication)
	pubDeleteRouter := pubRouter.PathPrefix("").Subrouter()
	pubDeleteRouter.Use(middleware.RequireCanDeletePublication)
	pubRouter.HandleFunc("", publicationsController.Show).
		Methods("GET").
		Name("publication")
	pubRouter.HandleFunc("/delete", publicationsController.ConfirmDelete).
		Methods("GET").
		Name("publication_confirm_delete")
	// TODO why doesn't a DELETE with methodoverride work with CAS?
	pubDeleteRouter.HandleFunc("/delete", publicationsController.Delete).
		Methods("POST").
		Name("publication_delete")
	pubRouter.HandleFunc("/orcid", publicationsController.ORCIDAdd).
		Methods("POST").
		Name("publication_orcid_add")
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
		Name("publication_projects")
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
		Name("publicationDepartments")
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

	// datasets
	datasetsRouter := r.PathPrefix("/dataset").Subrouter()
	datasetsRouter.Use(middleware.SetActiveMenu("datasets"))
	datasetsRouter.Use(requireUser)
	datasetsRouter.HandleFunc("", datasetsController.List).
		Methods("GET").
		Name("datasets")
	datasetsRouter.HandleFunc("/add", datasetsController.Add).
		Methods("GET").
		Name("dataset_add")
	datasetsRouter.HandleFunc("/add/import/confirm", datasetsController.AddImportConfirm).
		Methods("POST").
		Name("dataset_add_import_confirm")
	datasetsRouter.HandleFunc("/add/import", datasetsController.AddImport).
		Methods("POST").
		Name("dataset_add_import")

	datasetRouter := datasetsRouter.PathPrefix("/{id}").Subrouter()
	datasetRouter.Use(middleware.SetDataset(services.Repository))
	datasetRouter.Use(middleware.RequireCanViewDataset)
	datasetEditRouter := datasetRouter.PathPrefix("").Subrouter()
	datasetEditRouter.Use(middleware.RequireCanEditDataset)
	datasetPublishRouter := datasetRouter.PathPrefix("").Subrouter()
	datasetPublishRouter.Use(middleware.RequireCanPublishDataset)
	datasetDeleteRouter := datasetRouter.PathPrefix("").Subrouter()
	datasetDeleteRouter.Use(middleware.RequireCanDeleteDataset)
	datasetRouter.HandleFunc("/delete", datasetsController.ConfirmDelete).
		Methods("GET").
		Name("dataset_confirm_delete")
	datasetDeleteRouter.HandleFunc("/delete", datasetsController.Delete).
		Methods("POST").
		Name("dataset_delete")
	datasetEditRouter.HandleFunc("/publish", datasetsController.Publish).
		Methods("POST").
		Name("dataset_publish")
	datasetEditRouter.HandleFunc("/add/description", datasetsController.AddDescription).
		Methods("GET").
		Name("dataset_add_description")
	datasetEditRouter.HandleFunc("/add/confirm", datasetsController.AddConfirm).
		Methods("GET").
		Name("dataset_add_confirm")
	datasetEditRouter.HandleFunc("/add/publish", datasetsController.AddPublish).
		Methods("POST").
		Name("dataset_add_publish")
	// Dataset details HTMX fragments
	datasetEditRouter.HandleFunc("/htmx/details", datasetDetailsController.Show).
		Methods("GET").
		Name("dataset_details")
	datasetEditRouter.HandleFunc("/htmx/details/edit", datasetDetailsController.Edit).
		Methods("GET").
		Name("dataset_edit_details")
	datasetEditRouter.HandleFunc("/htmx/details/access_level", datasetDetailsController.AccessLevel).
		Methods("PUT").
		Name("dataset_edit_details_access_level")
	datasetEditRouter.HandleFunc("/htmx/details/edit", datasetDetailsController.Update).
		Methods("PATCH").
		Name("dataset_details_save_form")
	// Dataset departments HTMX fragments
	datasetEditRouter.HandleFunc("/htmx/departments/list", datasetDepartmentsController.List).
		Methods("GET").
		Name("datasetDepartments")
	datasetEditRouter.HandleFunc("/htmx/departments/list/activesearch", datasetDepartmentsController.ActiveSearch).
		Methods("POST").
		Name("datasetDepartments_activesearch")
	datasetEditRouter.HandleFunc("/htmx/departments/add/{department_id}", datasetDepartmentsController.Add).
		Methods("PATCH").
		Name("datasetDepartments_add_to_dataset")
	datasetEditRouter.HandleFunc("/htmx/departments/remove/{department_id}", datasetDepartmentsController.ConfirmRemove).
		Methods("GET").
		Name("datasetDepartments_confirm_remove_from_dataset")
	datasetEditRouter.HandleFunc("/htmx/departments/remove/{department_id}", datasetDepartmentsController.Remove).
		Methods("PATCH").
		Name("datasetDepartments_remove_from_dataset")
	// Dataset contributors HTMX fragments
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/add", datasetEditingHandlerontributorsController.Add).
		Methods("GET").
		Name("dataset_contributors_add")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}", datasetEditingHandlerontributorsController.Create).
		Methods("POST").
		Name("dataset_contributors_create")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/order", datasetEditingHandlerontributorsController.Order).
		Methods("POST").
		Name("dataset_contributors_order")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/remove", datasetEditingHandlerontributorsController.ConfirmRemove).
		Methods("GET").
		Name("dataset_contributors_confirm_remove")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/{position}", datasetEditingHandlerontributorsController.Remove).
		Methods("DELETE").
		Name("dataset_contributors_remove")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/edit", datasetEditingHandlerontributorsController.Edit).
		Methods("GET").
		Name("dataset_contributors_edit")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/choose", datasetEditingHandlerontributorsController.Choose).
		Methods("GET").
		Name("dataset_contributors_choose")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/demote", datasetEditingHandlerontributorsController.Demote).
		Methods("PUT").
		Name("dataset_contributors_demote")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/{position}/promote", datasetEditingHandlerontributorsController.Promote).
		Methods("PUT").
		Name("dataset_contributors_promote")
	datasetEditRouter.HandleFunc("/htmx/contributors/{role}/{position}", datasetEditingHandlerontributorsController.Update).
		Methods("PUT").
		Name("dataset_contributors_update")
	// Dataset publications HTMX fragments
	datasetEditRouter.HandleFunc("/htmx/publications/choose", datasetPublicationsController.Choose).
		Methods("GET").
		Name("dataset_publications_choose")
	datasetEditRouter.HandleFunc("/htmx/publications/activesearch", datasetPublicationsController.ActiveSearch).
		Methods("POST").
		Name("dataset_publications_activesearch")
	datasetEditRouter.HandleFunc("/htmx/publications/add/{publication_id}", datasetPublicationsController.Add).
		Methods("PATCH").
		Name("dataset_publications_add")
	datasetEditRouter.HandleFunc("/htmx/publications/remove/{publication_id}", datasetPublicationsController.ConfirmRemove).
		Methods("GET").
		Name("dataset_publications_confirm_remove")
	datasetEditRouter.HandleFunc("/htmx/publications/remove/{publication_id}", datasetPublicationsController.Remove).
		Methods("PATCH").
		Name("dataset_publications_remove")

	licensesRouter := r.PathPrefix("/license").Subrouter()
	licensesRouter.HandleFunc("/htmx/choose", licensesController.Choose).
		Methods("GET").
		Name("license_choose")

	mediaTypesRouter := r.PathPrefix("/media_types").Subrouter()
	mediaTypesRouter.HandleFunc("/htmx/choose", mediaTypesController.Choose).
		Methods("GET").
		Name("media_type_choose")
}
