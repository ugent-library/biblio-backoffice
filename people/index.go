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

	"github.com/elastic/go-elasticsearch/v6"
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
		Total int `json:"total"`
		Hits  []struct {
			ID     string `json:"_id"`
			Source struct {
				Record T
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

var boosts = map[string]string{
	"identifiers":  "100",
	"phrase_ngram": "0.05",
	"ngram":        "0.01",
}

func (idx *Index) SearchPeople(ctx context.Context, q string) ([]*Person, error) {
	query := map[string]any{
		"match_all": map[string]any{},
	}

	if q = strings.TrimSpace(q); q != "" {
		dismaxQueries := make([]map[string]any, 0, len(boosts))
		for field, boost := range boosts {
			dismaxQuery := map[string]any{
				"match": map[string]any{
					field: map[string]any{
						"query":    q,
						"operator": "AND",
						"boost":    boost,
					},
				},
			}
			dismaxQueries = append(dismaxQueries, dismaxQuery)
		}
		query = map[string]any{
			"dis_max": map[string]any{
				"queries": dismaxQueries,
			},
		}
	}

	reqBody := map[string]any{
		"query": query,
		"size":  20,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return nil, err
	}

	res, err := idx.client.Search(
		idx.client.Search.WithContext(ctx),
		idx.client.Search.WithIndex(idx.alias),
		idx.client.Search.WithTrackTotalHits(true),
		idx.client.Search.WithBody(&buf),
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

	indexer, err := index.NewIndexer(idx.client, switcher.Name(), index.IndexerConfig{
		OnError: func(err error) {
			idx.logger.ErrorContext(ctx, "index error", slog.Any("error", err))
		},
		OnIndexFailure: func(docID string, err error) {
			idx.logger.ErrorContext(ctx, "index failure", slog.String("doc_id", docID), slog.Any("error", err))
		},
	})
	if err != nil {
		return err
	}
	defer indexer.Close(ctx)

	var indexErr error
	err = iter(ctx, func(p *Person) bool {
		doc, err := json.Marshal(newIndexPerson(p))
		if err != nil {
			indexErr = err
			return false
		}
		indexErr = indexer.Index(ctx, p.Identifiers[0].String(), doc)
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
