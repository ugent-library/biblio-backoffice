package publicationviewing

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backoffice/internal/bind"
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

// func (h *Handler) FileThumbnail(w http.ResponseWriter, r *http.Request, ctx Context) {
// 	f := ctx.Publication.GetFile(bind.PathValues(r).Get("file_id"))

// 	if f == nil {
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	if f.SHA256 == "" {
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	if f.ContentType != "application/pdf" || f.Size > 25000000 {
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	filePath := h.FileStore.RelativeFilePath(f.SHA256)

// 	params := imagorpath.Params{
// 		Image:  filePath,
// 		FitIn:  true,
// 		Width:  156,
// 		Height: 156,
// 	}
// 	imgPath := imagorpath.Generate(params, imagorpath.NewDefaultSigner(viper.GetString("imagor-secret")))
// 	imagorURL, _ := url.Parse(viper.GetString("imagor-url"))
// 	imgURL, _ := url.Parse(viper.GetString("imagor-url"))
// 	imgURL.Path = imgPath

// 	r.URL.Host = imgURL.Host
// 	r.URL.Scheme = imgURL.Scheme
// 	r.URL.Path = strings.Replace(imgURL.Path, imagorURL.Path, "", 1)
// 	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
// 	r.Header.Del("Cookie")
// 	r.Host = imgURL.Host

// 	proxy := httputil.NewSingleHostReverseProxy(imagorURL)
// 	proxy.ServeHTTP(w, r)
// }
