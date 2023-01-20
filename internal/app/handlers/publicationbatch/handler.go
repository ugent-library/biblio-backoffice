package publicationbatch

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
)

type Handler struct {
	handlers.BaseHandler
	Repository     backends.Repository
	ProjectService backends.ProjectService
}

type Context struct {
	handlers.BaseContext
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		if !ctx.User.CanCurate() {
			render.Forbidden(w, r)
			return
		}

		context := Context{
			BaseContext: ctx,
		}

		fn(w, r, context)
	})
}

var reSplitID = regexp.MustCompile(`[\s,;]+`)

type BindAddProjects struct {
	ProjectID      string `form:"project_id"`
	PublicationIDs string `form:"publication_ids"`
}

type YieldShow struct {
	Context
	PageTitle string
	ActiveNav string
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "layouts/default", "publication/batch/show", YieldShow{
		Context:   ctx,
		PageTitle: "Batch",
		ActiveNav: "publications",
	})
}

func (h *Handler) AddProjects(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAddProjects{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	project, err := h.ProjectService.GetProject(b.ProjectID)
	if err != nil {
		h.AddSessionFlash(r, w, *flash.SimpleFlash().
			WithLevel("error").
			WithBody(template.HTML(fmt.Sprintf("<p>could not find project %s</p>", b.ProjectID))).
			DismissedAfter(0))
		http.Redirect(w, r, h.PathFor("publication_batch").String(), http.StatusFound)
		return
	}

	done := 0
	var errorMsgs []string

	for _, id := range reSplitID.Split(b.PublicationIDs, -1) {
		pub, err := h.Repository.GetPublication(id)
		if err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>could not find publication %s</p>", id))
			continue
		}
		pub.AddProject(&models.PublicationProject{
			ID:   project.ID,
			Name: project.Title,
		})
		if err = h.Repository.UpdatePublication(pub.SnapshotID, pub, ctx.User); err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>could not update publication %s</p>", id))
			continue
		}
		done++
	}

	if done > 0 {
		h.AddSessionFlash(r, w, *flash.SimpleFlash().
			WithLevel("success").
			WithBody(template.HTML(fmt.Sprintf("<p>Successfully added project to %d publications.</p>", done))))
	}
	if len(errorMsgs) > 0 {
		h.AddSessionFlash(r, w, *flash.SimpleFlash().
			WithLevel("error").
			WithBody(template.HTML(strings.Join(errorMsgs, ""))).
			DismissedAfter(0))
	}
	http.Redirect(w, r, h.PathFor("publication_batch").String(), http.StatusFound)
}
