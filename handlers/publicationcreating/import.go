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
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/render/form"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/biblio-backoffice/views/publication/pages"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindImportSingle struct {
	Source          string `form:"source"`
	Identifier      string `form:"identifier"`
	PublicationType string `form:"publication_type"`
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

func Add(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	switch r.URL.Query().Get("method") {
	case "identifier":
		pages.AddIdentifier(c, pages.AddIdentifierArgs{
			Step: 1,
		}).Render(r.Context(), w)
	case "manual":
		pages.AddManual(c, 1).Render(r.Context(), w)
	case "wos":
		pages.AddWebOfScience(c, 1).Render(r.Context(), w)
	case "bibtex":
		pages.AddBibTeX(c, 1).Render(r.Context(), w)
	default:
		pages.Add(c, 1).Render(r.Context(), w)
	}
}

func AddSingleImportConfirm(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindImportSingle{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("import confirm single publication: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// check for duplicates
	if b.Source == "crossref" && b.Identifier != "" {
		args := models.NewSearchArgs().
			WithFilter("identifier", strings.ToLower(b.Identifier)).
			WithFilter("status", "public")

		existing, err := c.PublicationSearchIndex.Search(args)
		if err != nil {
			c.Log.Warnw("import single publication: could not execute search for duplicates", "errors", err, "args", args, "user", c.User.ID)
			c.HandleError(w, r, httperror.InternalServerError)
			return
		}

		if existing.Total > 0 {
			pages.AddIdentifier(c, pages.AddIdentifierArgs{
				Step:                 1,
				Source:               b.Source,
				Identifier:           b.Identifier,
				DuplicatePublication: existing.Hits[0],
			}).Render(r.Context(), w)
			return
		}
	}

	AddSingleImport(w, r)
}

func AddSingleImport(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindImportSingle{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("import single publication: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	var (
		p   *models.Publication
		err error
	)

	// import by identifier
	if b.Identifier != "" {
		p, err = fetchPublicationByIdentifier(c, b.Source, b.Identifier)
		if err != nil {
			c.Log.Warnw("import single publication: could not fetch publication", "errors", err, "publication", b.Identifier, "user", c.User.ID)

			flash := flash.SimpleFlash().
				WithLevel("error").
				WithBody(template.HTML(c.Loc.Get("publication.single_import.import_by_id.import_failed")))
			c.Flash = append(c.Flash, *flash)

			pages.AddIdentifier(c, pages.AddIdentifierArgs{
				Step:       1,
				Source:     b.Source,
				Identifier: b.Identifier,
			}).Render(r.Context(), w)
			return
		}
	} else {
		// or start with empty publication
		p = &models.Publication{Type: b.PublicationType}
	}

	p.ID = ulid.Make().String()
	p.CreatorID = c.User.ID
	p.Creator = c.User
	p.UserID = c.User.ID
	p.User = c.User
	p.Status = "private"
	p.Classification = "U"

	if len(c.User.Affiliations) > 0 {
		p.AddOrganization(c.User.Affiliations[0].Organization)
	}

	if validationErrs := p.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))

		pages.AddIdentifier(c, pages.AddIdentifierArgs{
			Step:       1,
			Source:     b.Source,
			Identifier: b.Identifier,
			Errors:     errors,
		}).Render(r.Context(), w)
		return
	}

	err = c.Repo.SavePublication(p, c.User)
	if err != nil {
		c.Log.Errorf("import single publication: could not save the publication:", "error", err, "identifier", b.Identifier, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}

	pages.AddSingleDescription(c, pages.AddSingleDescriptionArgs{
		Step:         2,
		ActiveSubNav: subNav,
		Publication:  p,
	}).Render(r.Context(), w)
}

func AddSingleDescription(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}

	pages.AddSingleDescription(c, pages.AddSingleDescriptionArgs{
		Step:         2,
		Publication:  ctx.GetPublication(r),
		ActiveSubNav: subNav,
	}).Render(r.Context(), w)
}

func AddSingleConfirm(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	pages.AddSingleConfirm(c, pages.AddSingleConfirmArgs{
		Step:           3,
		Publication:    publication,
		PublicationURL: c.PathTo("publication_add_single_description", "id", publication.ID),
	}).Render(r.Context(), w)
}

func AddSinglePublish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	if !c.User.CanPublishPublication(publication) {
		c.Log.Warnw("add single publication publish: user has no permission to publish publication.", "publication", publication.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	publication.Status = "public"

	if validationErrs := publication.Validate(); validationErrs != nil {
		errors := form.Errors(localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)))
		views.ShowModal(views.FormErrorsDialog(c, "Unable to publish this publication due to the following errors", errors)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), publication, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("publication.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("add single publication publish: could not save the publication:", "errors", err, "publication", publication.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	redirectURL := views.URL(c.PathTo("publication_add_single_finish", "id", publication.ID)).Query(r.URL.Query()).String()
	w.Header().Set("HX-Redirect", redirectURL)
}

func AddSingleFinish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)

	pages.AddSingleFinish(c, pages.AddSingleFinishArgs{
		Step:           4,
		Publication:    publication,
		PublicationURL: c.PathTo("publication", "id", publication.ID),
	}).Render(r.Context(), w)
}

