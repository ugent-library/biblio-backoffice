package proxies

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	proxyviews "github.com/ugent-library/biblio-backoffice/views/proxy"
)

func Proxies(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	proxyviews.List(c).Render(r.Context(), w)
}
