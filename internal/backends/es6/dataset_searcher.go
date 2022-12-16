package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type DatasetSearcher struct {
	Client
	scopes  []M
	maxSize int
}

func NewDatasetSearcher(c Client, maxSize int) *DatasetSearcher {
	return &DatasetSearcher{
		Client:  c,
		maxSize: maxSize,
	}
}

func (searcher *DatasetSearcher) GetMaxSize() int {
	return searcher.maxSize
}

func (searcher *DatasetSearcher) SetMaxSize(maxSize int) {
	searcher.maxSize = maxSize
}

func (searcher *DatasetSearcher) Searcher(searchArgs *models.SearchArgs, cb func(*models.Dataset)) error {

	nProcessed := 0
	start := 0
	limit := 200
	query := searcher.buildQuery(searchArgs)

	sortValue := "0" //lowest sort value when sortin on id?
	queryFilters := query["query"].(M)["bool"].(M)["filter"].([]M)
	queryFilters = append(queryFilters, M{
		"range": M{
			"id": M{
				"gt": sortValue,
			},
		},
	})
	query["query"].(M)["bool"].(M)["filter"] = queryFilters

	for {
		//filter by range greater than instead of via from and size
		query["from"] = start
		query["size"] = limit
		queryFilters[len(queryFilters)-1]["range"].(M)["id"].(M)["gt"] = sortValue
		query["query"].(M)["bool"].(M)["filter"] = queryFilters

		opts, optsErr := searcher.buildEsOpts(query)

		if optsErr != nil {
			return optsErr
		}

		hits, hitsErr := searcher.esSearch(opts...)

		if hitsErr != nil {
			return hitsErr
		}

		if len(hits.Hits) == 0 {
			return nil
		}

		for _, hit := range hits.Hits {
			nProcessed++
			if nProcessed > searcher.maxSize {
				return nil
			}
			cb(hit)
		}

		if len(hits.Hits) > 0 {
			sortValue = hits.Hits[len(hits.Hits)-1].ID
		}

		if len(hits.Hits) < limit {
			return nil
		}
	}
}

func (searcher *DatasetSearcher) WithScope(field string, terms ...string) backends.DatasetSearcherService {
	d := searcher.Clone()
	d.scopes = append(d.scopes, ParseScope(field, terms...))
	return d
}

func (searcher *DatasetSearcher) Clone() *DatasetSearcher {
	newScopes := make([]M, 0, len(searcher.scopes))
	newScopes = append(newScopes, searcher.scopes...)
	return &DatasetSearcher{
		Client:  searcher.Client,
		scopes:  newScopes,
		maxSize: searcher.maxSize,
	}
}

func (searcher *DatasetSearcher) buildQuery(searchArgs *models.SearchArgs) M {
	query := buildDatasetUserQuery(searchArgs)
	queryFilters := query["query"].(M)["bool"].(M)["filter"].([]M)
	queryFilters = append(queryFilters, searcher.scopes...)
	query["query"].(M)["bool"].(M)["filter"] = queryFilters
	return query
}

func (searcher *DatasetSearcher) buildEsOpts(query M) ([]func(*esapi.SearchRequest), error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "es dataset search: %s\n", buf.String())
	opts := []func(*esapi.SearchRequest){
		searcher.Client.es.Search.WithContext(context.Background()),
		searcher.Client.es.Search.WithIndex(searcher.Client.Index),
		searcher.Client.es.Search.WithTrackTotalHits(true),
		searcher.Client.es.Search.WithSort("id:asc"),
		searcher.Client.es.Search.WithBody(&buf),
	}
	return opts, nil
}

func (searcher *DatasetSearcher) esSearch(opts ...func(*esapi.SearchRequest)) (*models.DatasetHits, error) {
	var envelop datasetResEnvelope
	err := searcher.Client.searchWithOpts(opts, &envelop)
	if err != nil {
		return nil, err
	}
	return decodeDatasetRes(&envelop, []string{})
}
