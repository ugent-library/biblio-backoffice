package projects

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	index "github.com/ugent-library/index/es6"
)

//go:embed *.json
var indexSettingsFS embed.FS

type IndexConfig struct {
	Conn      string
	Name      string
	Retention int
	Logger    *slog.Logger
}

type Index struct {
	client    *elasticsearch.Client
	alias     string
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
		alias:     c.Name,
		retention: c.Retention,
		logger:    c.Logger,
	}, nil
}

type getResponseBody[T any] struct {
	Source struct {
		Record T `json:"record"`
	} `json:"_source"`
}

type searchResponseBody[T any] struct {
	Hits struct {
		// Total int `json:"total"`
		Hits []struct {
			ID     string `json:"_id"`
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
	"size": 20
}`

const matchAllQuery = `{{define "query"}}{"match_all": {}}{{end}}`

const queryStringQuery = `{{define "query"}}{
	"dis_max": {
		"queries": [
			{
				"match": {
					"identifiers": {
						"query": "{{.QueryString}}",
						"operator": "AND",
						"boost": "100"
					}
				}
			},
			{
				"match": {
					"phrase_ngram": {
						"query": "{{.QueryString}}",
						"operator": "AND",
						"boost": "0.05"
					}
				}
			},
			{
				"match": {
					"ngram": {
						"query": "{{.QueryString}}",
						"operator": "AND",
						"boost": "0.01"
					}
				}
			}
		]
	}
}{{end}}`

var (
	matchAllTmpl    = template.Must(template.New("").Parse(matchAllQuery + searchBody))
	queryStringTmpl = template.Must(template.New("").Parse(queryStringQuery + searchBody))
)

func (idx *Index) GetProject(ctx context.Context, id string) (*Project, error) {
	res, err := idx.client.Get(idx.alias, id,
		idx.client.Get.WithSource("record"),
	)
	if err != nil {
		return nil, err
	}
	resBody := getResponseBody[*Project]{}
	if err := decodeResponseBody(res, &resBody); err != nil {
		return nil, err
	}
	return resBody.Source.Record, nil
}

func (idx *Index) SearchProjects(ctx context.Context, qs string) ([]*Project, error) {
	qs = strings.TrimSpace(qs)

	b := bytes.Buffer{}
	tmpl := matchAllTmpl
	if qs != "" {
		tmpl = queryStringTmpl
	}

	err := tmpl.Execute(&b, struct {
		QueryString string
	}{
		QueryString: qs,
	})
	if err != nil {
		return nil, err
	}

	res, err := idx.client.Search(
		idx.client.Search.WithContext(ctx),
		idx.client.Search.WithIndex(idx.alias),
		// idx.client.Search.WithTrackTotalHits(true),
		idx.client.Search.WithBody(strings.NewReader(b.String())),
		idx.client.Search.WithSort("_score:desc"),
		idx.client.Search.WithSource("record"),
	)
	if err != nil {
		return nil, err
	}

	resBody := searchResponseBody[*Project]{}
	if err := decodeResponseBody(res, &resBody); err != nil {
		return nil, err
	}

	recs := make([]*Project, len(resBody.Hits.Hits))

	for i, hit := range resBody.Hits.Hits {
		recs[i] = hit.Source.Record
	}

	return recs, nil
}

func (idx *Index) ReindexProjects(ctx context.Context, iter ProjectIter) error {
	b, err := indexSettingsFS.ReadFile("projects_index_settings.json")
	if err != nil {
		return err
	}

	switcher, err := index.NewSwitcher(idx.client, idx.alias, string(b))
	if err != nil {
		return err
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
		return err
	}
	defer bi.Close(ctx)

	var indexErr error
	err = iter(ctx, func(p *Project) bool {
		doc, err := json.Marshal(newIndexProject(p))
		if err != nil {
			indexErr = err
			return false
		}
		indexErr = bi.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action:       "index",
				DocumentID:   p.ID(),
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
		return err
	}
	if indexErr != nil {
		return indexErr
	}

	return switcher.Switch(ctx, idx.retention)
}

type indexProject struct {
	Names        []string `json:"names"`
	Descriptions []string `json:"descriptions"`
	Identifiers  []string `json:"identifiers"`
	Record       *Project `json:"record"`
}

func newIndexProject(p *Project) *indexProject {
	ip := &indexProject{
		Names:        make([]string, len(p.Names)),
		Descriptions: make([]string, len(p.Descriptions)),
		Identifiers:  make([]string, len(p.Identifiers)),
		Record:       p,
	}

	for i, name := range p.Names {
		ip.Names[i] = name.Value
	}

	for i, desc := range p.Descriptions {
		ip.Descriptions[i] = desc.Value
	}

	for i, id := range p.Identifiers {
		ip.Identifiers[i] = id.Value
	}

	return ip
}
