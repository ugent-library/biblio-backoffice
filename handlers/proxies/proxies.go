package proxies

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
	proxyviews "github.com/ugent-library/biblio-backoffice/views/proxy"
)

func Proxies(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	proxyviews.List(c).Render(r.Context(), w)
}

func AddProxy(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	views.ShowModal(proxyviews.Add(c)).Render(r.Context(), w)
}
