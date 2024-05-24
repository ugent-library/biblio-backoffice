package candidaterecords

import (
	"fmt"
	"io"
	"log"
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
		c.Log.Warnw("preview candidate record: could not bind request arguments", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	rec, err := c.Repo.GetCandidateRecord(r.Context(), b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	f := rec.Publication.GetFile(bind.PathValue(r, "file_id"))

	if f == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	rc, err := c.FileStore.Get(r.Context(), f.SHA256)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(f.Name)),
	)

	io.Copy(w, rc)

	fmt.Fprintf(w, "test")
}
