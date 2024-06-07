package frontoffice

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/jpillora/ipfilter"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/frontoffice"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
	internal_time "github.com/ugent-library/biblio-backoffice/time"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httpx"
)

type Handler struct {
	Log       *slog.Logger
	Repo      *repositories.Repo
	FileStore backends.FileStore
	IPRanges  string
	IPFilter  *ipfilter.IPFilter
}

type Hits[T any] struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
	Hits   []T `json:"hits"`
}

func (h *Handler) GetPublication(w http.ResponseWriter, r *http.Request) {
	id := bind.PathValue(r, "id")
	p, err := h.Repo.GetPublication(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
		} else {
			h.Log.Error("unable to fetch publication", "id", id, "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	httpx.RenderJSON(w, 200, frontoffice.MapPublication(p, h.Repo))
}

func (h *Handler) GetDataset(w http.ResponseWriter, r *http.Request) {
	id := bind.PathValue(r, "id")
	p, err := h.Repo.GetDataset(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
		} else {
			h.Log.Error("unable to fetch dataset", "id", id, "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	httpx.RenderJSON(w, 200, frontoffice.MapDataset(p, h.Repo))
}

type BindGetAll struct {
	Limit        int    `query:"limit"`
	Offset       int    `query:"offset"`
	UpdatedSince string `query:"updated_since"`
}

func (h *Handler) GetAllPublications(w http.ResponseWriter, r *http.Request) {
	b := BindGetAll{}
	if err := bind.Query(r, &b); err != nil {
		h.Log.Error("unable to decode query", "error", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var updatedSince time.Time
	if b.UpdatedSince != "" {
		t, err := internal_time.ParseTimeUTC(b.UpdatedSince)
		if err != nil {
			h.Log.Error("unable to parse updatedSince", "error", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		updatedSince = t.Local()
	}

	n, publications, err := h.Repo.PublicationsAfter(updatedSince, b.Limit, b.Offset)
	if err != nil {
		h.Log.Error("unable to retrieve publications after", "updatedSince", updatedSince.String(), "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hits := &Hits[*frontoffice.Record]{
		Limit:  b.Limit,
		Offset: b.Offset,
		Total:  n,
		Hits:   make([]*frontoffice.Record, 0, len(publications)),
	}
	for _, p := range publications {
		hits.Hits = append(hits.Hits, frontoffice.MapPublication(p, h.Repo))
	}

	w.Header().Set("Cache-Control", "no-cache")
	httpx.RenderJSON(w, 200, hits)
}

func (h *Handler) GetAllDatasets(w http.ResponseWriter, r *http.Request) {
	b := BindGetAll{}
	if err := bind.Query(r, &b); err != nil {
		h.Log.Error("unable to decode query", "error", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var updatedSince time.Time
	if b.UpdatedSince != "" {
		t, err := internal_time.ParseTimeUTC(b.UpdatedSince)
		if err != nil {
			h.Log.Error("unable to parse updatedSince", "error", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		updatedSince = t.Local()
	}

	n, datasets, err := h.Repo.DatasetsAfter(updatedSince, b.Limit, b.Offset)
	if err != nil {
		h.Log.Error("unable to retrieve datasets after", "updatedSince", updatedSince.String(), "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hits := &Hits[*frontoffice.Record]{
		Limit:  b.Limit,
		Offset: b.Offset,
		Total:  n,
		Hits:   make([]*frontoffice.Record, 0, len(datasets)),
	}
	for _, d := range datasets {
		hits.Hits = append(hits.Hits, frontoffice.MapDataset(d, h.Repo))
	}

	w.Header().Set("Cache-Control", "no-cache")
	httpx.RenderJSON(w, 200, hits)
}

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := bind.PathValue(r, "id")
	p, err := h.Repo.GetPublication(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
		} else {
			h.Log.Error("unable to get publication", "id", id, "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if p.Status != "public" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	fileID := bind.PathValue(r, "file_id")
	f := p.GetFile(fileID)
	if f == nil {
		http.NotFound(w, r)
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
			h.Log.Warn("ip not allowed, allowed", "ip", ip, "allowed", h.IPRanges)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var reader io.ReadCloser
	var readerErr error

	if r.Method == "GET" {
		reader, readerErr = h.FileStore.Get(r.Context(), f.SHA256)
		if readerErr != nil {
			h.Log.Error("unable get file for publication", "fileID", fileID, "id", id, "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer reader.Close()
	}

	responseHeaders := [][]string{
		{"Content-Type", f.ContentType},
		{"Content-Length", fmt.Sprintf("%d", f.Size)},
		{"Last-Modified", f.DateUpdated.UTC().Format(http.TimeFormat)},
		{"ETag", f.SHA256},
		{"Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(f.Name))},
	}

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
			h.Log.Error("unable parse header If-Modified-Since", "error", err)
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
