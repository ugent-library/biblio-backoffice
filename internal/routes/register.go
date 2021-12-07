package routes

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/controllers"
	"github.com/ugent-library/biblio-backend/internal/middleware"
	"github.com/ugent-library/go-locale/locale"
)

func Register(c controllers.Context) {
	router := c.Router
	basePath := c.BaseURL.Path

	router.Use(handlers.RecoveryHandler())

	// static files
	router.PathPrefix(basePath + "/static/").Handler(http.StripPrefix(basePath+"/static/", http.FileServer(http.Dir("./static"))))

	requireUser := middleware.RequireUser(c.BaseURL.Path + "/login")
	setUser := middleware.SetUser(c.Engine, c.SessionName, c.SessionStore)

	homeController := controllers.NewHome(c)

	authController := controllers.NewAuth(c)
	usersController := controllers.NewUsers(c)

	publicationsController := controllers.NewPublications(c)
	publicationFilesController := controllers.NewPublicationFiles(c)
	publicationDetailsController := controllers.NewPublicationDetails(c)
	publicationConferenceController := controllers.NewPublicationConference(c)
	publicationProjectsController := controllers.NewPublicationProjects(c)
	publicationDepartmentsController := controllers.NewPublicationDepartments(c)
	publicationAbstractsController := controllers.NewPublicationAbstracts(c)
	publicationLinksController := controllers.NewPublicationLinks(c)
	publicationContributorsController := controllers.NewPublicationContributors(c)
	publicationDatasetsController := controllers.NewPublicationDatasets(c)
	publicationAdditionalInfoController := controllers.NewPublicationAdditionalInfo(c)

	datasetsController := controllers.NewDatasets(c)
	datasetDetailsController := controllers.NewDatasetDetails(c)
	datasetProjectsController := controllers.NewDatasetProjects(c)
	datasetDepartmentsController := controllers.NewDatasetDepartments(c)
	datasetAbstractsController := controllers.NewDatasetAbstracts(c)
	datasetContributorsController := controllers.NewDatasetContributors(c)
	datasetPublicationsController := controllers.NewDatasetPublications(c)

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
		csrf.Secure(c.BaseURL.Scheme == "https"),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.FieldName("csrf-token"),
	)
	// TODO restrict to POST,PUT,PATCH
	r.Use(csrfMiddleware)

	// r.Use(handlers.HTTPMethodOverrideHandler)
	r.Use(locale.Detect(c.Localizer))

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
	pubsRouter.HandleFunc("/add-single", publicationsController.AddSingle).
		Methods("GET").
		Name("publication_add_single")
	pubsRouter.HandleFunc("/add-single/start", publicationsController.AddSingleStart).
		Methods("GET").
		Name("publication_add_single_start")
	pubsRouter.HandleFunc("/add-single/import", publicationsController.AddSingleImport).
		Methods("POST").
		Name("publication_add_single_import")
	pubsRouter.HandleFunc("/add-multiple", publicationsController.AddMultiple).
		Methods("GET").
		Name("publication_add_multiple")
	pubsRouter.HandleFunc("/add-multiple/start", publicationsController.AddMultipleStart).
		Methods("GET").
		Name("publication_add_multiple_start")
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
	pubRouter.Use(middleware.SetPublication(c.Engine))
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
	pubRouter.HandleFunc("/thumbnail", publicationsController.Thumbnail).
		Methods("GET").
		Name("publication_thumbnail")
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
	pubRouter.HandleFunc("/file/{file_id}/thumbnail", publicationFilesController.Thumbnail).
		Methods("GET").
		Name("publication_file_thumbnail")
	pubEditRouter.HandleFunc("/htmx/file", publicationFilesController.Upload).
		Methods("POST").
		Name("upload_publication_file")
	pubEditRouter.HandleFunc("/htmx/file/{file_id}/edit", publicationFilesController.Edit).
		Methods("GET").
		Name("publication_file_edit")
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
	pubEditRouter.HandleFunc("/htmx/projects/add/{project_id}", publicationProjectsController.Add).
		Methods("PATCH").
		Name("publication_projects_add_to_publication")
	pubEditRouter.HandleFunc("/htmx/projects/remove/{project_id}", publicationProjectsController.ConfirmRemove).
		Methods("GET").
		Name("publication_projects_confirm_remove_from_publication")
	pubEditRouter.HandleFunc("/htmx/projects/remove/{project_id}", publicationProjectsController.Remove).
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
	pubEditRouter.HandleFunc("/htmx/{role}/add", publicationContributorsController.Add).
		Methods("GET").
		Name("publication_contributors_add")
	pubEditRouter.HandleFunc("/htmx/{role}", publicationContributorsController.Create).
		Methods("POST").
		Name("publication_contributors_create")
	pubEditRouter.HandleFunc("/htmx/{role}/{position}/remove", publicationContributorsController.ConfirmRemove).
		Methods("GET").
		Name("publication_contributors_confirm_remove")
	pubEditRouter.HandleFunc("/htmx/{role}/{position}", publicationContributorsController.Remove).
		Methods("DELETE").
		Name("publication_contributors_remove")
	pubEditRouter.HandleFunc("/htmx/{role}/{position}/edit", publicationContributorsController.Edit).
		Methods("GET").
		Name("publication_contributors_edit")
	pubEditRouter.HandleFunc("/htmx/{role}/{position}/choose", publicationContributorsController.Choose).
		Methods("GET").
		Name("publication_contributors_choose")
	pubEditRouter.HandleFunc("/htmx/{role}/{position}", publicationContributorsController.Update).
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
	datasetsRouter.HandleFunc("/add/import", datasetsController.AddImport).
		Methods("POST").
		Name("dataset_add_import")

	datasetRouter := datasetsRouter.PathPrefix("/{id}").Subrouter()
	datasetRouter.Use(middleware.SetDataset(c.Engine))
	datasetRouter.Use(middleware.RequireCanViewDataset)
	datasetEditRouter := datasetRouter.PathPrefix("").Subrouter()
	datasetEditRouter.Use(middleware.RequireCanEditDataset)
	datasetPublishRouter := datasetRouter.PathPrefix("").Subrouter()
	datasetPublishRouter.Use(middleware.RequireCanPublishDataset)
	datasetDeleteRouter := datasetRouter.PathPrefix("").Subrouter()
	datasetDeleteRouter.Use(middleware.RequireCanDeleteDataset)
	datasetRouter.HandleFunc("", datasetsController.Show).
		Methods("GET").
		Name("dataset")
	datasetRouter.HandleFunc("/delete", datasetsController.ConfirmDelete).
		Methods("GET").
		Name("dataset_confirm_delete")
	// TODO why doesn't a DELETE with methodoverride work with CAS?
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
		Name("dataset_details_edit_form")
	datasetEditRouter.HandleFunc("/htmx/details/edit", datasetDetailsController.Update).
		Methods("PATCH").
		Name("dataset_details_save_form")
	// Dataset projects HTMX fragments
	datasetEditRouter.HandleFunc("/htmx/projects/list", datasetProjectsController.Choose).
		Methods("GET").
		Name("dataset_projects")
	datasetEditRouter.HandleFunc("/htmx/projects/list/activesearch", datasetProjectsController.ActiveSearch).
		Methods("POST").
		Name("dataset_projects_activesearch")
	datasetEditRouter.HandleFunc("/htmx/projects/add/{project_id}", datasetProjectsController.Add).
		Methods("PATCH").
		Name("dataset_projects_add")
	datasetEditRouter.HandleFunc("/htmx/projects/remove/{project_id}", datasetProjectsController.ConfirmRemove).
		Methods("GET").
		Name("dataset_projects_confirm_remove")
	datasetEditRouter.HandleFunc("/htmx/projects/remove/{project_id}", datasetProjectsController.Remove).
		Methods("PATCH").
		Name("dataset_projects_remove")
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
	// Publication contributors HTMX fragments
	datasetEditRouter.HandleFunc("/htmx/{role}/add", datasetContributorsController.Add).
		Methods("GET").
		Name("dataset_contributors_add")
	datasetEditRouter.HandleFunc("/htmx/{role}", datasetContributorsController.Create).
		Methods("POST").
		Name("dataset_contributors_create")
	datasetEditRouter.HandleFunc("/htmx/{role}/{position}/remove", datasetContributorsController.ConfirmRemove).
		Methods("GET").
		Name("dataset_contributors_confirm_remove")
	datasetEditRouter.HandleFunc("/htmx/{role}/{position}", datasetContributorsController.Remove).
		Methods("DELETE").
		Name("dataset_contributors_remove")
	datasetEditRouter.HandleFunc("/htmx/{role}/{position}/edit", datasetContributorsController.Edit).
		Methods("GET").
		Name("dataset_contributors_edit")
	datasetEditRouter.HandleFunc("/htmx/{role}/{position}/choose", datasetContributorsController.Choose).
		Methods("GET").
		Name("dataset_contributors_choose")
	datasetEditRouter.HandleFunc("/htmx/{role}/{position}", datasetContributorsController.Update).
		Methods("PUT").
		Name("dataset_contributors_update")
	// Dataset abstracts HTMX fragments
	datasetEditRouter.HandleFunc("/htmx/abstracts/add", datasetAbstractsController.Add).
		Methods("GET").
		Name("dataset_abstracts_add_abstract")
	datasetEditRouter.HandleFunc("/htmx/abstracts/create", datasetAbstractsController.Create).
		Methods("POST").
		Name("dataset_abstracts_create_abstract")
	datasetEditRouter.HandleFunc("/htmx/abstracts/edit/{delta}", datasetAbstractsController.Edit).
		Methods("GET").
		Name("dataset_abstracts_edit_abstract")
	datasetEditRouter.HandleFunc("/htmx/abstracts/update/{delta}", datasetAbstractsController.Update).
		Methods("PUT").
		Name("dataset_abstracts_update_abstract")
	datasetEditRouter.HandleFunc("/htmx/abstracts/remove/{delta}", datasetAbstractsController.ConfirmRemove).
		Methods("GET").
		Name("dataset_abstracts_confirm_remove_from_dataset")
	datasetEditRouter.HandleFunc("/htmx/abstracts/remove/{delta}", datasetAbstractsController.Remove).
		Methods("DELETE").
		Name("dataset_abstracts_remove_abstract")
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
}
