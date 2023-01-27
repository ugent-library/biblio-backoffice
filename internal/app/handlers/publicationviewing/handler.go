package publicationviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/backends/filestore"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/render"
)

type Handler struct {
	handlers.BaseHandler
	Repository backends.Repository
	FileStore  *filestore.Store
}

type Context struct {
	handlers.BaseContext
	Publication *models.Publication
	RedirectURL string
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		p, err := h.Repository.GetPublication(bind.PathValues(r).Get("id"))
		if err != nil {
			if err == backends.ErrNotFound {
				h.NotFound(w, r, ctx)
			} else {
				render.InternalServerError(w, r, err)
			}
			return
		}

		if !ctx.User.CanViewPublication(p) {
			h.Logger.Warn("publication viewing: user isn't allowed to ivew the publication:", "errors", err, "publication", p.ID, "user", ctx.User.ID)
			render.Forbidden(w, r)
			return
		}

		redirectURL := r.URL.Query().Get("redirect-url")
		if redirectURL == "" {
			redirectURL = h.PathFor("publications").String()
		}

		fn(w, r, Context{
			BaseContext: ctx,
			Publication: p,
			RedirectURL: redirectURL,
		})
	})
}
