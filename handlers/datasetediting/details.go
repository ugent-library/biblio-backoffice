package datasetediting

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindDetails struct {
	AccessLevel             string   `form:"access_level"`
	AccessLevelAfterEmbargo string   `form:"access_level_after_embargo"`
	EmbargoDate             string   `form:"embargo_date"`
	Format                  []string `form:"format"`
	Identifier              string   `form:"identifier"`
	IdentifierType          string   `form:"identifier_type"`
	Keyword                 []string `form:"keyword"`
	Language                []string `form:"language"`
	License                 string   `form:"license"`
	OtherLicense            string   `form:"other_license"`
	Publisher               string   `form:"publisher"`
	Title                   string   `form:"title"`
	Year                    string   `form:"year"`
}

func EditDetails(w http.ResponseWriter, r *http.Request) {
	views.ShowModal(datasetviews.EditDetailsDialog(
		ctx.Get(r), ctx.GetDataset(r), false, nil,
	)).Render(r.Context(), w)
}

func RefreshEditDetails(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindDetails{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update dataset details: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		b.EmbargoDate = ""
		b.AccessLevelAfterEmbargo = ""
	}

	dataset.AccessLevel = b.AccessLevel
	dataset.AccessLevelAfterEmbargo = b.AccessLevelAfterEmbargo
	dataset.EmbargoDate = b.EmbargoDate
	dataset.Format = b.Format
	dataset.Identifiers = models.Values{b.IdentifierType: []string{b.Identifier}}
	dataset.Keyword = b.Keyword
	dataset.Language = b.Language
	dataset.License = b.License
	dataset.OtherLicense = b.OtherLicense
	dataset.Publisher = b.Publisher
	dataset.Title = b.Title
	dataset.Year = b.Year

	views.ReplaceModal(datasetviews.EditDetailsDialog(c, dataset, false, nil)).Render(r.Context(), w)
}

func UpdateDetails(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	dataset := ctx.GetDataset(r)

	b := BindDetails{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update dataset details: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// @note decoding the form into a model omits empty values
	//   removing "omitempty" in the model doesn't make a difference.
	if b.AccessLevel != "info:eu-repo/semantics/embargoedAccess" {
		b.EmbargoDate = ""
		b.AccessLevelAfterEmbargo = ""
	}

	dataset.AccessLevel = b.AccessLevel
	dataset.EmbargoDate = b.EmbargoDate
	dataset.AccessLevelAfterEmbargo = b.AccessLevelAfterEmbargo
	dataset.Format = b.Format
	dataset.Language = b.Language
	dataset.Keyword = b.Keyword
	dataset.Identifiers = models.Values{b.IdentifierType: []string{b.Identifier}}
	dataset.License = b.License
	dataset.OtherLicense = b.OtherLicense
	dataset.Publisher = b.Publisher
	dataset.Title = b.Title
	dataset.Year = b.Year

	validationErrs := dataset.Validate()
	// check EmbargoDate is in the future at time of submit
	if dataset.EmbargoDate != "" {
		t, e := time.Parse("2006-01-02", dataset.EmbargoDate)
		if e == nil && !t.After(time.Now()) {
			validationErrs = okay.Add(validationErrs, okay.NewError("/embargo_date", "dataset.embargo_date.expired"))
		}
	}

	if validationErrs != nil {
		views.ReplaceModal(datasetviews.EditDetailsDialog(
			c, dataset, false, validationErrs.(*okay.Errors),
		)).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), dataset, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(datasetviews.EditDetailsDialog(
			c, dataset, true, nil,
		)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("update dataset details: Could not save the dataset:", "errors", err, "dataset", dataset.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace("#details-body", datasetviews.DetailsBody(c, dataset)).Render(r.Context(), w)
}
