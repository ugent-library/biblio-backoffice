package datasetcreating

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetpages "github.com/ugent-library/biblio-backoffice/views/dataset/pages"
	"github.com/ugent-library/biblio-backoffice/views/flash"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindAdd struct {
	Method string `form:"method"`
}

type BindImport struct {
	Source     string `form:"source"`
	Identifier string `form:"identifier"`
}

func Add(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindAdd{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	switch b.Method {
	case "identifier":
		datasetpages.AddIdentifier(c, datasetpages.AddIdentifierArgs{}).Render(r.Context(), w)
	case "manual":
		ConfirmImport(w, r)
		return
	default:
		datasetpages.Add(c).Render(r.Context(), w)
	}
}

func ConfirmImport(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	c.SubNav = "description" //TODO: ok?

	b := BindImport{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	// check for duplicates
	if b.Source == "datacite" {
		args := models.NewSearchArgs().WithFilter("identifier", strings.ToLower(b.Identifier)).WithFilter("status", "public")

		existing, err := c.DatasetSearchIndex.Search(args)

		if err != nil {
			c.HandleError(w, r, httperror.InternalServerError.Wrap(fmt.Errorf("could not execute search: %w", err)))
			return
		}

		if existing.Total > 0 {
			datasetpages.AddIdentifier(c, datasetpages.AddIdentifierArgs{
				Source:           b.Source,
				Identifier:       b.Identifier,
				Dataset:          existing.Hits[0],
				DuplicateDataset: true,
			}).Render(r.Context(), w)
			return
		}
	}

	AddImport(w, r)
}

func AddImport(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindImport{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	var (
		d   *models.Dataset
		err error
	)

	if b.Identifier != "" {
		d, err = fetchDatasetByIdentifier(*c.Services, b.Source, b.Identifier)
		if err != nil {
			flash := flash.SimpleFlash().
				WithLevel("error").
				WithTitle("Failed to save draft").
				WithBody(c.Loc.Get("dataset.single_import.import_by_id.import_failed"))

			c.Flash = append(c.Flash, *flash)

			datasetpages.AddIdentifier(c, datasetpages.AddIdentifierArgs{
				Source:     b.Source,
				Identifier: b.Identifier,
			}).Render(r.Context(), w)
			return
		}
	} else {
		// or start with empty dataset
		d = &models.Dataset{}
	}

	d.ID = ulid.Make().String()
	d.CreatorID = c.User.ID
	d.Creator = c.User
	d.UserID = c.User.ID
	d.User = c.User
	d.Status = "private"

	if len(c.User.Affiliations) > 0 {
		d.AddOrganization(c.User.Affiliations[0].Organization)
	}

	if validationErrs := d.Validate(); validationErrs != nil {
		datasetpages.AddIdentifier(c, datasetpages.AddIdentifierArgs{
			Source:     b.Source,
			Identifier: b.Identifier,
			Errors:     localize.ValidationErrors(c.Loc, validationErrs.(*okay.Errors)),
		}).Render(r.Context(), w)
		return
	}

	err = c.Repo.SaveDataset(d, c.User)

	if err != nil {
		c.HandleError(w, r, httperror.InternalServerError.Wrap(fmt.Errorf("could not save dataset: %w", err)))
		return
	}

	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}
	c.SubNav = subNav

	datasetpages.AddDescription(c, d).Render(r.Context(), w)
}

func AddDescription(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	subNav := r.URL.Query().Get("show")
	if subNav == "" {
		subNav = "description"
	}
	c.SubNav = subNav

	datasetpages.AddDescription(c, ctx.GetDataset(r)).Render(r.Context(), w)
}

func AddSaveDraft(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	flash := flash.SimpleFlash().
		WithLevel("success").
		WithBody("<p>Dataset successfully saved as a draft.</p>")

	c.PersistFlash(w, *flash)

	w.Header().Set("HX-Redirect", c.PathTo("datasets").String())
}

func AddConfirm(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	datasetpages.AddConfirm(c, ctx.GetDataset(r)).Render(r.Context(), w)
}

func AddPublish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	dataset.Status = "public"

	if err := dataset.Validate(); err != nil {
		errors := localize.ValidationErrors(c.Loc, err.(*okay.Errors))
		views.ShowModal(views.FormErrorsDialog("Unable to publish this dataset due to the following errors", errors)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, httperror.InternalServerError.Wrap(fmt.Errorf("could not save the dataset: %w", err)))
		return
	}

	redirectURL := c.PathTo("dataset_add_finish", "id", dataset.ID)
	redirectURL.RawQuery = r.URL.Query().Encode()

	w.Header().Set("HX-Redirect", redirectURL.String())
}

func AddFinish(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	datasetpages.AddFinish(c, ctx.GetDataset(r)).Render(r.Context(), w)
}

func fetchDatasetByIdentifier(services backends.Services, source string, identifier string) (*models.Dataset, error) {
	s, ok := services.DatasetSources[source]

	if !ok {
		return nil, fmt.Errorf("unkown dataset source: %s", source)
	}

	d, err := s.GetDataset(identifier)
	if err != nil {
		return nil, err
	}

	return d, nil
}
