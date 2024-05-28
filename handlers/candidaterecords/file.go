package candidaterecords

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if !c.User.CanCurate() {
		c.HandleError(w, r, httperror.Unauthorized)
		return
	}

	b := bindCandidateRecord{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	rec, err := c.Repo.GetCandidateRecord(r.Context(), b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := rec.Publication.GetFile(bind.PathValue(r, "file_id"))

	if f == nil {
		c.HandleError(w, r, httperror.NotFound)
		return
	}

	rc, err := c.FileStore.Get(r.Context(), f.SHA256)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}
	defer rc.Close()

	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(f.Name)),
	)

	io.Copy(w, rc)
}
