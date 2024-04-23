package publicationviewing

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/bind"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	publication := ctx.GetPublication(r)
	f := publication.GetFile(bind.PathValue(r, "file_id"))

	if f == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	b, err := c.FileStore.Get(r.Context(), f.SHA256)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer b.Close()

	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(f.Name)),
	)

	io.Copy(w, b)
}
