package publicationcreating

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/ulid"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindImport struct {
	Source          string `form:"source"`
	Identifier      string `form:"identifier"`
	PublicationType string `form:"publication_type"`
}

type YieldAddPublication struct {
	Context
	PageTitle            string
	Step                 int
	Source               string
	Identifier           string
	Publication          *models.Publication
	DuplicatePublication bool
	PublicationDatasets  []*models.Dataset
	ActiveNav            string
	SubNavs              []string // needed to render show_description
	ActiveSubNav         string   // needed to render show_description
	RedirectURL          string   // needed to render show_description
	DisplayDetails       *display.Display
	Errors               *YieldValidationErrors
}

type YieldValidationErrors struct {
	Title  string
	Errors form.Errors
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request, ctx Context) {
	tmpl := ""
	switch r.URL.Query().Get("method") {
	case "identifier":
		tmpl = "publication/add_identifier"
	case "manual":
		tmpl = "publication/add_manual"
	case "wos":
		tmpl = "publication/add_wos"
	case "bibtex":
		tmpl = "publication/add_bibtex"
	default:
		tmpl = "publication/add"
	}

	render.Wrap(w, "layouts/default", tmpl, YieldAddPublication{
		Context:   ctx,
		PageTitle: "Add - Publications - Biblio",
		Step:      1,
		ActiveNav: "publications",
	})
}

func (h *Handler) AddSingleImportConfirm(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindImport{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	// check for duplicates
	if b.Source == "crossref" && b.Identifier != "" {
		args := models.NewSearchArgs().
			WithFilter("doi", b.Identifier)

		existing, err := h.PublicationSearchService.Search(args)

		if err != nil {
			render.InternalServerError(w, r, err)
			return
		}

		if existing.Total > 0 {
			render.Wrap(w, "layouts/default", "dataset/add_identifier", YieldAddPublication{
				Context:              ctx,
				PageTitle:            "Add - Publications - Biblio",
				Step:                 1,
				ActiveNav:            "datasets",
				Source:               b.Source,
				Identifier:           b.Identifier,
				Publication:          existing.Hits[0],
				DuplicatePublication: true,
			})
			return
		}
	}

	h.AddSingleImport(w, r, ctx)
}

func (h *Handler) AddSingleImport(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindImport{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	var (
		p   *models.Publication
		err error
	)

	// import by identifier
	if b.Identifier != "" {
		p, err = h.fetchPublicationByIdentifier(b.Source, b.Identifier)
		if err != nil {
			log.Println(err)
			flash := flash.Flash{
				Type:         "error",
				Body:         template.HTML(ctx.T("publication.single_import.import_by_id.import_failed")),
				DismissAfter: 5 * time.Second,
			}

			ctx.Flash = append(ctx.Flash, flash)

			render.Wrap(w, "layouts/default", "publication/add_identifier", YieldAddPublication{
				Context:    ctx,
				PageTitle:  "Add - Publications - Biblio",
				Step:       1,
				ActiveNav:  "publications",
				Source:     b.Source,
				Identifier: b.Identifier,
			})
			return
		}
	} else {
		// or start with empty publication
		p = &models.Publication{Type: b.PublicationType}
	}

	p.ID = ulid.MustGenerate()
	p.CreatorID = ctx.User.ID
	p.UserID = ctx.User.ID
	p.Status = "private"
	p.Classification = "U"

	if validationErrs := p.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, err.(validation.Errors)))
		render.Wrap(w, "layouts/default", "publication/add_identifier", YieldAddPublication{
			Context:    ctx,
			PageTitle:  "Add - Publications - Biblio",
			Step:       1,
			ActiveNav:  "publications",
			Source:     b.Source,
			Identifier: b.Identifier,
			Errors: &YieldValidationErrors{
				Title:  "Unable to import this publication due to the following errors",
				Errors: errors,
			},
		})
		return
	}

	err = h.Repository.SavePublication(p)

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Wrap(w, "layouts/default", "publication/add_single_description", YieldAddPublication{
		Context:        ctx,
		PageTitle:      "Add - Publications - Biblio",
		Step:           2,
		ActiveNav:      "publications",
		SubNavs:        []string{"description", "files", "contributors", "datasets"},
		ActiveSubNav:   "description",
		Publication:    p,
		DisplayDetails: displays.PublicationDetails(ctx.Locale, p),
	})
}

func (h *Handler) AddSingleDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Wrap(w, "layouts/default", "publication/add_single_description", YieldAddPublication{
		Context:        ctx,
		PageTitle:      "Add - Publications - Biblio",
		Step:           2,
		ActiveNav:      "publications",
		SubNavs:        []string{"description", "files", "contributors", "datasets"},
		ActiveSubNav:   "description",
		Publication:    ctx.Publication,
		DisplayDetails: displays.PublicationDetails(ctx.Locale, ctx.Publication),
	})
}

func (h *Handler) AddSingleConfirm(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Wrap(w, "layouts/default", "publication/add_single_confirm", YieldAddPublication{
		Context:     ctx,
		PageTitle:   "Add - Publications - Biblio",
		Step:        3,
		ActiveNav:   "publications",
		Publication: ctx.Publication,
	})
}

func (h *Handler) AddSinglePublish(w http.ResponseWriter, r *http.Request, ctx Context) {
	if !ctx.User.CanPublishPublication(ctx.Publication) {
		render.Forbidden(w, r)
		return
	}

	ctx.Publication.Status = "public"

	if err := ctx.Publication.Validate(); err != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, err.(validation.Errors)))
		render.Render(w, "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to publish this publication due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Render(w, "error_dialog", ctx.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	redirectURL := h.PathFor("publication_add_finish", "id", ctx.Publication.ID)
	redirectURL.RawQuery = r.URL.Query().Encode()

	w.Header().Set("HX-Redirect", redirectURL.String())
}

func (h *Handler) AddSingleFinish(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Wrap(w, "layouts/default", "publication/add_single_finish", YieldAddPublication{
		Context:     ctx,
		PageTitle:   "Add - Publications - Biblio",
		Step:        4,
		ActiveNav:   "publications",
		Publication: ctx.Publication,
	})
}

func (h *Handler) fetchPublicationByIdentifier(source, identifier string) (*models.Publication, error) {
	s, ok := h.PublicationSources[source]

	if !ok {
		return nil, fmt.Errorf("unkown publication source: %s", source)
	}

	d, err := s.GetPublication(identifier)
	if err != nil {
		return nil, err
	}

	return d, nil
}
