package datasetediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	datasetviews "github.com/ugent-library/biblio-backoffice/views/dataset"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindAbstract struct {
	AbstractID string `path:"abstract_id"`
	Text       string `form:"text"`
	Lang       string `form:"lang"`
}

type BindDeleteAbstract struct {
	AbstractID string `path:"abstract_id"`
	SnapshotID string `path:"snapshot_id"`
}

func AddAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	views.ShowModal(datasetviews.EditAbstractDialog(c, datasetviews.EditAbstractDialogArgs{
		IsNew:   true,
		Dataset: ctx.GetDataset(r),
		Index:   -1,
	})).Render(r.Context(), w)
}

func CreateAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("create dataset abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	abstract := models.Text{
		Lang: b.Lang,
		Text: b.Text,
	}
	d.AddAbstract(&abstract)
	index := getAbstractIndex(d.Abstract, &abstract)

	if validationErrs := d.Validate(); validationErrs != nil {
		views.ReplaceModal(datasetviews.EditAbstractDialog(c, datasetviews.EditAbstractDialogArgs{
			IsNew:    true,
			Dataset:  d,
			Abstract: abstract,
			Index:    index,
			Errors:   validationErrs.(*okay.Errors),
		})).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(datasetviews.EditAbstractDialog(c, datasetviews.EditAbstractDialogArgs{
			IsNew:    true,
			Dataset:  d,
			Abstract: abstract,
			Index:    index,
			Conflict: true,
		})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Errorf("create dataset abstract: could not save the dataset:", "errors", err, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(datasetviews.AbstractsBodySelector, datasetviews.AbstractsBody(c, d)).Render(r.Context(), w)
}

func EditAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("edit dataset abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	abstract := d.GetAbstract(b.AbstractID)

	// TODO catch non-existing item in UI
	if abstract == nil {
		c.Log.Warnf("edit dataset abstract: Could not fetch the abstract:", "dataset", d.ID, "abstract", b.AbstractID, "user", c.User.ID)
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ShowModal(datasetviews.EditAbstractDialog(c, datasetviews.EditAbstractDialogArgs{
		IsNew:    false,
		Dataset:  d,
		Abstract: *abstract,
		Index:    getAbstractIndex(d.Abstract, abstract),
	})).Render(r.Context(), w)
}

func UpdateAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := BindAbstract{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.Log.Warnw("update dataset abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	// get pointer to abstract and manipulate in place
	d := ctx.GetDataset(r)
	abstract := d.GetAbstract(b.AbstractID)
	index := getAbstractIndex(d.Abstract, abstract)

	if abstract == nil {
		views.ReplaceModal(datasetviews.EditAbstractDialog(c, datasetviews.EditAbstractDialogArgs{
			IsNew:   false,
			Dataset: d,
			Abstract: models.Text{
				Text: b.Text,
				Lang: b.Lang,
			},
			Index:    index,
			Conflict: true,
		})).Render(r.Context(), w)
		return
	}

	abstract.Text = b.Text
	abstract.Lang = b.Lang
	d.SetAbstract(abstract)

	if validationErrs := d.Validate(); validationErrs != nil {
		views.ReplaceModal(datasetviews.EditAbstractDialog(c, datasetviews.EditAbstractDialogArgs{
			IsNew:    false,
			Dataset:  d,
			Abstract: *abstract,
			Index:    index,
			Errors:   validationErrs.(*okay.Errors),
		})).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(datasetviews.EditAbstractDialog(c, datasetviews.EditAbstractDialogArgs{
			IsNew:    false,
			Dataset:  d,
			Abstract: *abstract,
			Index:    index,
			Conflict: true,
		})).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Warnf("update dataset abstract: Could not save the dataset:", "errors", err, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(datasetviews.AbstractsBodySelector, datasetviews.AbstractsBody(c, d)).Render(r.Context(), w)
}

func ConfirmDeleteAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	d := ctx.GetDataset(r)

	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("confirm delete dataset: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	if b.SnapshotID != d.SnapshotID {
		views.ShowModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	views.ConfirmDelete(views.ConfirmDeleteArgs{
		Context:    c,
		Question:   "Are you sure you want to remove this abstract?",
		DeleteUrl:  c.PathTo("dataset_delete_abstract", "id", d.ID, "abstract_id", b.AbstractID),
		SnapshotID: d.SnapshotID,
	}).Render(r.Context(), w)
}

func DeleteAbstract(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	var b BindDeleteAbstract
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("delete datase abstract: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	d := ctx.GetDataset(r)
	d.RemoveAbstract(b.AbstractID)

	err := c.Repo.UpdateDataset(r.Header.Get("If-Match"), d, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(views.ErrorDialog(c.Loc.Get("dataset.conflict_error_reload"))).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.Log.Warnf("delete dataset abstract: Could not save the dataset:", "errors", err, "dataset", d.ID, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	views.CloseModalAndReplace(datasetviews.AbstractsBodySelector, datasetviews.AbstractsBody(c, d)).Render(r.Context(), w)
}

func getAbstractIndex(abstracts []*models.Text, abstract *models.Text) int {
	for i, a := range abstracts {
		if a.ID == abstract.ID {
			return i
		}
	}

	return -1
}
