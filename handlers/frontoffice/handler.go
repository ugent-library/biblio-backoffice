package frontoffice

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/jpillora/ipfilter"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/frontoffice"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/render"
	"github.com/ugent-library/biblio-backoffice/repositories"
	internal_time "github.com/ugent-library/biblio-backoffice/time"
	"github.com/ugent-library/bind"
)

type Handler struct {
	handlers.BaseHandler
	Repo             *repositories.Repo
	FileStore        backends.FileStore
	PeopleIndex      *people.Index
	IPRanges         string
	IPFilter         *ipfilter.IPFilter
	FrontendUsername string
	FrontendPassword string
}

// safe basic auth handling
// see https://www.alexedwards.net/blog/basic-authentication-in-go
func (h *Handler) BasicAuth(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(h.FrontendUsername))
			expectedPasswordHash := sha256.Sum256([]byte(h.FrontendPassword))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				fn(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type Hits struct {
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
	Total  int                   `json:"total"`
	Hits   []*frontoffice.Record `json:"hits"`
}

func (h *Handler) GetPublication(w http.ResponseWriter, r *http.Request) {
	p, err := h.Repo.GetPublication(bind.PathValue(r, "id"))
	if err != nil {
		if err == models.ErrNotFound {
			render.NotFound(w, r, err)
		} else {
			render.InternalServerError(w, r, err)
		}
		return
	}
	j, err := json.Marshal(frontoffice.MapPublication(p, h.Repo))
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (h *Handler) GetDataset(w http.ResponseWriter, r *http.Request) {
	p, err := h.Repo.GetDataset(bind.PathValue(r, "id"))
	if err != nil {
		if err == models.ErrNotFound {
			render.NotFound(w, r, err)
		} else {
			render.InternalServerError(w, r, err)
		}
		return
	}
	j, err := json.Marshal(frontoffice.MapDataset(p, h.Repo))
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	ident, err := people.NewIdentifier(bind.PathValue(r, "id"))
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	p, err := h.PeopleIndex.GetPersonByIdentifier(r.Context(), ident.Kind, ident.Value)
	if err == people.ErrNotFound {
		render.NotFound(w, r, err)
		return
	}
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	j, err := json.Marshal(p)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (h *Handler) GetActivePerson(w http.ResponseWriter, r *http.Request) {
	ident, err := people.NewIdentifier(bind.PathValue(r, "id"))
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	p, err := h.PeopleIndex.GetActivePersonByIdentifier(r.Context(), ident.Kind, ident.Value)
	if err == people.ErrNotFound {
		render.NotFound(w, r, err)
		return
	}
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	j, err := json.Marshal(p)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (h *Handler) GetActivePersonByUsername(w http.ResponseWriter, r *http.Request) {
	p, err := h.PeopleIndex.GetActivePersonByUsername(r.Context(), bind.PathValue(r, "username"))
	if err == people.ErrNotFound {
		render.NotFound(w, r, err)
		return
	}
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	j, err := json.Marshal(p)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

type BindGetAll struct {
	Limit        int    `query:"limit"`
	Offset       int    `query:"offset"`
	UpdatedSince string `query:"updated_since"`
}

func (h *Handler) GetAllPublications(w http.ResponseWriter, r *http.Request) {
	b := BindGetAll{}
	if err := bind.Query(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	var updatedSince time.Time
	if b.UpdatedSince != "" {
		t, err := internal_time.ParseTimeUTC(b.UpdatedSince)
		if err != nil {
			h.Logger.Errorw("updatedSince error", err)
			render.InternalServerError(w, r, err)
			return
		}
		updatedSince = t.Local()
	}

	mappedHits := &Hits{
		Limit:  b.Limit,
		Offset: b.Offset,
	}

	n, publications, err := h.Repo.PublicationsAfter(updatedSince, b.Limit, b.Offset)
	if err != nil {
		h.Logger.Errorw("select error", err)
		render.InternalServerError(w, r, err)
		return
	}

	mappedHits.Total = n
	mappedHits.Hits = make([]*frontoffice.Record, 0, len(publications))
	for _, p := range publications {
		mappedHits.Hits = append(mappedHits.Hits, frontoffice.MapPublication(p, h.Repo))
	}

	j, err := json.Marshal(mappedHits)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(j)
}

func (h *Handler) GetAllDatasets(w http.ResponseWriter, r *http.Request) {
	b := BindGetAll{}
	if err := bind.Query(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	var updatedSince time.Time
	if b.UpdatedSince != "" {
		t, err := internal_time.ParseTimeUTC(b.UpdatedSince)
		if err != nil {
			h.Logger.Errorw("updatedSince error", err)
			render.InternalServerError(w, r, err)
			return
		}
		updatedSince = t.Local()
	}

	mappedHits := &Hits{
		Limit:  b.Limit,
		Offset: b.Offset,
	}

	n, datasets, err := h.Repo.DatasetsAfter(updatedSince, b.Limit, b.Offset)
	if err != nil {
		h.Logger.Errorw("select error", err)
		render.InternalServerError(w, r, err)
		return
	}

	mappedHits.Total = n
	mappedHits.Hits = make([]*frontoffice.Record, 0, len(datasets))
	for _, d := range datasets {
		mappedHits.Hits = append(mappedHits.Hits, frontoffice.MapDataset(d, h.Repo))
	}

	j, err := json.Marshal(mappedHits)
	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(j)
}

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	p, err := h.Repo.GetPublication(bind.PathValue(r, "id"))
	if err != nil {
		if err == models.ErrNotFound {
			render.NotFound(w, r, err)
		} else {
			render.InternalServerError(w, r, err)
		}
		return
	}

	if p.Status != "public" {
		render.Forbidden(w, r)
		return
	}

	f := p.GetFile(bind.PathValue(r, "file_id"))
	if f == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	accessLevel := f.AccessLevel
	if accessLevel == "info:eu-repo/semantics/embargoedAccess" {
		accessLevel = f.AccessLevelDuringEmbargo
	}

	switch accessLevel {
	case "info:eu-repo/semantics/openAccess":
		// ok
	case "info:eu-repo/semantics/restrictedAccess":
		// check ip
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
			ip = remoteIP
		}
		if !h.IPFilter.Allowed(ip) {
			h.Logger.Warnw("ip not allowed, allowed", "ip", ip, "allowed", h.IPRanges)
			render.Forbidden(w, r)
			return
		}
	default:
		render.Forbidden(w, r)
		return
	}

	var reader io.ReadCloser
	var readerErr error

	if r.Method == "GET" {
		reader, readerErr = h.FileStore.Get(r.Context(), f.SHA256)
		if readerErr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer reader.Close()
	}

	responseHeaders := [][]string{}
	responseHeaders = append(responseHeaders, []string{"Content-Type", f.ContentType})
	responseHeaders = append(responseHeaders, []string{"Content-Length", fmt.Sprintf("%d", f.Size)})
	responseHeaders = append(responseHeaders, []string{"Last-Modified", f.DateUpdated.UTC().Format(http.TimeFormat)})
	responseHeaders = append(responseHeaders, []string{"ETag", f.SHA256})
	responseHeaders = append(responseHeaders, []string{"Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(f.Name))})

	/*
		Important: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/304 dictates that all
		headers in a 304 response should be sent, that would have been sent in 200 response,
		but https://github.com/golang/go/issues/6157 dictates that Content-Length
		and Content-Type should not. Weird.
	*/

	// Status 304: If-Modified-Since (Last-Modified)
	if r.Header.Get("If-Modified-Since") != "" {
		sinceTime, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// http time format does not register milliseconds
		if !f.DateUpdated.Truncate(time.Second).After(sinceTime) {
			for _, pairs := range responseHeaders {
				w.Header().Set(pairs[0], pairs[1])
			}
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// Status 304: If-None-Match (ETag)
	if r.Header.Get("If-None-Match") == f.SHA256 {
		for _, pairs := range responseHeaders {
			w.Header().Set(pairs[0], pairs[1])
		}
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Status 200
	for _, pairs := range responseHeaders {
		w.Header().Set(pairs[0], pairs[1])
	}
	w.WriteHeader(http.StatusOK)

	if r.Method == "GET" {
		io.Copy(w, reader)
	}

}
