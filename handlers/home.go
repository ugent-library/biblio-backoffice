package handlers

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.HomePage(c).Render(r.Context(), w)
}
