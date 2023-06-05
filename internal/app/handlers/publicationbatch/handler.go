package publicationbatch

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backoffice/internal/app/handlers"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/render"
	"github.com/ugent-library/biblio-backoffice/internal/render/flash"
)

type Handler struct {
	handlers.BaseHandler
	Repository backends.Repository
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

func (h *Handler) Process(w http.ResponseWriter, r *http.Request, ctx Context) {
	lines := strings.Split(strings.ReplaceAll(r.FormValue("ops"), "\r\n", "\n"), "\n")

	if len(lines) > 500 {
		h.AddFlash(r, w, *flash.SimpleFlash().
			WithLevel("error").
			WithBody("No more than 500 operations can be processed at one time.").
			DismissedAfter(0))
		http.Redirect(w, r, h.PathFor("publication_batch").String(), http.StatusFound)
		return
	}

	done := 0
	var errorMsgs []string

	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		rdr := csv.NewReader(strings.NewReader(strings.TrimSpace(line)))
		rec, err := rdr.Read()

		if err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>error parsing line %d</p>", i))
			continue
		}

		if len(rec) < 2 {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>error parsing line %d</p>", i))
			continue
		}

		id := strings.TrimSpace(rec[0])
		op := strings.TrimSpace(rec[1])
		args := rec[2:]
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)
		}

		err = h.Repository.MutatePublication(id, ctx.User, backends.Mutation{
			Op:   op,
			Args: args,
		})

		if err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>could not process publication %s at line %d</p>", id, i))
			continue
		}

		done++
	}

	if done > 0 {
		h.AddFlash(r, w, *flash.SimpleFlash().
			WithLevel("success").
			WithBody(template.HTML(fmt.Sprintf("<p>Successfully processed %d publications.</p>", done))))
	}
	if len(errorMsgs) > 0 {
		h.AddFlash(r, w, *flash.SimpleFlash().
			WithLevel("error").
			WithBody(template.HTML(strings.Join(errorMsgs, ""))).
			DismissedAfter(0))
	}

	http.Redirect(w, r, h.PathFor("publication_batch").String(), http.StatusFound)
}
