package datasetviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/bind"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
)

type Handler struct {
	handlers.BaseHandler
	Repo *repositories.Repo
}

type Context struct {
	handlers.BaseContext
	Dataset     *models.Dataset
	RedirectURL string
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		d, err := h.Repo.GetDataset(bind.PathValues(r).Get("id"))
		if err != nil {
			if err == models.ErrNotFound {
				h.NotFound(w, r, ctx)
			} else {
				render.InternalServerError(w, r, err)
			}
			return
		}

		if !ctx.User.CanViewDataset(d) {
			h.Logger.Warn("view dataset: user isn't allowed to view the dataset:", "error", err, "dataset", d.ID, "user", ctx.User.ID)
			render.Forbidden(w, r)
			return
		}

		redirectURL := r.URL.Query().Get("redirect-url")
		if redirectURL == "" {
			redirectURL = h.PathFor("datasets").String()
		}

		fn(w, r, Context{
			BaseContext: ctx,
			Dataset:     d,
			RedirectURL: redirectURL,
		})
	})
}
