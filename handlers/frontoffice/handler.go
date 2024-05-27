package frontoffice

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/jpillora/ipfilter"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/frontoffice"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/projects"
	"github.com/ugent-library/biblio-backoffice/repositories"
	internal_time "github.com/ugent-library/biblio-backoffice/time"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httpx"
	"go.uber.org/zap"
)

type Handler struct {
	Log           *zap.SugaredLogger
	Repo          *repositories.Repo
	FileStore     backends.FileStore
	PeopleRepo    *people.Repo
	PeopleIndex   *people.Index
	ProjectsIndex *projects.Index
	IPRanges      string
	IPFilter      *ipfilter.IPFilter
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
			h.Log.Errorw("unable to fetch publication %s: %w", id, err)
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
			h.Log.Errorw("unable to fetch dataset %s: %w", id, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	httpx.RenderJSON(w, 200, frontoffice.MapDataset(p, h.Repo))
}

// TODO this gets way too many data
// TODO materialize sort order
// TODO constrain to those with publications
// TODO a-z sorting by id isn't the best order
func (h *Handler) GetAllOrganizations(w http.ResponseWriter, r *http.Request) {
	results, err := h.PeopleIndex.SearchOrganizations(r.Context(), people.SearchParams{Limit: 1000})
	if err != nil {
		h.Log.Errorw("unable to search organizations: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	recs := make([]*frontoffice.Organization, len(results.Hits))
	for i, o := range results.Hits {
		recs[i] = frontoffice.MapOrganization(o)
	}

	httpx.RenderJSON(w, 200, recs)
}

func (h *Handler) GetAllOrganizationTrees(w http.ResponseWriter, r *http.Request) {
	results, err := h.PeopleIndex.SearchOrganizations(r.Context(), people.SearchParams{Limit: 1000})
	if err != nil {
		h.Log.Errorw("unable to search organizations: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	recs := make([]*frontoffice.Organization, len(results.Hits))
	for i, o := range results.Hits {
		recs[i] = frontoffice.MapOrganization(o)
	}

	treeNodes := frontoffice.ToOrganizationTrees(recs)
	httpx.RenderJSON(w, 200, treeNodes)
}

func (h *Handler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	ident, err := people.NewIdentifier(bind.PathValue(r, "id"))
	if err != nil {
		h.Log.Errorw("unable to decode identifier: %w", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	o, err := h.PeopleIndex.GetOrganizationByIdentifier(r.Context(), ident.Kind, ident.Value)
	if errors.Is(err, models.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		h.Log.Errorw("unable to get organization by identifier: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	httpx.RenderJSON(w, 200, frontoffice.MapOrganization(o))
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	ident, err := people.NewIdentifier(bind.PathValue(r, "id"))
	if err != nil {
		h.Log.Errorw("unable to decode identifier: %w", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	p, err := h.PeopleIndex.GetPersonByIdentifier(r.Context(), ident.Kind, ident.Value)
	if errors.Is(err, models.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		h.Log.Errorw("unable to get person by identifier: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	httpx.RenderJSON(w, 200, frontoffice.MapPerson(p))
}

// TODO optimize
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]
	recs := make([]*frontoffice.Person, 0, len(ids))
	for _, id := range ids {
		ident, err := people.NewIdentifier(id)
		if err != nil {
			h.Log.Errorw("unable to decode identifier: %w", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		p, err := h.PeopleIndex.GetPersonByIdentifier(r.Context(), ident.Kind, ident.Value)
		if errors.Is(err, models.ErrNotFound) {
			h.Log.Warnf("unable to find person with identifier %s", ident.String())
			continue
		}
		if err != nil {
			h.Log.Errorw("unable to get person with identifier %s: %w", ident.String(), err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		recs = append(recs, frontoffice.MapPerson(p))
	}

	httpx.RenderJSON(w, 200, recs)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := bind.PathValue(r, "id")
	ident, err := people.NewIdentifier(id)
	if err != nil {
		h.Log.Errorw("unable to decode identifier %s: %w", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	p, err := h.PeopleRepo.GetActivePersonByIdentifier(r.Context(), ident.Kind, ident.Value)

	if errors.Is(err, models.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		h.Log.Errorw("unable to get person by identifier: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	httpx.RenderJSON(w, 200, frontoffice.MapPerson(p))
}

func (h *Handler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	p, err := h.PeopleRepo.GetActivePersonByUsername(r.Context(), bind.PathValue(r, "username"))
	if errors.Is(err, models.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		h.Log.Errorw("unable to get person by username: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	httpx.RenderJSON(w, 200, frontoffice.MapPerson(p))
}

type BindQuery struct {
	Limit   int      `query:"limit"`
	Offset  int      `query:"offset"`
	Query   string   `query:"q"`
	Filters []string `query:"f"`
	Sort    string   `query:"sort"`
}

func (h *Handler) SearchPeople(w http.ResponseWriter, r *http.Request) {
	b := BindQuery{}
	if err := bind.Query(r, &b); err != nil {
		h.Log.Errorw("unable to decode query: %w", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	params := people.SearchParams{
		Limit:  b.Limit,
		Offset: b.Offset,
		Query:  b.Query,
		Sort:   b.Sort,
	}
	for _, f := range b.Filters {
		if err := params.AddFilter(f); err != nil {
			h.Log.Errorw("unable to add filter from query: %w", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	results, err := h.PeopleIndex.SearchPeople(r.Context(), params)
	if err != nil {
		h.Log.Errorw("unable search people: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hits := &Hits[*frontoffice.Person]{
		Limit:  results.Limit,
		Offset: results.Offset,
		Total:  results.Total,
		Hits:   make([]*frontoffice.Person, len(results.Hits)),
	}
	for i, p := range results.Hits {
		hits.Hits[i] = frontoffice.MapPerson(p)
	}

	httpx.RenderJSON(w, 200, hits)
}

type BindSetPersonPreferredName struct {
	ID         string `path:"id"`
	GivenName  string `form:"given_name"`
	FamilyName string `form:"family_name"`
}

func (h *Handler) SetPersonPreferredName(w http.ResponseWriter, r *http.Request) {
	b := BindSetPersonPreferredName{}
	if err := bind.Request(r, &b); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ident, err := people.NewIdentifier(b.ID)
	if err != nil {
		h.Log.Errorw("unable to decode identifier %s: %w", b.ID, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.PeopleRepo.SetPersonPreferredName(r.Context(), ident.Kind, ident.Value, b.GivenName, b.FamilyName)
	if err != nil {
		h.Log.Errorw("unable to set person preferred name: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := bind.PathValue(r, "id")
	ident, err := projects.NewIdentifier(id) // TODO don't use function from people ns
	if err != nil {
		h.Log.Errorw("unable to decode identifier %s: %w", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	p, err := h.ProjectsIndex.GetProjectByIdentifier(r.Context(), ident.Kind, ident.Value)
	if errors.Is(err, models.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		h.Log.Errorw("unable to get project by identifier: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	httpx.RenderJSON(w, 200, frontoffice.MapProject(p))
}

// TODO constrain to those with publications
func (h *Handler) BrowseProjects(w http.ResponseWriter, r *http.Request) {
	b := BindQuery{}
	if err := bind.Query(r, &b); err != nil {
		h.Log.Errorw("unable to decode query: %w", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	results, err := h.ProjectsIndex.BrowseProjects(r.Context(), projects.SearchParams{
		Query:  b.Query,
		Limit:  b.Limit,
		Offset: b.Offset,
	})
	if err != nil {
		h.Log.Errorw("unable to search projects: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hits := &Hits[*frontoffice.Project]{
		Limit:  results.Limit,
		Offset: results.Offset,
		Total:  results.Total,
		Hits:   make([]*frontoffice.Project, len(results.Hits)),
	}
	for i, p := range results.Hits {
		hits.Hits[i] = frontoffice.MapProject(p)
	}

	httpx.RenderJSON(w, 200, hits)
}

type BindGetAll struct {
	Limit        int    `query:"limit"`
	Offset       int    `query:"offset"`
	UpdatedSince string `query:"updated_since"`
}

func (h *Handler) GetAllPublications(w http.ResponseWriter, r *http.Request) {
	b := BindGetAll{}
	if err := bind.Query(r, &b); err != nil {
		h.Log.Errorw("unable to decode query: %w", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var updatedSince time.Time
	if b.UpdatedSince != "" {
		t, err := internal_time.ParseTimeUTC(b.UpdatedSince)
		if err != nil {
			h.Log.Errorw("unable to parse updatedSince: %w", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		updatedSince = t.Local()
	}

	n, publications, err := h.Repo.PublicationsAfter(updatedSince, b.Limit, b.Offset)
	if err != nil {
		h.Log.Errorw("unable to retrieve publications after %s: %w", updatedSince.String(), err)
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
		h.Log.Errorw("unable to decode query: %w", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var updatedSince time.Time
	if b.UpdatedSince != "" {
		t, err := internal_time.ParseTimeUTC(b.UpdatedSince)
		if err != nil {
			h.Log.Errorw("unable to parse updatedSince: %w", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		updatedSince = t.Local()
	}

	n, datasets, err := h.Repo.DatasetsAfter(updatedSince, b.Limit, b.Offset)
	if err != nil {
		h.Log.Errorw("unable to retrieve datasets after %s: %w", updatedSince.String(), err)
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
			h.Log.Errorw("unable to get publication %s: %w", id, err)
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
			h.Log.Warnw("ip not allowed, allowed", "ip", ip, "allowed", h.IPRanges)
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
			h.Log.Errorw("unable get file %s for publication %s: %w", fileID, id, err)
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
			h.Log.Errorw("unable parse header If-Modified-Since: %w", err)
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
