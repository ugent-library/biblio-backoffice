package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.User != nil {
		http.Redirect(w, r, c.PathTo("dashboard").String(), http.StatusSeeOther)
	} else {
		views.Home(c).Render(r.Context(), w)
	}
}