func AddMultipleImport(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2000000000)

	source := r.FormValue("source")

	file, _, err := r.FormFile("file")
	if err != nil {
		c.Log.Warnw("add multiple import publication: could not retrieve file from request", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}
	defer file.Close()

	// TODO why does the code imports zero entries without this?
	_, _ = file.Seek(0, io.SeekStart)

	batchID, err := importPublications(c, source, file)
	if err != nil {
		c.Log.Warnw("add multiple import publication: could not import publications", "errors", err, "batch", batchID, "user", c.User.ID)

		flash := flash.SimpleFlash().
			WithLevel("error").
			WithBody(template.HTML("<p>Sorry, something went wrong. Could not import the publication(s).</p>"))
		c.Flash = append(c.Flash, *flash)

		switch source {
		case "wos":
			pages.AddWebOfScience(c, 1).Render(r.Context(), w)
		case "bibtex":
			pages.AddBibTeX(c, 1).Render(r.Context(), w)
		}
		return
	}

	// TODO wait for index refresh, do something more elegant
	time.Sleep(1 * time.Second)

	// redirect to batch page so that pagination links work
	http.Redirect(w, r, c.PathTo("publication_add_multiple_confirm", "batch_id", batchID).String(), http.StatusFound)
}

func (h *Handler) AddMultipleSave(w http.ResponseWriter, r *http.Request, ctx Context) {
	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody(template.HTML("<p>Publications successfully saved as draft.</p>"))

	h.AddFlash(r, w, *flash)

	redirectURL := h.PathFor("publications")
	w.Header().Set("HX-Redirect", redirectURL.String())
}

// TODO after changing tabs, the wrong url is pushed in the history
func AddMultipleShow(w http.ResponseWriter, r *http.Request) {
	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}

	pages.AddMultipleShow(ctx.Get(r), pages.AddMultipleShowArgs{
		Step:         2,
		ActiveSubNav: subNav,
		RedirectURL:  r.URL.Query().Get("redirect-url"),
		Publication:  ctx.GetPublication(r),
	}).Render(r.Context(), w)
}

func AddMultipleConfirm(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	searchArgs := models.NewSearchArgs()
	if err := bind.Query(r, searchArgs); err != nil {
		c.Log.Warnw("add multiple confirm publication: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	searchArgs.WithFacetLines(vocabularies.Facets["publication"])

	batchID := bind.PathValue(r, "batch_id")

	hits, err := c.Services.PublicationSearchIndex.
		WithScope("status", "private", "public").
		WithScope("creator_id", c.User.ID).
		WithScope("batch_id", batchID).
		Search(searchArgs)

	if err != nil {
		c.Log.Errorw("add multiple confirm publication: could not execute search", "errors", err, "batch", batchID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	pages.AddMultipleConfirm(c, pages.AddMultipleConfirmArgs{
		Step:        2,
		RedirectURL: r.URL.String(),
		BatchID:     batchID,
		SearchArgs:  searchArgs,
		Hits:        hits,
	}).Render(r.Context(), w)
}

func (h *Handler) AddMultiplePublish(w http.ResponseWriter, r *http.Request, ctx Context) {
	batchID := bind.PathValue(r, "batch_id")

	err := h.batchPublishPublications(batchID, ctx.User)

	// TODO this is useless to the user unless we point to the publication in
	// question
	var validationErrs *okay.Errors
	if errors.As(err, &validationErrs) {
		h.Logger.Warnw("add multiple publish publication: could not validate abstract:", "errors", validationErrs, "batch", batchID, "user", ctx.User.ID)
		errors := form.Errors(localize.ValidationErrors(ctx.Loc, validationErrs))
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
	if err := bind.Query(r, searchArgs); err != nil {
		h.Logger.Warnw("add multiple finish publication: could not bind request arguments", "errors", err, "request", r)
		render.BadRequest(w, r, err)
		return
	}

	searchArgs.WithFacetLines(vocabularies.Facets["publication"])

	batchID := bind.PathValue(r, "batch_id")

	hits, err := h.PublicationSearchIndex.
		WithScope("status", "private", "public").
		WithScope("creator_id", ctx.User.ID).
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

func fetchPublicationByIdentifier(c *ctx.Ctx, source, identifier string) (*models.Publication, error) {
	s, ok := c.Services.PublicationSources[source]

	if !ok {
		return nil, fmt.Errorf("unkown publication source: %s", source)
	}

	d, err := s.GetPublication(identifier)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func importPublications(c *ctx.Ctx, source string, file io.Reader) (string, error) {
	batchID := ulid.Make().String()

	decFactory, ok := c.Services.PublicationDecoders[source]
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
			CreatorID:      c.User.ID,
			Creator:        c.User,
			UserID:         c.User.ID,
			User:           c.User,
		}

		if len(c.User.Affiliations) > 0 {
			p.AddOrganization(c.User.Affiliations[0].Organization)
		}

		if err := dec.Decode(&p); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			importErr = err
			break
		}

		if err := c.Repo.SavePublication(&p, c.User); err != nil {
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
func (h *Handler) batchPublishPublications(batchID string, user *models.Person) (err error) {
	searcher := h.PublicationSearchIndex.
		WithScope("status", "private", "public").
		WithScope("creator_id", user.ID).
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
			if err = h.Repo.SavePublication(pub, user); err != nil {
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
