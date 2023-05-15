package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
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

func (searcher *DatasetSearcher) Searcher(searchArgs *models.SearchArgs, cb func(*models.Dataset)) error {

	nProcessed := 0
	start := 0
	limit := 200

	query := buildDatasetUserQuery(searchArgs)

	queryFilters := query["query"].(M)["bool"].(M)["filter"].([]M)

	// Set the searcher scopes
	queryFilters = append(queryFilters, searcher.scopes...)

	// Set the range to ID = 0, this value gets updated with each 200 hits
	// fetched from ES in the loop
	sortValue := "0" //lowest sort value when sortin on id?
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

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "es dataset search: %s\n", buf.String())

		opts := []func(*esapi.SearchRequest){
			searcher.Client.es.Search.WithContext(context.Background()),
			searcher.Client.es.Search.WithIndex(searcher.Client.Index),
			searcher.Client.es.Search.WithTrackTotalHits(true),
			searcher.Client.es.Search.WithSort("id:asc"),
			searcher.Client.es.Search.WithBody(&buf),
		}

		var envelop datasetResEnvelope

		err := searcher.Client.SearchWithOpts(opts, &envelop)
		if err != nil {
			return err
		}

		hits, err := decodeDatasetRes(&envelop, []string{})
		if err != nil {
			return err
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
	newScopes := make([]M, 0, len(searcher.scopes))

	// Copy existing scopes
	newScopes = append(newScopes, searcher.scopes...)

	// Add new scopes
	newScopes = append(newScopes, ParseScope(field, terms...))

	return &DatasetSearcher{
		Client:  searcher.Client,
		scopes:  newScopes,
		maxSize: searcher.maxSize,
	}
}
