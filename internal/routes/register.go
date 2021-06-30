package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/ugent-library/biblio-backend/internal/controllers"
)

func Register(router chi.Router, publicationController *controllers.Publication) {
	// general middleware
	router.Use(chimw.RequestID)
	router.Use(chimw.RealIP)
	router.Use(chimw.Logger)
	router.Use(chimw.Recoverer)

	// static files
	router.Mount("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// home
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/publication", http.StatusFound)
	})

	// publication
	router.Get("/publication", publicationController.List)
}
