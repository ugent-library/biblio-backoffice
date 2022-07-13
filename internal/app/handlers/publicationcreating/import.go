// TODO this replicates to much of publicationviewing and publicationsearching
package publicationcreating

import (
	"errors"
	"fmt"
	"html/template"
	"io"
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

type BindImportSingle struct {
	Source          string `form:"source"`
	Identifier      string `form:"identifier"`
	PublicationType string `form:"publication_type"`
}

type YieldAddSingle struct {
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

type YieldAddMultiple struct {
	Context
	PageTitle   string
	Step        int
	ActiveNav   string
	RedirectURL string
	BatchID     string
	SearchArgs  *models.SearchArgs
	Hits        *models.PublicationHits
}

type YieldHit struct {
	Context
	Publication *models.Publication
}

func (y YieldAddMultiple) YieldHit(d *models.Publication) YieldHit {
	return YieldHit{y.Context, d}
}

type YieldAddMultipleShow struct {
	Context
	PageTitle    string
	Step         int
	ActiveNav    string
	SubNavs      []string // needed to render show_description
	ActiveSubNav string   // needed to render show_description
	RedirectURL  string   // needed to render show_description
	BatchID      string
}

type YieldValidationErrors struct {
	Title  string
	Errors form.Errors
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request, ctx Context) {
	tmpl := ""
	switch r.URL.Query().Get("method") {
	case "identifier":
		tmpl = "publication/pages/add_identifier"
	case "manual":
		tmpl = "publication/pages/add_manual"
	case "wos":
		tmpl = "publication/pages/add_wos"
	case "bibtex":
		tmpl = "publication/pages/add_bibtex"
	default:
		tmpl = "publication/pages/add"
	}

	render.Layout(w, "layouts/default", tmpl, YieldAddSingle{
		Context:   ctx,
		PageTitle: "Add - Publications - Biblio",
		Step:      1,
		ActiveNav: "publications",
	})
}

func (h *Handler) AddSingleImportConfirm(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindImportSingle{}
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
			render.Layout(w, "layouts/default", "publication/pages/add_identifier", YieldAddSingle{
				Context:              ctx,
				PageTitle:            "Add - Publications - Biblio",
				Step:                 1,
				ActiveNav:            "publications",
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
	b := BindImportSingle{}
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
			ctx.Flash = append(ctx.Flash, flash.Flash{
				Type: "error",
				Body: template.HTML(ctx.Locale.T("publication.single_import.import_by_id.import_failed")),
			})

			render.Layout(w, "layouts/default", "publication/pages/add_identifier", YieldAddSingle{
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
		render.Layout(w, "layouts/default", "publication/pages/add_identifier", YieldAddSingle{
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

	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}

	render.Layout(w, "layouts/default", "publication/pages/add_single_description", YieldAddSingle{
		Context:        ctx,
		PageTitle:      "Add - Publications - Biblio",
		Step:           2,
		ActiveNav:      "publications",
		SubNavs:        []string{"description", "files", "contributors", "datasets"},
		ActiveSubNav:   subNav,
		Publication:    p,
		DisplayDetails: displays.PublicationDetails(ctx.Locale, p),
	})
}

func (h *Handler) AddSingleDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}

	render.Layout(w, "layouts/default", "publication/pages/add_single_description", YieldAddSingle{
		Context:        ctx,
		PageTitle:      "Add - Publications - Biblio",
		Step:           2,
		ActiveNav:      "publications",
		SubNavs:        []string{"description", "files", "contributors", "datasets"},
		ActiveSubNav:   subNav,
		Publication:    ctx.Publication,
		DisplayDetails: displays.PublicationDetails(ctx.Locale, ctx.Publication),
	})
}

func (h *Handler) AddSingleConfirm(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "layouts/default", "publication/pages/add_single_confirm", YieldAddSingle{
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
		render.Layout(w, "show_modal", "form_errors_dialog", struct {
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
		render.Layout(w, "show_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	redirectURL := h.PathFor("publication_add_single_finish", "id", ctx.Publication.ID)
	redirectURL.RawQuery = r.URL.Query().Encode()

	w.Header().Set("HX-Redirect", redirectURL.String())
}

func (h *Handler) AddSingleFinish(w http.ResponseWriter, r *http.Request, ctx Context) {
	render.Layout(w, "layouts/default", "publication/pages/add_single_finish", YieldAddSingle{
		Context:     ctx,
		PageTitle:   "Add - Publications - Biblio",
		Step:        4,
		ActiveNav:   "publications",
		Publication: ctx.Publication,
	})
}

func (h *Handler) AddMultipleImport(w http.ResponseWriter, r *http.Request, ctx Context) {
	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2000000000)

	// buffer limit of 32MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	source := r.FormValue("source")

	file, _, err := r.FormFile("file")
	if err != nil {
		render.BadRequest(w, r, err)
		return
	}
	defer file.Close()

	//TODO: why does the code imports zero entries without this?
	_, _ = file.Seek(0, io.SeekStart)

	batchID, err := h.importPublications(ctx.User.ID, source, file)
	if err != nil {
		ctx.Flash = append(ctx.Flash, flash.Flash{
			Type: "error",
			Body: "Sorry, something went wrong. Could not import the publications.",
		})

		tmpl := ""
		switch source {
		case "wos":
			tmpl = "publication/pages/add_wos"
		case "bibtex":
			tmpl = "publication/pages/add_bibtex"
		}

		render.Layout(w, "layouts/default", tmpl, YieldAddSingle{
			Context:   ctx,
			PageTitle: "Add - Publications - Biblio",
			Step:      1,
			ActiveNav: "publications",
			Source:    source,
		})
		return
	}

	// TODO wait for index refresh, do something more elegant
	time.Sleep(time.Second)

	searchArgs := models.NewSearchArgs()

	hits, err := h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator_id", ctx.User.ID).
		WithScope("batch_id", batchID).
		IncludeFacets(true).
		Search(searchArgs)

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/add_multiple_description", YieldAddMultiple{
		Context:     ctx,
		PageTitle:   "Add - Publications - Biblio",
		Step:        2,
		ActiveNav:   "publications",
		RedirectURL: h.PathFor("publication_add_multiple_description", "batch_id", batchID).String(),
		BatchID:     batchID,
		SearchArgs:  searchArgs,
		Hits:        hits,
	})
}

func (h *Handler) AddMultipleDescription(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.RequestQuery(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	batchID := bind.PathValues(r).Get("batch_id")

	hits, err := h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator_id", ctx.User.ID).
		WithScope("batch_id", batchID).
		IncludeFacets(true).
		Search(searchArgs)

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/add_multiple_description", YieldAddMultiple{
		Context:     ctx,
		PageTitle:   "Add - Publications - Biblio",
		Step:        2,
		ActiveNav:   "publications",
		RedirectURL: r.URL.String(),
		BatchID:     batchID,
		SearchArgs:  searchArgs,
		Hits:        hits,
	})
}

// TODO after changing tabs, the wrong url is pushed in the history
func (h *Handler) AddMultipleShow(w http.ResponseWriter, r *http.Request, ctx Context) {
	batchID := bind.PathValues(r).Get("batch_id")
	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}

	render.Layout(w, "layouts/default", "publication/pages/add_multiple_show", YieldAddMultipleShow{
		Context:      ctx,
		PageTitle:    "Add - Publications - Biblio",
		Step:         2,
		ActiveNav:    "publications",
		SubNavs:      []string{"description", "files", "contributors", "datasets"},
		ActiveSubNav: subNav,
		RedirectURL:  r.URL.Query().Get("redirect-url"),
		BatchID:      batchID,
	})
}

func (h *Handler) AddMultipleConfirm(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.RequestQuery(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	batchID := bind.PathValues(r).Get("batch_id")

	hits, err := h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator_id", ctx.User.ID).
		WithScope("batch_id", batchID).
		IncludeFacets(true).
		Search(searchArgs)

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/add_multiple_confirm", YieldAddMultiple{
		Context:     ctx,
		PageTitle:   "Add - Publications - Biblio",
		Step:        3,
		ActiveNav:   "publications",
		RedirectURL: r.URL.String(),
		BatchID:     batchID,
		SearchArgs:  searchArgs,
		Hits:        hits,
	})
}

func (h *Handler) AddMultiplePublish(w http.ResponseWriter, r *http.Request, ctx Context) {
	batchID := bind.PathValues(r).Get("batch_id")

	err := h.batchPublishPublications(batchID, ctx.User.ID)

	// TODO this is useless to the user unless we point to the publication in
	// question
	var validationErrs validation.Errors
	if errors.As(err, &validationErrs) {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs))
		render.Layout(w, "show_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to publish a publication due to the following errors",
			Errors: errors,
		})
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	redirectURL := h.PathFor("publication_add_multiple_finish", "batch_id", batchID)

	w.Header().Set("HX-Redirect", redirectURL.String())
}

func (h *Handler) AddMultipleFinish(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.RequestQuery(r, searchArgs); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	batchID := bind.PathValues(r).Get("batch_id")

	hits, err := h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator_id", ctx.User.ID).
		WithScope("batch_id", batchID).
		IncludeFacets(true).
		Search(searchArgs)

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/add_multiple_finish", YieldAddMultiple{
		Context:     ctx,
		PageTitle:   "Add - Publications - Biblio",
		Step:        4,
		ActiveNav:   "publications",
		RedirectURL: r.URL.String(),
		BatchID:     batchID,
		SearchArgs:  searchArgs,
		Hits:        hits,
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

func (h *Handler) importPublications(userID, source string, file io.Reader) (string, error) {
	batchID := ulid.MustGenerate()

	decFactory, ok := h.PublicationDecoders[source]
	if !ok {
		return "", errors.New("unknown publication source")
	}
	dec := decFactory(file)

	var importErr error
	for {
		p := models.Publication{
			ID:             ulid.MustGenerate(),
			BatchID:        batchID,
			Status:         "private",
			Classification: "U",
			CreatorID:      userID,
			UserID:         userID,
		}
		if err := dec.Decode(&p); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			importErr = err
			break
		}

		if err := h.Repository.SavePublication(&p); err != nil {
			importErr = err
			break
		}
	}

	// TODO rollback if error
	if importErr != nil {
		return "", importErr
	}

	return batchID, nil
}

// TODO check conflicts?
func (h *Handler) batchPublishPublications(batchID, userID string) (err error) {
	searcher := h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator_id", userID).
		WithScope("batch_id", batchID)
	args := models.NewSearchArgs()

	var hits *models.PublicationHits
	for {
		hits, err = searcher.Search(args)
		if err != nil {
			return
		}
		for _, pub := range hits.Hits {
			// TODO check CanPublishPublication
			pub.Status = "public"
			if err = pub.Validate(); err != nil {
				return
			}
			if err = h.Repository.SavePublication(pub); err != nil {
				return
			}
		}
		if !hits.NextPage() {
			break
		}
		args.Page = args.Page + 1
	}
	return
}
