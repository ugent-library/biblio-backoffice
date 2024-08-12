package proxies

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
	proxyviews "github.com/ugent-library/biblio-backoffice/views/proxy"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

type bindProxy struct {
	ProxyID string `path:"proxy_id"`
}

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

func Edit(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := bindProxy{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	proxy, err := c.UserService.GetUser(b.ProxyID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	hits, err := c.PersonSearchService.SuggestPeople("")
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ReplaceModal(proxyviews.Edit(c, proxy, hits)).Render(r.Context(), w)
}

func SuggestPeople(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	q := r.URL.Query().Get("proxy_query")

	hits, err := c.PersonSearchService.SuggestPeople(q)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	proxyviews.PeopleSuggestions(c, hits).Render(r.Context(), w)
}
