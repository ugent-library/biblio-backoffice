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
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	index "github.com/ugent-library/index/es6"
)

//go:embed *.json
var indexSettingsFS embed.FS

const (
	peopleIndexName        = "people"
	organizationsIndexName = "organizations"
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

func (idx *Index) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	return get[Organization](ctx, idx, organizationsIndexName, id)
}

func (idx *Index) GetPerson(ctx context.Context, id string) (*Person, error) {
	return get[Person](ctx, idx, peopleIndexName, id)
}

func get[T any](_ context.Context, idx *Index, indexName string, id string) (*T, error) {
	res, err := idx.client.Get(idx.prefix+indexName, id,
		idx.client.Get.WithSource("record"),
	)
	if err != nil {
		return nil, err
	}
	resBody := getResponseBody[*T]{}
	if err := decodeResponseBody(res, &resBody); err != nil {
		return nil, err
	}
	return resBody.Source.Record, nil
}

func (idx *Index) SearchOrganizations(ctx context.Context, qs string) ([]*Organization, error) {
	return search[Organization](ctx, idx, organizationsIndexName, qs)
}

func (idx *Index) SearchPeople(ctx context.Context, qs string) ([]*Person, error) {
	return search[Person](ctx, idx, peopleIndexName, qs)
}

func search[T any](ctx context.Context, idx *Index, indexName string, qs string) ([]*T, error) {
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
		idx.client.Search.WithIndex(idx.prefix+indexName),
		// idx.client.Search.WithTrackTotalHits(true),
		idx.client.Search.WithBody(strings.NewReader(b.String())),
		idx.client.Search.WithSort("_score:desc"),
		idx.client.Search.WithSource("record"),
	)
	if err != nil {
		return nil, err
	}

	resBody := searchResponseBody[*T]{}
	if err := decodeResponseBody(res, &resBody); err != nil {
		return nil, err
	}

	recs := make([]*T, len(resBody.Hits.Hits))

	for i, hit := range resBody.Hits.Hits {
		recs[i] = hit.Source.Record
	}

	return recs, nil
}

func (idx *Index) ReindexOrganizations(ctx context.Context, iter Iter[*Organization]) error {
	return reindex(ctx, idx, organizationsIndexName, iter, toOrganizationDoc)
}

func (idx *Index) ReindexPeople(ctx context.Context, iter Iter[*Person]) error {
	return reindex(ctx, idx, peopleIndexName, iter, toPersonDoc)
}

func reindex[T any](ctx context.Context, idx *Index, indexName string, iter Iter[T], docFn func(T) (string, []byte, error)) error {
	b, err := indexSettingsFS.ReadFile(indexName + "_index_settings.json")
	if err != nil {
		return err
	}

	switcher, err := index.NewSwitcher(idx.client, idx.prefix+indexName, string(b))
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
	err = iter(ctx, func(p T) bool {
		docID, doc, err := docFn(p)
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
		return err
	}
	if indexErr != nil {
		return indexErr
	}

	return switcher.Switch(ctx, idx.retention)
}

type organizationDoc struct {
	Names       []string      `json:"names"`
	Identifiers []string      `json:"identifiers"`
	Record      *Organization `json:"record"`
}

func toOrganizationDoc(o *Organization) (string, []byte, error) {
	od := &organizationDoc{
		Names:       make([]string, 0, len(o.Names)),
		Identifiers: make([]string, 0, len(o.Identifiers)*2),
		Record:      o,
	}

	for _, text := range o.Names {
		od.Names = append(od.Names, text.Value)
	}

	for _, id := range o.Identifiers {
		od.Identifiers = append(od.Identifiers, id.String(), id.Value)
	}

	doc, err := json.Marshal(od)
	if err != nil {
		return "", nil, err
	}

	return o.Identifiers.Get(idKind), doc, nil
}

type personDoc struct {
	Names       []string `json:"names"`
	Identifiers []string `json:"identifiers"`
	Record      *Person  `json:"record"`
}

func toPersonDoc(p *Person) (string, []byte, error) {
	pd := &personDoc{
		Names:       []string{p.Name},
		Identifiers: make([]string, 0, len(p.Identifiers)*2),
		Record:      p,
	}

	for _, name := range []string{p.PreferredName, p.GivenName, p.PreferredGivenName, p.FamilyName, p.PreferredFamilyName} {
		if name != "" {
			pd.Names = append(pd.Names, name)
		}
	}

	for _, id := range p.Identifiers {
		pd.Identifiers = append(pd.Identifiers, id.String(), id.Value)
	}

	doc, err := json.Marshal(pd)
	if err != nil {
		return "", nil, err
	}

	return p.Identifiers.Get(idKind), doc, nil
}
