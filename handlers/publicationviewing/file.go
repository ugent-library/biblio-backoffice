package publicationviewing

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/bind"
)

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request, ctx Context) {
	f := ctx.Publication.GetFile(bind.PathValues(r).Get("file_id"))

	if f == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	b, err := h.FileStore.Get(r.Context(), f.SHA256)
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
