package proxies

import (
	"context"
	"net/http"

	"github.com/samber/lo"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/views"
	proxyviews "github.com/ugent-library/biblio-backoffice/views/proxy"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/htmx"
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

	if c.UserRole != "curator" {
		userProxies(w, r)
		return
	}

	proxies, err := findProxies(r.Context(), c, "", 0, 20)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	proxyviews.Index(c, proxies).Render(r.Context(), w)
}

func List(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := struct {
		ProxiesFilter string `query:"proxies_filter"`
		Offset        int    `query:"offset"`
	}{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	proxies, err := findProxies(r.Context(), c, b.ProxiesFilter, 20, b.Offset)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	proxyviews.List(c, proxies).Render(r.Context(), w)
}

func findProxies(rc context.Context, c *ctx.Ctx, q string, limit, offset int) ([][]*models.Person, error) {
	var personIDs []string
	if q != "" {
		hits, err := c.UserSearchService.SuggestUsers(q)
		if err != nil {
			return nil, err
		}
		// only exact matches
		if len(hits) != 1 {
			return nil, nil
		}
		personIDs = lo.Map(hits, func(p *models.Person, _ int) string {
			return p.ID
		})
	}

	var proxies [][]*models.Person
	pairs, err := c.Repo.FindProxies(rc, personIDs, limit, offset)
	if err != nil {
		return nil, err
	}
	for _, pair := range pairs {
		proxy := make([]*models.Person, 2)
		if p, err := c.PersonService.GetPerson(pair[0]); err == nil {
			proxy[0] = p
		} else {
			return nil, err
		}
		if p, err := c.PersonService.GetPerson(pair[1]); err == nil {
			proxy[1] = p
		} else {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

// TODO this makes way too many calls, all sequentially
func userProxies(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	ids, err := c.Repo.ProxyPersonIDs(r.Context(), c.User.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	proxies := make([]proxyviews.ProxiedPerson, 0, len(ids))
	for _, id := range ids {
		p, err := c.PersonService.GetPerson(id)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}

		withdrawnPublicationHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
			WithPageSize(0).
			WithFilter("creator_id|author_id", p.ID).
			WithFilter("status", "returned").
			WithFilter("locked", "false"))
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		draftPublicationHits, err := c.PublicationSearchIndex.Search(models.NewSearchArgs().
			WithPageSize(0).
			WithFilter("creator_id|author_id", p.ID).
			WithFilter("status", "private").
			WithFilter("locked", "false"))
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		withdrawnDatasetHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
			WithPageSize(0).
			WithFilter("creator_id|author_id", p.ID).
			WithFilter("status", "returned").
			WithFilter("locked", "false"))
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		draftDatasetHits, err := c.DatasetSearchIndex.Search(models.NewSearchArgs().
			WithPageSize(0).
			WithFilter("creator_id|author_id", p.ID).
			WithFilter("status", "private").
			WithFilter("locked", "false"))
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		proxies = append(proxies, proxyviews.ProxiedPerson{
			Person:                     p,
			WithdrawnPublicationsCount: withdrawnPublicationHits.Total,
			DraftPublicationsCount:     draftPublicationHits.Total,
			WithdrawnDatasetsCount:     withdrawnDatasetHits.Total,
			DraftDatasetsCount:         draftDatasetHits.Total,
		})
	}

	proxyviews.UserList(c, proxies).Render(r.Context(), w)
}

func AddProxy(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	hits, err := c.UserSearchService.SuggestUsers("")
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	hits = lo.Reject(hits, func(p *models.Person, _ int) bool {
		return p.ID == c.User.ID
	})

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

	hits = lo.Reject(hits, func(p *models.Person, _ int) bool {
		return p.ID == c.User.ID
	})

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

	hits = lo.Reject(hits, func(p *models.Person, _ int) bool {
		return p.ID == c.User.ID || p.ID == proxy.ID
	})

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

func People(w http.ResponseWriter, r *http.Request) {
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

	proxyviews.People(c, proxy, people).Render(r.Context(), w)
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

	hits = lo.Reject(hits, func(p *models.Person, _ int) bool {
		return p.ID == c.User.ID || p.ID == proxy.ID
	})

	peopleIDs, err := c.Repo.ProxyPersonIDs(r.Context(), b.ProxyID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	proxiedPeople := lo.Associate(peopleIDs, func(id string) (string, struct{}) { return id, struct{}{} })

	proxyviews.PeopleSuggestions(c, proxy, hits, proxiedPeople).Render(r.Context(), w)
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

	htmx.Trigger(w, "proxyChanged")
	w.WriteHeader(200)
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

	htmx.Trigger(w, "proxyChanged")
	w.WriteHeader(200)
}
