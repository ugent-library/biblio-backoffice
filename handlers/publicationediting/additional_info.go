package publicationediting

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/snapstore"
	"github.com/ugent-library/biblio-backoffice/views"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type BindAdditionalInfo struct {
	AdditionalInfo string   `form:"additional_info"`
	Keyword        []string `form:"keyword"`
	ResearchField  []string `form:"research_field"`
}

func EditAdditionalInfo(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)
	views.ShowModal(publicationviews.EditAdditionalInfoDialog(c, p, false, nil)).Render(r.Context(), w)
}

func UpdateAdditionalInfo(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	p := ctx.GetPublication(r)

	b := BindAdditionalInfo{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	p.AdditionalInfo = b.AdditionalInfo
	p.Keyword = b.Keyword
	p.ResearchField = b.ResearchField

	if validationErrs := p.Validate(); validationErrs != nil {
		c.Log.Warnw("update publication additional info: could not validate additional info:", "errors", validationErrs, "publication", p.ID, "user", c.User.ID)
		views.ReplaceModal(publicationviews.EditAdditionalInfoDialog(c, p, false, validationErrs.(*okay.Errors))).Render(r.Context(), w)
		return
	}

	err := c.Repo.UpdatePublication(r.Header.Get("If-Match"), p, c.User)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		views.ReplaceModal(publicationviews.EditAdditionalInfoDialog(c, p, true, nil)).Render(r.Context(), w)
		return
	}

	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.CloseModalAndReplace(publicationviews.AdditionalInfoBodySelector, publicationviews.AdditionalInfoBody(c, p)).Render(r.Context(), w)
}
