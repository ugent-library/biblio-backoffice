package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/controllers"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/middleware"
	"github.com/ugent-library/go-oidc/oidc"
	"github.com/unrolled/render"
)

func Register(e *engine.Engine, r *mux.Router, renderer *render.Render, sessionName string, sessionStore sessions.Store, oidcClient *oidc.Client) {
	requireUser := middleware.RequireUser("/logout")
	setUser := middleware.SetUser(e, sessionName, sessionStore)
	authController := controllers.NewAuth(e, sessionName, sessionStore, oidcClient)
	publicationController := controllers.NewPublications(e, renderer)
	datasetController := controllers.NewDatasets(e, renderer)
	publicationFilesController := controllers.NewPublicationsFiles(e, renderer)
	publicationDetailsController := controllers.NewPublicationsDetails(e, renderer)

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
}
