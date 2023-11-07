package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
)

func DashBoard(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.UserRole == "curator" {
		// TODO port and render here as CuratorDashboard
		http.Redirect(w, r, c.PathTo("dashboard_publications", "type", "faculties").String(), http.StatusSeeOther)
	} else {
		views.UserDashboard(c).Render(r.Context(), w)
	}
}
