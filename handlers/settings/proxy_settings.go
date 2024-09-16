package settings

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	settingsviews "github.com/ugent-library/biblio-backoffice/views/settings"
)

func ProxySettings(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	proxyIDs, err := c.Repo.ProxyIDs(r.Context(), c.User.IDs)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	proxies := make([]*models.Person, len(proxyIDs))
	for i, id := range proxyIDs {
		person, err := c.PersonService.GetPerson(id)
		if err != nil {
			c.HandleError(w, r, err)
			return
		}
		proxies[i] = person
	}

	settingsviews.ProxySettings(c, proxies).Render(r.Context(), w)
}
