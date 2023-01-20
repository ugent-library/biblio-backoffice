// TODO this replicates to much of publicationviewing and publicationsearching
package publicationcreating

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backend/internal/app/displays"
	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/display"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
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
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("import confirm single publication: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	// check for duplicates
	if b.Source == "crossref" && b.Identifier != "" {
		args := models.NewSearchArgs().WithFilter("doi", strings.ToLower(b.Identifier)).WithFilter("status", "public")

		existing, err := h.PublicationSearchService.Search(args)

		if err != nil {
			h.Logger.Warnw("import single publication: could not execute search for duplicates", "errors", err, "args", args, "user", ctx.User.ID)
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
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		h.Logger.Warnw("import single publication: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
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
			h.Logger.Warnw("import single publication: could not fetch publication", "errors", err, "publication", b.Identifier, "user", ctx.User.ID)

			flash := flash.SimpleFlash().
				WithLevel("error").
				WithBody(template.HTML(ctx.Locale.T("publication.single_import.import_by_id.import_failed")))

			ctx.Flash = append(ctx.Flash, *flash)

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

	p.ID = ulid.Make().String()
	p.Creator = &models.PublicationUser{ID: ctx.User.ID, Name: ctx.User.FullName}
	p.User = &models.PublicationUser{ID: ctx.User.ID, Name: ctx.User.FullName}
	p.Status = "private"
	p.Classification = "U"

	// Set the first department of the user if the user resides under at least one department
	// TODO: this should be centralized
	if len(ctx.User.Department) > 0 {
		org, orgErr := h.OrganizationService.GetOrganization(ctx.User.Department[0].ID)
		if orgErr != nil {
			h.Logger.Warnw("import single publication: could not fetch user department", "errors", orgErr, "user", ctx.User.ID)
		} else {
			p.AddDepartmentByOrg(org)
		}
	}

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

	err = h.Repository.SavePublication(p, ctx.User)

	if err != nil {
		h.Logger.Errorf("import single publication: -could not save the publication:", "error", err, "identifier", b.Identifier, "user", ctx.User.ID)
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
		DisplayDetails: displays.PublicationDetails(ctx.User, ctx.Locale, p),
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
		DisplayDetails: displays.PublicationDetails(ctx.User, ctx.Locale, ctx.Publication),
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
	if !ctx.User.CanEditPublication(ctx.Publication) {
		h.Logger.Warnw("add single publication publish: user has no permission to publish publication.", "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.Forbidden(w, r)
		return
	}

	ctx.Publication.Status = "public"

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(ctx.Locale, validationErrs.(validation.Errors)))
		render.Layout(w, "show_modal", "form_errors_dialog", struct {
			Title  string
			Errors form.Errors
		}{
			Title:  "Unable to publish this publication due to the following errors",
			Errors: errors,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication, ctx.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "show_modal", "error_dialog", handlers.YieldErrorDialog{
			Message: ctx.Locale.T("publication.conflict_error_reload"),
		})
		return
	}

	if err != nil {
		h.Logger.Errorf("add single publication publish: could not save the publication:", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
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

	source := r.FormValue("source")

	file, _, err := r.FormFile("file")
	if err != nil {
		h.Logger.Warnw("add multiple import publication: could not retrieve file from request", "errors", err, "publication", ctx.Publication.ID, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}
	defer file.Close()

	// TODO why does the code imports zero entries without this?
	_, _ = file.Seek(0, io.SeekStart)

	batchID, err := h.importPublications(ctx.User, source, file)
	if err != nil {
		h.Logger.Warnw("add multiple import publication: could not import publications", "errors", err, "batch", batchID, "user", ctx.User.ID)

		flash := flash.SimpleFlash().
			WithLevel("error").
			WithBody(template.HTML("<p>Sorry, something went wrong. Could not import the publications.</p>"))

		ctx.Flash = append(ctx.Flash, *flash)

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

	// redirect to batch page so that pagination links work
	http.Redirect(
		w,
		r,
		h.PathFor("publication_add_multiple_confirm", "batch_id", batchID).String(),
		http.StatusFound)
}

func (h *Handler) AddMultipleSave(w http.ResponseWriter, r *http.Request, ctx Context) {
	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Publications successfully saved as draft.</p>"))

	h.AddSessionFlash(r, w, *flash)

	redirectURL := h.PathFor("publications")
	w.Header().Set("HX-Redirect", redirectURL.String())
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
		h.Logger.Warnw("add multiple confirm publication: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	searchArgs.WithFacets(vocabularies.Map["publication_facets"]...)

	batchID := bind.PathValues(r).Get("batch_id")

	hits, err := h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator.id", ctx.User.ID).
		WithScope("batch_id", batchID).
		Search(searchArgs)

	if err != nil {
		h.Logger.Errorw("add multiple confirm publication: could not execute search", "errors", err, "batch", batchID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/add_multiple_confirm", YieldAddMultiple{
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

func (h *Handler) AddMultiplePublish(w http.ResponseWriter, r *http.Request, ctx Context) {
	batchID := bind.PathValues(r).Get("batch_id")

	err := h.batchPublishPublications(batchID, ctx.User)

	// TODO this is useless to the user unless we point to the publication in
	// question
	var validationErrs validation.Errors
	if errors.As(err, &validationErrs) {
		h.Logger.Warnw("add multiple publish publication: could not validate abstract:", "errors", validationErrs, "batch", batchID, "user", ctx.User.ID)
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
		h.Logger.Errorw("add multiple publish publication: could not publish publications", "errors", err, "batch", batchID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	redirectURL := h.PathFor("publication_add_multiple_finish", "batch_id", batchID)

	w.Header().Set("HX-Redirect", redirectURL.String())
}

func (h *Handler) AddMultipleFinish(w http.ResponseWriter, r *http.Request, ctx Context) {
	searchArgs := models.NewSearchArgs()
	if err := bind.RequestQuery(r, searchArgs); err != nil {
		h.Logger.Warnw("add multiple finish publication: could not bind request arguments", "errors", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	searchArgs.WithFacets(vocabularies.Map["publication_facets"]...)

	batchID := bind.PathValues(r).Get("batch_id")

	hits, err := h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator.id", ctx.User.ID).
		WithScope("batch_id", batchID).
		Search(searchArgs)

	if err != nil {
		h.Logger.Errorw("add multiple finish publication: could not execute search", "errors", err, "batch", batchID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.Layout(w, "layouts/default", "publication/pages/add_multiple_finish", YieldAddMultiple{
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

func (h *Handler) importPublications(user *models.User, source string, file io.Reader) (string, error) {
	batchID := ulid.Make().String()

	decFactory, ok := h.PublicationDecoders[source]
	if !ok {
		return "", errors.New("unknown publication source")
	}
	dec := decFactory(file)

	var importErr error
	for {
		p := models.Publication{
			ID:             ulid.Make().String(),
			BatchID:        batchID,
			Status:         "private",
			Classification: "U",
			Creator:        &models.PublicationUser{ID: user.ID, Name: user.FullName},
			User:           &models.PublicationUser{ID: user.ID, Name: user.FullName},
		}

		// Set the department if the user was assigned to at least one department
		// TODO: this should be centralized
		if len(user.Department) > 0 {
			org, orgErr := h.OrganizationService.GetOrganization(user.Department[0].ID)
			if orgErr != nil {
				h.Logger.Warnw("add multiple publications: could not fetch user department", "errors", orgErr, "user", user.ID)
			} else {
				p.AddDepartmentByOrg(org)
			}
		}

		if err := dec.Decode(&p); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			importErr = err
			break
		}

		if err := h.Repository.SavePublication(&p, user); err != nil {
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
func (h *Handler) batchPublishPublications(batchID string, user *models.User) (err error) {
	searcher := h.PublicationSearchService.
		WithScope("status", "private", "public").
		WithScope("creator.id", user.ID).
		WithScope("batch_id", batchID)
	args := models.NewSearchArgs()

	var hits *models.PublicationHits
	for {
		hits, err = searcher.Search(args)
		if err != nil {
			return
		}
		for _, pub := range hits.Hits {
			// TODO check CanEditPublication
			pub.Status = "public"
			if err = pub.Validate(); err != nil {
				return
			}
			if err = h.Repository.SavePublication(pub, user); err != nil {
				return
			}
		}
		if !hits.HasNextPage() {
			break
		}
		args.Page = args.Page + 1
	}
	return
}
