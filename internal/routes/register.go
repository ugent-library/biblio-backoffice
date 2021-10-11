package routes

import (
	"log"
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

func Register(baseURL string, e *engine.Engine, r *mux.Router, renderer *render.Render, sessionName string, sessionStore sessions.Store, oidcClient *oidc.Client) {
	requireUser := middleware.RequireUser("/logout")
	setUser := middleware.SetUser(e, sessionName, sessionStore)
	authController := controllers.NewAuth(e, sessionName, sessionStore, oidcClient)
	publicationController := controllers.NewPublications(e, renderer)
	datasetController := controllers.NewDatasets(e, renderer)
	publicationFilesController := controllers.NewPublicationsFiles(e, renderer, r)
	publicationDetailsController := controllers.NewPublicationsDetails(e, renderer)
	publicationProjectsController := controllers.NewPublicationProjects(e, renderer)
	publicationDepartmentsController := controllers.NewPublicationDepartments(e, renderer)
	publicationAuthorsController := controllers.NewPublicationAuthors(e, renderer)
	datasetDetailsController := controllers.NewDatasetDetails(e, renderer)
	datasetProjectsController := controllers.NewDatasetProjects(e, renderer)

	// build route urls from base url
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	// var schemes []string
	// if u.Scheme == "http" {
	// 	schemes = []string{"http", "https"}
	// } else {
	// 	schemes = []string{"https", "http"}
	// }

	// r = r.Schemes(schemes...).Host(u.Host).PathPrefix(u.Path).Subrouter()
	r = r.PathPrefix(u.Path).Subrouter()

	// static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// home
	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/publication", http.StatusFound)
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

	// Publication projects HTMX fragmetns
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
	// publicationRouter.HandleFunc("/{id}/htmx/authors/list", publicationAuthorsController.Listauthors).
	// 	Methods("GET").
	// 	Name("publicationAuthors")
	// publicationRouter.HandleFunc("/{id}/htmx/authors/list/activesearch", publicationAuthorsController.ActiveSearch).
	// 	Methods("POST").
	// 	Name("publicationAuthors_activesearch")
	publicationRouter.HandleFunc("/{id}/htmx/authors/add/cancel", publicationAuthorsController.CancelAddAuthorToTable).
		Methods("GET").
		Name("publication_authors_cancel_add_to_publication")
	publicationRouter.HandleFunc("/{id}/htmx/authors/add/{author_delta}", publicationAuthorsController.AddAuthorToTable).
		Methods("GET").
		Name("publication_authors_add_to_publication")
	publicationRouter.HandleFunc("/{id}/htmx/authors/save", publicationAuthorsController.SaveAuthorToPublication).
		Methods("PATCH").
		Name("publication_authors_save_to_publication")
	// publicationRouter.HandleFunc("/{id}/htmx/authors/remove/{author_id}", publicationAuthorsController.ConfirmRemoveFromPublication).
	// 	Methods("GET").
	// 	Name("publicationAuthors_confirm_remove_from_publication")
	// publicationRouter.HandleFunc("/{id}/htmx/authors/remove/{author_id}", publicationAuthorsController.RemoveFromPublication).
	// 	Methods("PATCH").
	// 	Name("publicationAuthors_remove_from_publication")

	// datasets
	datasetRouter := r.PathPrefix("/dataset").Subrouter()
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
