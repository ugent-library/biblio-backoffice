package people

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

type responseBody[T any] struct {
	Hits struct {
		// Total int `json:"total"`
		Hits []struct {
			ID     string `json:"_id"`
			Source struct {
				Record T
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

const searchBody = `{
	"query": {{template "query" .}},
	"size": 20
}`

const matchAllQuery = `{{define "query"}}{"match_all": {}}{{end}}`

const queryStringQuery = `{{define "query"}}{
	"dis_max": {
		"queries": {
			"match": {
				"identifiers": {
					"query": {{.QueryString}},
					"operator": "AND",
					"boost": "100"
				},
				"phrase_ngram": {
					"query": {{.QueryString}},
					"operator": "AND",
					"boost": "0.05"
				},
				"ngram": {
					"query": {{.QueryString}},
					"operator": "AND",
					"boost": "0.01"
				}
			}
		}
	}
}{{end}}`

var (
	matchAllTmpl    = template.Must(template.New("").Parse(matchAllQuery + searchBody))
	queryStringTmpl = template.Must(template.New("").Parse(queryStringQuery + searchBody))
)

func (idx *Index) SearchPeople(ctx context.Context, qs string) ([]*Person, error) {
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
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return nil, err
		}
		return nil, errors.New("elasticsearch: error response: " + buf.String())
	}

	resBody := &responseBody[*Person]{}

	if err := json.NewDecoder(res.Body).Decode(resBody); err != nil {
		return nil, fmt.Errorf("elasticsearch: error parsing response body: %w", err)
	}

	recs := make([]*Person, len(resBody.Hits.Hits))

	for i, hit := range resBody.Hits.Hits {
		recs[i] = hit.Source.Record
	}

	return recs, nil
}

func (idx *Index) ReindexPeople(ctx context.Context, iter PersonIter) error {
	b, err := indexSettingsFS.ReadFile("people_index_settings.json")
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
	err = iter(ctx, func(p *Person) bool {
		doc, err := json.Marshal(newIndexPerson(p))
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

type indexPerson struct {
	Names       []string `json:"names"`
	Identifiers []string `json:"identifiers"`
	Record      *Person  `json:"record"`
}

func newIndexPerson(p *Person) *indexPerson {
	ip := &indexPerson{
		Names:       []string{p.Name},
		Identifiers: make([]string, len(p.Identifiers)),
		Record:      p,
	}

	for _, name := range []string{p.PreferredName, p.GivenName, p.PreferredGivenName, p.FamilyName, p.PreferredFamilyName} {
		if name != "" {
			ip.Names = append(ip.Names, name)
		}
	}

	for i, id := range p.Identifiers {
		ip.Identifiers[i] = id.String()
	}

	return ip
}
