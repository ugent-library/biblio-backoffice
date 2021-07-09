package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/controllers"
)

func Register(r *mux.Router, publicationController *controllers.Publication) {
	// static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// home
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/publications", http.StatusFound)
	}).Methods("GET").Name("home")

	// publications
	r.HandleFunc("/publications", publicationController.List).
		Methods("GET").
		Name("publications")
}
