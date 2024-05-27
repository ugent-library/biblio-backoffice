package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
)

type DatasetIndex struct {
	client *elasticsearch.Client
	index  string
	scopes []M
}

func newDatasetIndex(c *elasticsearch.Client, i string) backends.DatasetIDIndex {
	return &DatasetIndex{
		client: c,
		index:  i,
	}
}

func (di *DatasetIndex) Search(args *models.SearchArgs) (*models.SearchHits, error) {
	// BUILD QUERY AND FILTERS FROM USER INPUT
	query := buildDatasetUserQuery(args)

	queryFilters := query["query"].(M)["bool"].(M)["filter"].([]M)
	queryMust := query["query"].(M)["bool"].(M)["must"].(M)

	// extra internal filters
	// internalFilters := []M{
	// 	{
	// 		"bool": M{
	// 			"must_not": M{
	// 				"exists": M{
	// 					"field": "date_until",
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// FACETS
	// 	create global bucket so that not all buckets are influenced by query and filters
	// 	name "facets" is not important
	if args.Facets != nil {
		query["aggs"] = M{
			"facets": M{
				"global": M{},
				"aggs":   M{},
			},
		}

		// facet filter contains all query and all filters except itself
		for _, field := range args.Facets {
			filters := make([]M, 0, len(di.scopes)+1)

			// add all internal filters
			filters = append(filters, queryMust)
			filters = append(filters, di.scopes...)
			// filters = append(filters, internalFilters...)

			// TODO: cleanup messy difference between regular filters and
			// facet based filters (based on the existence of "terms")
			for _, filter := range queryFilters {
				terms := filter["terms"]
				//regular filter
				if terms == nil {
					filters = append(filters, filter)
					continue
				}
				//facet based filter: add filter only if not matching
				if _, found := terms.(M)[field]; found {
					continue
				} else {
					filters = append(filters, filter)
				}
			}

			facet := M{
				"filter": M{"bool": M{"must": filters}},
				"aggs": M{
					"facet": M{
						"terms": M{
							"field":         field,
							"order":         M{"_key": "asc"},
							"size":          200,
							"min_doc_count": 0,
						},
					},
				},
			}

			if includeFields, e := fixedFacetValues[field]; e {
				facet["aggs"].(M)["facet"].(M)["terms"].(M)["include"] = includeFields
			}

			query["aggs"].(M)["facets"].(M)["aggs"].(M)[field] = facet
		}

		// Dynamically add variable facet values
		preIncludeFields := []string{}
		for _, facet := range args.Facets {
			if slices.Contains([]string{"reviewer_tags"}, facet) {
				preIncludeFields = append(preIncludeFields, facet)
			}
		}
		if len(preIncludeFields) > 0 {
			facetValues, err := di.getScopedFacetValues(preIncludeFields...)
			if err != nil {
				return nil, err
			}
			for field, values := range facetValues {
				if len(values) == 0 {
					continue
				}
				query["aggs"].(M)["facets"].(M)["aggs"].(M)[field].(M)["aggs"].(M)["facet"].(M)["terms"].(M)["include"] = values
			}
		}

	}

	// ADD QUERY FILTERS
	queryFilters = append(queryFilters, di.scopes...)
	// queryFilters = append(queryFilters, internalFilters...)
	query["query"].(M)["bool"].(M)["filter"] = queryFilters

	// ADD SORTS
	sorts := []string{"date_updated:desc", "year:desc"}
	if len(args.Sort) > 0 {
		switch args.Sort[0] {
		case "date-updated-desc":
			// sorts = []string{"date_updated:desc", "year:desc"}
		case "date-updated-asc":
			sorts = []string{"date_updated:asc", "year:asc"}
		case "date-created-desc":
			sorts = []string{"date_created:desc", "year:desc"}
		case "date-created-asc":
			sorts = []string{"date_created:asc", "year:asc"}
		case "year-desc":
			sorts = []string{"year:desc"}
		case "id-asc":
			sorts = []string{"id:asc"}
		}
	}

	// SEND QUERY TO ES
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	opts := []func(*esapi.SearchRequest){
		di.client.Search.WithContext(context.Background()),
		di.client.Search.WithIndex(di.index),
		di.client.Search.WithTrackTotalHits(true),
		di.client.Search.WithSort(sorts...),
		di.client.Search.WithBody(&buf),
	}

	var res datasetResEnvelope

	err := di.searchWithOpts(opts, func(r io.ReadCloser) error {
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return fmt.Errorf("datasetindex.Search: failed to parse es6 response body: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("datasetindex.Search: %w", err)
	}

	// READ RESPONSE FROM ES
	hits, err := decodeDatasetRes(&res, args.Facets)
	if err != nil {
		return nil, fmt.Errorf("datasetindex.Search: failed to parse es6 response body: %w", err)
	}

	hits.Limit = args.Limit()
	hits.Offset = args.Offset()

	return hits, nil
}

func (di *DatasetIndex) getScopedFacetValues(fields ...string) (map[string][]string, error) {
	req := M{
		"query": M{
			"bool": M{
				"filter": di.scopes,
			},
		},
		"size": 0,
		"aggs": M{},
	}
	for _, field := range fields {
		req["aggs"].(M)[field] = M{
			"terms": M{
				"field":         field,
				"order":         M{"_key": "asc"},
				"size":          999,
				"min_doc_count": 1,
			},
		}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return nil, fmt.Errorf("datasetindex.getScopedFacetValues: %w", err)
	}

	opts := []func(*esapi.SearchRequest){
		di.client.Search.WithContext(context.Background()),
		di.client.Search.WithIndex(di.index),
		di.client.Search.WithTrackTotalHits(false),
		di.client.Search.WithBody(&buf),
	}

	var res map[string]any
	err := di.searchWithOpts(opts, func(r io.ReadCloser) error {
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return fmt.Errorf("datasetindex.getScopedFacetValues: failed to parse es6 response body: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("datasetindex.getScopedFacetValues: %w", err)
	}

	m := map[string][]string{}

	for field, agg := range res["aggregations"].(map[string]any) {
		buckets := agg.(map[string]any)["buckets"].([]any)
		m[field] = make([]string, 0, len(buckets))
		for _, bucket := range buckets {
			fv := bucket.(map[string]any)
			if v, e := fv["key_as_string"]; e {
				m[field] = append(m[field], v.(string))
			} else {
				switch v := fv["key"].(type) {
				case string:
					m[field] = append(m[field], v)
				case int:
					m[field] = append(m[field], fmt.Sprintf("%d", v))
				case float64:
					m[field] = append(m[field], fmt.Sprintf("%.2f", v))
				}
			}
		}
	}

	return m, nil
}

func (di *DatasetIndex) Each(searchArgs *models.SearchArgs, maxSize int, cb func(string)) error {
	nProcessed := 0
	start := 0
	limit := 200

	query := buildDatasetUserQuery(searchArgs)

	queryFilters := query["query"].(M)["bool"].(M)["filter"].([]M)

	// Set the searcher scopes
	queryFilters = append(queryFilters, di.scopes...)

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
			return fmt.Errorf("datasetindex.Each: failed to encode json body: %w", err)
		}

		opts := []func(*esapi.SearchRequest){
			di.client.Search.WithContext(context.Background()),
			di.client.Search.WithIndex(di.index),
			di.client.Search.WithTrackTotalHits(true),
			di.client.Search.WithSort("id:asc"),
			di.client.Search.WithBody(&buf),
		}

		var res datasetResEnvelope

		err := di.searchWithOpts(opts, func(r io.ReadCloser) error {
			if err := json.NewDecoder(r).Decode(&res); err != nil {
				return fmt.Errorf("datasetindex.Each: failed to parse es6 response body: %w", err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("datasetindex.Each: %w", err)
		}

		hits, err := decodeDatasetRes(&res, []string{})
		if err != nil {
			return fmt.Errorf("datasetindex.Each: failed to decode es6 response body: %w", err)
		}

		if len(hits.Hits) == 0 {
			return nil
		}

		for _, hit := range hits.Hits {
			nProcessed++
			if nProcessed > maxSize {
				return nil
			}
			cb(hit)
		}

		if len(hits.Hits) > 0 {
			sortValue = hits.Hits[len(hits.Hits)-1]
		}

		if len(hits.Hits) < limit {
			return nil
		}
	}
}

func (di *DatasetIndex) Delete(id string) error {
	ctx := context.Background()
	res, err := esapi.DeleteRequest{
		Index:      di.index,
		DocumentID: id,
	}.Do(ctx, di.client)
	if err != nil {
		return fmt.Errorf("datasetindex.Delete: es6 http error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return fmt.Errorf("datasetindex.Delete: io error while reading es6 error response body: %w", err)
		}
		return errors.New("datasetindex.Delete: es6 error response: " + buf.String())
	}

	return nil
}

func (di *DatasetIndex) DeleteAll() error {
	ctx := context.Background()
	req := esapi.DeleteByQueryRequest{
		Index: []string{di.index},
		Body: strings.NewReader(`{
			"query" : { 
				"match_all" : {}
			}
		}`),
	}
	res, err := req.Do(ctx, di.client)
	if err != nil {
		return fmt.Errorf("datasetindex.DeleteAll: es6 http error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return fmt.Errorf("datasetindex.DeleteAll: io error while reading es6 error response body: %w", err)
		}
		return errors.New("datasetindex.DeleteAll: es6 error response: " + buf.String())
	}

	return nil
}

func (di *DatasetIndex) WithScope(field string, terms ...string) backends.DatasetIDIndex {
	newScopes := make([]M, 0, len(di.scopes))

	// Copy existing scopes
	newScopes = append(newScopes, di.scopes...)

	// Add new scopes
	newScopes = append(newScopes, ParseScope(field, terms...))

	return &DatasetIndex{
		client: di.client,
		index:  di.index,
		scopes: newScopes,
	}
}

func (di *DatasetIndex) searchWithOpts(opts []func(*esapi.SearchRequest), fn func(r io.ReadCloser) error) error {
	res, err := di.client.Search(opts...)

	if err != nil {
		return fmt.Errorf("datasetindex.searchWithOpts: es6 http error: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return fmt.Errorf("datasetindex.searchWithOpts: io error while reading es6 error response body: %w", err)
		}
		return errors.New("datasetindex.searchWithOpts: es6 error response: " + buf.String())
	}

	return fn(res.Body)
}

func buildDatasetUserQuery(args *models.SearchArgs) M {
	var query M
	var queryMust M
	var queryFilters []M

	if len(args.Query) == 0 {
		queryMust = M{
			"match_all": M{},
		}
	} else {
		// use term based query
		// regular dis_max or multi_match are query based
		// and therefore will try to match full query over multiple fields
		queryMust = M{
			"simple_query_string": M{
				"query": args.Query,
				"fields": []string{
					"id^100",
					"identifier^50",
					"title^40",
					"organization_id^50",
					"contributor.phrase_ngram^0.05",
					"contributor.ngram^0.01",
					"all",
				},
				"lenient":                             true,
				"analyze_wildcard":                    false,
				"default_operator":                    "AND",
				"minimum_should_match":                "100%",
				"flags":                               "PHRASE",
				"auto_generate_synonyms_phrase_query": true,
			},
		}
	}

	/*
		query.bool.must: search with score
		query.bool.should: boost given search results with extra score
						   make sure minimum_should_match is 0
	*/
	if len(args.Query) > 0 {
		queryShould := []M{
			{
				"match_phrase": M{
					"title": M{
						"query": args.Query,
						"boost": 200,
					},
				},
			},
			{
				"match_phrase": M{
					"contributor": M{
						"query": args.Query,
						"boost": 200,
					},
				},
			},
			{
				"match_phrase": M{
					"all": M{
						"query": args.Query,
						"boost": 100,
					},
				},
			},
		}
		query = M{
			"query": M{
				"bool": M{
					"must":                 queryMust,
					"minimum_should_match": 0,
					"should":               queryShould,
				},
			},
		}
	} else {
		query = M{
			"query": M{
				"bool": M{
					"must": queryMust,
				},
			},
		}
	}

	if args.Filters != nil {
		for field, terms := range args.Filters {

			if qf := getRegularDatasetFilter(field, terms); qf != nil {
				if len(terms) == 0 {
					continue
				}
				/*
					TODO: invalid syntax is now solved by creating
					queries that cannot return any results.
					Error should be returned
				*/
				queryFilters = append(queryFilters, qf.ToQuery())
				continue
			}

			queryFilters = append(queryFilters, ParseScope(field, terms...))
		}
		query["query"].(M)["bool"].(M)["filter"] = queryFilters
	}

	query["size"] = args.Limit()
	query["from"] = args.Offset()

	return query
}

type datasetResEnvelope struct {
	// ScrollID string `json:"_scroll_id"`
	Hits struct {
		Total int
		Hits  []struct {
			ID string `json:"_id"`
			// Source    json.RawMessage `json:"_source"`
			// Highlight json.RawMessage
		}
	}
	Aggregations struct {
		Facets M
	}
}

func decodeDatasetRes(r *datasetResEnvelope, facets []string) (*models.SearchHits, error) {
	hits := models.SearchHits{}
	hits.Total = r.Hits.Total
	hits.Facets = make(map[string]models.FacetValues)

	//preallocate to ensure non zero slices
	for _, facet := range facets {
		hits.Facets[facet] = []models.Facet{}
	}

	for _, facet := range facets {
		if _, found := r.Aggregations.Facets[facet]; !found {
			continue
		}

		for _, f := range r.Aggregations.Facets[facet].(map[string]any)["facet"].(map[string]any)["buckets"].([]any) {
			fv := f.(map[string]any)
			value := ""
			if v, e := fv["key_as_string"]; e {
				value = v.(string)
			} else {
				switch v := fv["key"].(type) {
				case string:
					value = v
				case int:
					value = fmt.Sprintf("%d", v)
				case float64:
					value = fmt.Sprintf("%.2f", v)
				}
			}
			hits.Facets[facet] = append(hits.Facets[facet], models.Facet{
				Value: value,
				Count: int(fv["doc_count"].(float64)),
			})
		}
	}

	//reorder facet values, if applicable
	for _, facetName := range facets {
		hits.Facets[facetName] = reorderFacets(facetName, hits.Facets[facetName])
	}

	for _, h := range r.Hits.Hits {
		hits.Hits = append(hits.Hits, h.ID)
	}

	return &hits, nil
}
