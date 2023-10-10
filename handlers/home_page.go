package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.HomePage(c).Render(r.Context(), w)
}

func ActionRequired(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.ActionRequired(c).Render(r.Context(), w)
}

func DraftsToComplete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.DraftsToComplete(c).Render(r.Context(), w)
}
