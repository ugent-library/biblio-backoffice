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

	hits, err := c.UserSearchService.SuggestUsers("")
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ShowModal(proxyviews.Add(c, hits)).Render(r.Context(), w)
}

func SuggestProxies(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	q := r.URL.Query().Get("proxy_query")

	hits, err := c.UserSearchService.SuggestUsers(q)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	proxyviews.Suggestions(c, hits).Render(r.Context(), w)
}
