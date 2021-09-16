package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/controllers"
)

func Register(r *mux.Router,
	requireUser mux.MiddlewareFunc,
	setUser mux.MiddlewareFunc,
	authController *controllers.Auth,
	publicationController *controllers.Publications,
	publicationDetailsController *controllers.PublicationsDetails) {

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

	// Publication details
	publicationRouter.HandleFunc("/{id}/details", publicationDetailsController.Show).
		Methods("GET").
		Name("publication_details")
	publicationRouter.HandleFunc("/{id}/details/edit", publicationDetailsController.OpenForm).
		Methods("GET").
		Name("publication_details_edit_form")
	publicationRouter.HandleFunc("/{id}/details/edit", publicationDetailsController.SaveForm).
		Methods("PATCH").
		Name("publication_details_save_form")
}
