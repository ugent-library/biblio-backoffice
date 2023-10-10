package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
)

func ActionRequired(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.ActionRequired(c).Render(r.Context(), w)
}
