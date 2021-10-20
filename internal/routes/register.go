package routes

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/controllers"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/middleware"
	"github.com/ugent-library/go-oidc/oidc"
	"github.com/unrolled/render"
)

func Register(baseURL *url.URL, e *engine.Engine, router *mux.Router, renderer *render.Render, sessionName string, sessionStore sessions.Store, oidcClient *oidc.Client) {
	// static files
	router.PathPrefix(baseURL.Path + "/static/").Handler(http.StripPrefix(baseURL.Path+"/static/", http.FileServer(http.Dir("./static"))))

	requireUser := middleware.RequireUser(baseURL.Path + "/logout")
	setUser := middleware.SetUser(e, sessionName, sessionStore)
	authController := controllers.NewAuth(e, sessionName, sessionStore, oidcClient, router)
	publicationController := controllers.NewPublications(e, renderer)
	datasetController := controllers.NewDatasets(e, renderer)
	publicationFilesController := controllers.NewPublicationsFiles(e, renderer, router)
	publicationDetailsController := controllers.NewPublicationsDetails(e, renderer)
	publicationConferenceController := controllers.NewPublicationConference(e, renderer)
	publicationProjectsController := controllers.NewPublicationProjects(e, renderer)
	publicationDepartmentsController := controllers.NewPublicationDepartments(e, renderer)
	publicationAuthorsController := controllers.NewPublicationAuthors(e, renderer)
	publicationDatasetsController := controllers.NewPublicationDatasets(e, renderer)
	datasetDetailsController := controllers.NewDatasetDetails(e, renderer)
	datasetProjectsController := controllers.NewDatasetProjects(e, renderer)

	// TODO fix absolute url generation
	// var schemes []string
	// if u.Scheme == "http" {
	// 	schemes = []string{"http", "https"}
	// } else {
	// 	schemes = []string{"https", "http"}
	// }
	// r = r.Schemes(schemes...).Host(u.Host).PathPrefix(u.Path).Subrouter()
	r := router.PathPrefix(baseURL.Path).Subrouter()

	// home
	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, baseURL.Path+"/publication", http.StatusFound)
	}).Methods("GET").Name("home")

	// auth
	r.HandleFunc("/login", authController.Login).
		Methods("GET").
		Name("login")
	r.HandleFunc("/auth/openid-connect/callback", authController.Callback).
		Methods("GET")
	r.HandleFunc("/logout", authController.Logout).
		Methods("GET").
		Name("logout")

	// publications
	publicationRouter := r.PathPrefix("/publication").Subrouter()
	publicationRouter.Use(middleware.SetActiveMenu("publications"))
	publicationRouter.Use(setUser)
	publicationRouter.Use(requireUser)
	publicationRouter.HandleFunc("", publicationController.List).
		Methods("GET").
		Name("publications")
	publicationRouter.HandleFunc("/new", publicationController.New).
		Methods("GET").
		Name("new_publication")
	publicationRouter.HandleFunc("/{id}", publicationController.Show).
		Methods("GET").
		Name("publication")
	publicationRouter.HandleFunc("/{id}/thumbnail", publicationController.Thumbnail).
		Methods("GET").
		Name("publication_thumbnail")

	// Publication files
	publicationRouter.HandleFunc("/{id}/file/{file_id}", publicationFilesController.Download).
		Methods("GET").
		Name("publication_file")
	publicationRouter.HandleFunc("/{id}/file/{file_id}/thumbnail", publicationFilesController.Thumbnail).
		Methods("GET").
		Name("publication_file_thumbnail")
	publicationRouter.HandleFunc("/{id}/file", publicationFilesController.Upload).
		Methods("POST").
		Name("upload_publication_file")

	// Publication HTMX fragments
	publicationRouter.HandleFunc("/{id}/htmx/summary", publicationController.Summary).
		Methods("GET").
		Name("publication_summary")

	// Publication details HTMX fragments
	publicationRouter.HandleFunc("/{id}/htmx", publicationDetailsController.Show).
		Methods("GET").
		Name("publication_details")
	publicationRouter.HandleFunc("/{id}/htmx/edit", publicationDetailsController.OpenForm).
		Methods("GET").
		Name("publication_details_edit_form")
	publicationRouter.HandleFunc("/{id}/htmx/edit", publicationDetailsController.SaveForm).
		Methods("PATCH").
		Name("publication_details_save_form")

	// Publication conference HTMX fragments
	publicationRouter.HandleFunc("/{id}/htmx/conference", publicationConferenceController.Show).
		Methods("GET").
		Name("publication_conference")
	publicationRouter.HandleFunc("/{id}/htmx/conference/edit", publicationConferenceController.OpenForm).
		Methods("GET").
		Name("publication_conference_edit_form")
	publicationRouter.HandleFunc("/{id}/htmx/conference/edit", publicationConferenceController.SaveForm).
		Methods("PATCH").
		Name("publication_conference_save_form")

	// Publication projects HTMX fragments
	publicationRouter.HandleFunc("/{id}/htmx/projects/list", publicationProjectsController.ListProjects).
		Methods("GET").
		Name("publication_projects")
	publicationRouter.HandleFunc("/{id}/htmx/projects/list/activesearch", publicationProjectsController.ActiveSearch).
		Methods("POST").
		Name("publication_projects_activesearch")
	publicationRouter.HandleFunc("/{id}/htmx/projects/add/{project_id}", publicationProjectsController.AddToPublication).
		Methods("PATCH").
		Name("publication_projects_add_to_publication")
	publicationRouter.HandleFunc("/{id}/htmx/projects/remove/{project_id}", publicationProjectsController.ConfirmRemoveFromPublication).
		Methods("GET").
		Name("publication_projects_confirm_remove_from_publication")
	publicationRouter.HandleFunc("/{id}/htmx/projects/remove/{project_id}", publicationProjectsController.RemoveFromPublication).
		Methods("PATCH").
		Name("publication_projects_remove_from_publication")

	// Publication departments HTMX fragments
	publicationRouter.HandleFunc("/{id}/htmx/departments/list", publicationDepartmentsController.ListDepartments).
		Methods("GET").
		Name("publicationDepartments")
	publicationRouter.HandleFunc("/{id}/htmx/departments/list/activesearch", publicationDepartmentsController.ActiveSearch).
		Methods("POST").
		Name("publicationDepartments_activesearch")
	publicationRouter.HandleFunc("/{id}/htmx/departments/add/{department_id}", publicationDepartmentsController.AddToPublication).
		Methods("PATCH").
		Name("publicationDepartments_add_to_publication")
	publicationRouter.HandleFunc("/{id}/htmx/departments/remove/{department_id}", publicationDepartmentsController.ConfirmRemoveFromPublication).
		Methods("GET").
		Name("publicationDepartments_confirm_remove_from_publication")
	publicationRouter.HandleFunc("/{id}/htmx/departments/remove/{department_id}", publicationDepartmentsController.RemoveFromPublication).
		Methods("PATCH").
		Name("publicationDepartments_remove_from_publication")

	// Publication authors HTMX fragments
	publicationRouter.HandleFunc("/{id}/htmx/authors/list", publicationAuthorsController.List).
		Methods("GET").
		Name("publication_authors_list")
	publicationRouter.HandleFunc("/{id}/htmx/authors/add/{delta}", publicationAuthorsController.AddRow).
		Methods("GET").
		Name("publication_authors_add_row")
	publicationRouter.HandleFunc("/{id}/htmx/authors/shift/{delta}", publicationAuthorsController.ShiftRow).
		Methods("GET").
		Name("publication_authors_shift_row")
	publicationRouter.HandleFunc("/{id}/htmx/authors/cancel/add/{delta}", publicationAuthorsController.CancelAddRow).
		Methods("DELETE").
		Name("publication_authors_cancel_add_row")
	publicationRouter.HandleFunc("/{id}/htmx/authors/create/{delta}", publicationAuthorsController.CreateAuthor).
		Methods("POST").
		Name("publication_authors_create_author")
	publicationRouter.HandleFunc("/{id}/htmx/authors/edit/{delta}", publicationAuthorsController.EditRow).
		Methods("GET").
		Name("publication_authors_edit_row")
	publicationRouter.HandleFunc("/{id}/htmx/authors/cancel/edit/{delta}", publicationAuthorsController.CancelEditRow).
		Methods("DELETE").
		Name("publication_authors_cancel_edit_row")
	publicationRouter.HandleFunc("/{id}/htmx/authors/update/{delta}", publicationAuthorsController.UpdateAuthor).
		Methods("POST").
		Name("publication_authors_update_author")
	publicationRouter.HandleFunc("/{id}/htmx/authors/remove/{delta}", publicationAuthorsController.ConfirmRemoveFromPublication).
		Methods("GET").
		Name("publication_authors_confirm_remove_from_publication")
	publicationRouter.HandleFunc("/{id}/htmx/authors/remove/{delta}", publicationAuthorsController.RemoveAuthor).
		Methods("DELETE").
		Name("publication_authors_remove_author")
	publicationRouter.HandleFunc("/{id}/htmx/authors/order/{start}/{end}", publicationAuthorsController.OrderAuthors).
		Methods("PUT").
		Name("publication_authors_order_authors")

	// Publication datasets HTMX fragments
	publicationRouter.HandleFunc("/{id}/htmx/datasets/choose", publicationDatasetsController.Choose).
		Methods("GET").
		Name("publication_datasets_choose")
	publicationRouter.HandleFunc("/{id}/htmx/datasets/activesearch", publicationDatasetsController.ActiveSearch).
		Methods("POST").
		Name("publication_datasets_activesearch")
	publicationRouter.HandleFunc("/{id}/htmx/datasets/add/{dataset_id}", publicationDatasetsController.Add).
		Methods("PATCH").
		Name("publication_datasets_add")
	publicationRouter.HandleFunc("/{id}/htmx/datasets/remove/{dataset_id}", publicationDatasetsController.ConfirmRemove).
		Methods("GET").
		Name("publication_datasets_confirm_remove")
	publicationRouter.HandleFunc("/{id}/htmx/datasets/remove/{dataset_id}", publicationDatasetsController.Remove).
		Methods("PATCH").
		Name("publication_datasets_remove")

	// datasets
	datasetRouter := r.PathPrefix("/dataset").Subrouter()
	datasetRouter.Use(middleware.SetActiveMenu("datasets"))
	datasetRouter.Use(setUser)
	datasetRouter.Use(requireUser)
	datasetRouter.HandleFunc("", datasetController.List).
		Methods("GET").
		Name("datasets")
	datasetRouter.HandleFunc("/{id}", datasetController.Show).
		Methods("GET").
		Name("dataset")

	// Dataset details HTMX fragments
	datasetRouter.HandleFunc("/{id}/htmx/details", datasetDetailsController.Show).
		Methods("GET").
		Name("dataset_details")
	datasetRouter.HandleFunc("/{id}/htmx/details/edit", datasetDetailsController.OpenForm).
		Methods("GET").
		Name("dataset_details_edit_form")
	datasetRouter.HandleFunc("/{id}/htmx/details/edit", datasetDetailsController.SaveForm).
		Methods("PATCH").
		Name("dataset_details_save_form")

	// Dataset projects HTMX fragmetns
	datasetRouter.HandleFunc("/{id}/htmx/projects/list", datasetProjectsController.ListProjects).
		Methods("GET").
		Name("dataset_projects")
	datasetRouter.HandleFunc("/{id}/htmx/projects/list/activesearch", datasetProjectsController.ActiveSearch).
		Methods("POST").
		Name("dataset_projects_activesearch")
	datasetRouter.HandleFunc("/{id}/htmx/projects/add/{project_id}", datasetProjectsController.AddToDataset).
		Methods("PATCH").
		Name("dataset_projects_add_to_dataset")
	datasetRouter.HandleFunc("/{id}/htmx/projects/remove/{project_id}", datasetProjectsController.ConfirmRemoveFromDataset).
		Methods("GET").
		Name("dataset_projects_confirm_remove_from_dataset")
	datasetRouter.HandleFunc("/{id}/htmx/projects/remove/{project_id}", datasetProjectsController.RemoveFromDataset).
		Methods("PATCH").
		Name("dataset_projects_remove_from_dataset")
}
