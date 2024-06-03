package projects

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"text/template"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	"github.com/ugent-library/biblio-backoffice/models"
	index "github.com/ugent-library/index/es6"
)

//go:embed *.json
var indexSettingsFS embed.FS

const (
	projectsIndexName = "projects"
)

type IndexConfig struct {
	Conn        string
	IndexPrefix string
	Retention   int
	Logger      *slog.Logger
}

type Index struct {
	client    *elasticsearch.Client
	prefix    string
	retention int
	logger    *slog.Logger
}

func NewIndex(c IndexConfig) (*Index, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{c.Conn},
	})
	if err != nil {
		return nil, err
	}

	return &Index{
		client:    client,
		prefix:    c.IndexPrefix,
		retention: c.Retention,
		logger:    c.Logger,
	}, nil
}

type searchResponseBody[T any] struct {
	Hits struct {
		Total int `json:"total"`
		Hits  []struct {
			Source struct {
				Record T `json:"record"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func decodeResponseBody[T any](res *esapi.Response, resBody T) error {
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return err
		}
		return errors.New("elasticsearch: error response: " + buf.String())
	}

	if err := json.NewDecoder(res.Body).Decode(resBody); err != nil {
		return fmt.Errorf("elasticsearch: error parsing response body: %w", err)
	}

	return nil
}

const searchBody = `{
	"query": {{template "query" .}},
	"size": {{.Limit}},
	"from": {{.Offset}},
	"_source": ["record"]
}`

const identifierQuery = `{{define "query"}}{
	"bool": {
		"filter": [
			{"term": {"identifiers": "{{.Identifier.Value}}"}}
		]
	}
}{{end}}`

const queryStringQuery = `{{define "query"}}{
	{{if .Query}}
	"dis_max": {
		"queries": [
			{
				"match": {
					"identifiers": {
						"query": "{{.Query}}",
						"operator": "AND",
						"boost": "100"
					}
				}
			},
			{
				"match": {
					"all": {
						"query": "{{.Query}}",
						"operator": "AND",
						"boost": "0.1"
					}
				}
			},
			{
				"match": {
					"phrase_ngram": {
						"query": "{{.Query}}",
						"operator": "AND",
						"boost": "0.05"
					}
				}
			},
			{
				"match": {
					"ngram": {
						"query": "{{.Query}}",
						"operator": "AND",
						"boost": "0.01"
					}
				}
			}
		]
	}
	{{else}}
	"match_all": {}
	{{end}}
}{{end}}`

const browseNameQuery = `{{define "query"}}{
	"bool": {
		"filter": [
			{"term": {"nameKey": "{{.Query}}"}}
		]
	}
}{{end}}`

var (
	identifierTmpl  = template.Must(template.New("").Parse(identifierQuery + searchBody))
	queryStringTmpl = template.Must(template.New("").Parse(queryStringQuery + searchBody))
	browseNameTmpl  = template.Must(template.New("").Parse(browseNameQuery + searchBody))
)

func (idx *Index) GetProjectByIdentifier(ctx context.Context, kind, value string) (*Project, error) {
	ident := Identifier{Kind: kind, Value: value}

	b := bytes.Buffer{}
	err := identifierTmpl.Execute(&b, struct {
		Limit      int
		Offset     int
		Identifier Identifier
	}{
		Limit:      1,
		Identifier: ident,
	})
	if err != nil {
		return nil, fmt.Errorf("index.GetProjectByIdentifier: identifier %q: %w", ident.String(), err)

	}

	res, err := idx.client.Search(
		idx.client.Search.WithContext(ctx),
		idx.client.Search.WithIndex(idx.prefix+projectsIndexName),
		idx.client.Search.WithTrackTotalHits(false),
		idx.client.Search.WithBody(strings.NewReader(b.String())),
	)
	if err != nil {
		return nil, fmt.Errorf("index.GetProjectByIdentifier: identifier %q: %w", ident.String(), err)
	}

	resBody := searchResponseBody[*Project]{}
	if err := decodeResponseBody(res, &resBody); err != nil {
		return nil, fmt.Errorf("index.GetProjectByIdentifier: identifier %q: %w", ident.String(), err)
	}

	if len(resBody.Hits.Hits) != 1 {
		return nil, fmt.Errorf("index.GetProjectByIdentifier: identifier %q: %w", ident.String(), models.ErrNotFound)
	}

	return resBody.Hits.Hits[0].Source.Record, nil
}

func (idx *Index) SearchProjects(ctx context.Context, params SearchParams) (*SearchResults[*Project], error) {
	results, err := search(ctx, idx, queryStringTmpl, params, "_score:desc")
	if err != nil {
		return nil, fmt.Errorf("index.SearchProjects: %w", err)
	}
	return results, nil
}

func (idx *Index) BrowseProjects(ctx context.Context, params SearchParams) (*SearchResults[*Project], error) {
	results, err := search(ctx, idx, browseNameTmpl, params, "sortName:asc")
	if err != nil {
		return nil, fmt.Errorf("index.BrowseProjects: %w", err)
	}
	return results, nil
}

func search(ctx context.Context, idx *Index, tmpl *template.Template, params SearchParams, sort string) (*SearchResults[*Project], error) {
	b := bytes.Buffer{}
	err := tmpl.Execute(&b, params)
	if err != nil {
		return nil, err
	}

	res, err := idx.client.Search(
		idx.client.Search.WithContext(ctx),
		idx.client.Search.WithIndex(idx.prefix+projectsIndexName),
		idx.client.Search.WithTrackTotalHits(true),
		idx.client.Search.WithBody(strings.NewReader(b.String())),
		idx.client.Search.WithSort(sort),
	)
	if err != nil {
		return nil, err
	}

	resBody := searchResponseBody[*Project]{}
	if err := decodeResponseBody(res, &resBody); err != nil {
		return nil, err
	}

	results := SearchResults[*Project]{
		Limit:  params.Limit,
		Offset: params.Offset,
		Total:  resBody.Hits.Total,
		Hits:   make([]*Project, len(resBody.Hits.Hits)),
	}

	for i, hit := range resBody.Hits.Hits {
		results.Hits[i] = hit.Source.Record
	}

	return &results, nil
}

func (idx *Index) ReindexProjects(ctx context.Context, iter ProjectIter) error {
	b, err := indexSettingsFS.ReadFile(projectsIndexName + "_index_settings.json")
	if err != nil {
		return fmt.Errorf("index.ReindexProjects: %w", err)
	}

	switcher, err := index.NewSwitcher(idx.client, idx.prefix+projectsIndexName, string(b))
	if err != nil {
		return fmt.Errorf("index.ReindexProjects: %w", err)
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:  idx.client,
		Index:   switcher.Name(),
		Refresh: "true",
		OnError: func(ctx context.Context, err error) {
			idx.logger.ErrorContext(ctx, "index error", slog.Any("error", err))
		},
	})
	if err != nil {
		return fmt.Errorf("index.ReindexProjects: %w", err)
	}
	defer bi.Close(ctx)

	var indexErr error
	err = iter(ctx, func(p *Project) bool {
		docID, doc, err := toProjectDoc(p)
		if err != nil {
			indexErr = err
			return false
		}
		indexErr = bi.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action:       "index",
				DocumentID:   docID,
				DocumentType: "_doc",
				Body:         bytes.NewReader(doc),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						err = fmt.Errorf("index error: %v", err)
					} else {
						err = fmt.Errorf("index error: %s: %s", res.Error.Type, res.Error.Reason)
					}
					idx.logger.ErrorContext(ctx, "index failure", slog.String("id", item.DocumentID), slog.Any("error", err))
				},
			},
		)
		return indexErr == nil
	})
	if err != nil {
		return fmt.Errorf("index.ReindexProjects: %w", err)
	}
	if indexErr != nil {
		return fmt.Errorf("index.ReindexProjects: %w", indexErr)
	}

	if err := switcher.Switch(ctx, idx.retention); err != nil {
		return fmt.Errorf("index.ReindexProjects: %w", err)
	}

	return nil
}

type projectDoc struct {
	NameKey      string   `json:"nameKey"`
	SortName     string   `json:"sortName"`
	Names        []string `json:"names"`
	Descriptions []string `json:"descriptions"`
	Identifiers  []string `json:"identifiers"`
	Deleted      bool     `json:"deleted"`
	Record       *Project `json:"record"`
}

func fixSortName(sort string) string {
	sort = strings.ToLower(sort)
	sort = strings.ReplaceAll(sort, " ", "")
	sort = strings.ReplaceAll(sort, "'", "")
	sort = strings.ReplaceAll(sort, "-", "")
	return sort
}

func toProjectDoc(p *Project) (string, []byte, error) {
	pd := &projectDoc{
		Names:        make([]string, len(p.Names)),
		Descriptions: make([]string, len(p.Descriptions)),
		Identifiers:  make([]string, len(p.Identifiers)),
		Deleted:      p.Deleted,
		Record:       p,
	}

	if name := p.Names.Get("und"); name != "" {
		pd.SortName = name
		pd.NameKey = name[0:1]
	}

	if pd.SortName != "" {
		pd.SortName = fixSortName(pd.SortName)
	}

	for i, name := range p.Names {
		pd.Names[i] = name.Value
	}

	for i, desc := range p.Descriptions {
		pd.Descriptions[i] = desc.Value
	}

	for i, id := range p.Identifiers {
		pd.Identifiers[i] = id.Value
	}

	doc, err := json.Marshal(pd)
	if err != nil {
		return "", nil, err
	}

	return p.Identifiers.Get(idKind), doc, nil
}
