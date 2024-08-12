package proxies

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
	proxyviews "github.com/ugent-library/biblio-backoffice/views/proxy"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

type bindProxy struct {
	ProxyID string `path:"proxy_id"`
}

type bindAddProxyPerson struct {
	ProxyID  string `path:"proxy_id"`
	PersonID string `form:"person_id"`
}

type bindRemoveProxyPerson struct {
	ProxyID  string `path:"proxy_id"`
	PersonID string `path:"person_id"`
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

	peopleIDs, err := c.Repo.ProxyPersonIDs(r.Context(), b.ProxyID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	people := make([]*models.Person, len(peopleIDs))
	for i, id := range peopleIDs {
		person, err := c.PersonService.GetPerson(id)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		people[i] = person
	}

	views.ReplaceModal(proxyviews.Edit(c, proxy, people, hits)).Render(r.Context(), w)
}

func SuggestPeople(w http.ResponseWriter, r *http.Request) {
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

	q := r.URL.Query().Get("proxy_query")

	hits, err := c.PersonSearchService.SuggestPeople(q)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	proxyviews.PeopleSuggestions(c, proxy, hits).Render(r.Context(), w)
}

func AddPerson(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := bindAddProxyPerson{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	if err := c.Repo.AddProxyPerson(r.Context(), b.ProxyID, b.PersonID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	proxy, err := c.UserService.GetUser(b.ProxyID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	peopleIDs, err := c.Repo.ProxyPersonIDs(r.Context(), b.ProxyID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	people := make([]*models.Person, len(peopleIDs))
	for i, id := range peopleIDs {
		person, err := c.PersonService.GetPerson(id)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		people[i] = person
	}

	views.ReplaceMany(map[string]templ.Component{
		"outerHTML:#people": proxyviews.People(c, proxy, people),
	}).Render(r.Context(), w)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := bindRemoveProxyPerson{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	if err := c.Repo.RemoveProxyPerson(r.Context(), b.ProxyID, b.PersonID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	proxy, err := c.UserService.GetUser(b.ProxyID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	peopleIDs, err := c.Repo.ProxyPersonIDs(r.Context(), b.ProxyID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	people := make([]*models.Person, len(peopleIDs))
	for i, id := range peopleIDs {
		person, err := c.PersonService.GetPerson(id)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		people[i] = person
	}

	views.ReplaceMany(map[string]templ.Component{
		"outerHTML:#people": proxyviews.People(c, proxy, people),
	}).Render(r.Context(), w)
}
