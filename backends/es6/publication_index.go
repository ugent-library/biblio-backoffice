package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
)

type PublicationIndex struct {
	client *elasticsearch.Client
	index  string
	scopes []M
}

func newPublicationIndex(c *elasticsearch.Client, i string) backends.PublicationIDIndex {
	return &PublicationIndex{
		client: c,
		index:  i,
	}
}

func (pi *PublicationIndex) Search(args *models.SearchArgs) (*models.SearchHits, error) {
	// BUILD QUERY AND FILTERS FROM USER INPUT
	query := buildPublicationUserQuery(args)

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

	// ADD FACETS
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

			filters := make([]M, 0, len(pi.scopes)+1)

			// add all internal filters
			filters = append(filters, queryMust)
			filters = append(filters, pi.scopes...)
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

			var conf M
			if c, ok := facetDefinitions[field]; ok {
				conf = c.config
			} else {
				conf = defaultFacetDefinition(field).config
			}

			facet := M{
				"filter": M{"bool": M{"must": filters}},
				"aggs": M{
					"facet": conf,
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
			if _, ok := facetDefinitions[facet]; ok {
				preIncludeFields = append(preIncludeFields, facet)
			}
		}
		if len(preIncludeFields) > 0 {
			facetValues, err := pi.getScopedFacetValues(preIncludeFields...)
			if err != nil {
				return nil, err
			}
			for field, values := range facetValues {
				query["aggs"].(M)["facets"].(M)["aggs"].(M)[field].(M)["aggs"].(M)["facet"].(M)["terms"].(M)["include"] = values
			}
		}
	}

	// ADD QUERY FILTERS
	queryFilters = append(queryFilters, pi.scopes...)
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
		return nil, fmt.Errorf("publicationindex.Search: %w", err)
	}

	opts := []func(*esapi.SearchRequest){
		pi.client.Search.WithContext(context.Background()),
		pi.client.Search.WithIndex(pi.index),
		pi.client.Search.WithTrackTotalHits(true),
		pi.client.Search.WithSort(sorts...),
		pi.client.Search.WithBody(&buf),
	}

	var res publicationResEnvelope

	err := pi.searchWithOpts(opts, func(r io.ReadCloser) error {
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return fmt.Errorf("publicationindex.Search: failed to parse es6 response body: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("publicationindex.Search: %w", err)
	}

	// READ RESPONSE FROM ES
	hits, err := decodePublicationRes(&res, args.Facets)
	if err != nil {
		return nil, fmt.Errorf("publicationindex.Search: %w", err)
	}

	hits.Limit = args.Limit()
	hits.Offset = args.Offset()

	return hits, nil
}

func (pi *PublicationIndex) Each(searchArgs *models.SearchArgs, maxSize int, cb func(string)) error {
	nProcessed := 0
	start := 0
	limit := 200

	query := buildPublicationUserQuery(searchArgs)

	queryFilters := query["query"].(M)["bool"].(M)["filter"].([]M)

	// Set the searcher scopes
	queryFilters = append(queryFilters, pi.scopes...)

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
			return fmt.Errorf("publicationindex.Each: %w", err)
		}

		opts := []func(*esapi.SearchRequest){
			pi.client.Search.WithContext(context.Background()),
			pi.client.Search.WithIndex(pi.index),
			pi.client.Search.WithTrackTotalHits(true),
			pi.client.Search.WithSort("id:asc"),
			pi.client.Search.WithBody(&buf),
		}

		var res publicationResEnvelope

		err := pi.searchWithOpts(opts, func(r io.ReadCloser) error {
			if err := json.NewDecoder(r).Decode(&res); err != nil {
				return fmt.Errorf("publicationindex.Each: failed to parse es6 response body: %w", err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("publicationindex.Each: %w", err)
		}

		hits, err := decodePublicationRes(&res, []string{})
		if err != nil {
			return fmt.Errorf("publicationindex.Each: %w", err)
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

func (pi *PublicationIndex) Delete(id string) error {
	ctx := context.Background()
	res, err := esapi.DeleteRequest{
		Index:      pi.index,
		DocumentID: id,
	}.Do(ctx, pi.client)
	if err != nil {
		return fmt.Errorf("publicationindex.Delete: es6 http error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return fmt.Errorf("publicationindex.Delete: io error while reading es6 error response body: %w", err)
		}
		return errors.New("publicationindex.Delete: es6 error response: " + buf.String())
	}

	return nil
}

func (pi *PublicationIndex) DeleteAll() error {
	ctx := context.Background()
	req := esapi.DeleteByQueryRequest{
		Index: []string{pi.index},
		Body: strings.NewReader(`{
			"query" : {
				"match_all" : {}
			}
		}`),
	}
	res, err := req.Do(ctx, pi.client)
	if err != nil {
		return fmt.Errorf("publicationindex.DeleteAll: es6 http error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return fmt.Errorf("publicationindex.DeleteAll: io error while reading es6 error response body: %w", err)
		}
		return errors.New("publicationindex.DeleteAll: es6 error response: " + buf.String())
	}

	return nil
}

func (pi *PublicationIndex) WithScope(field string, terms ...string) backends.PublicationIDIndex {
	newScopes := make([]M, 0, len(pi.scopes))

	// Copy existing scopes
	newScopes = append(newScopes, pi.scopes...)

	// Add new scopes
	newScopes = append(newScopes, ParseScope(field, terms...))

	return &PublicationIndex{
		client: pi.client,
		index:  pi.index,
		scopes: newScopes,
	}
}

func (pi *PublicationIndex) searchWithOpts(opts []func(*esapi.SearchRequest), fn func(r io.ReadCloser) error) error {
	res, err := pi.client.Search(opts...)

	if err != nil {
		return fmt.Errorf("publicationindex.searchWithOpts: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return fmt.Errorf("publicationindex.searchWithOpts: io error while reading es6 error response body: %w", err)
		}
		return errors.New("publicationindex.searchWithOpts: es6 error response: " + buf.String())
	}

	return fn(res.Body)
}

func (pi *PublicationIndex) getScopedFacetValues(fields ...string) (map[string][]string, error) {
	req := M{
		"query": M{
			"bool": M{
				"filter": pi.scopes,
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
		return nil, fmt.Errorf("publicationindex.getScopedFacetValues: %w", err)
	}

	opts := []func(*esapi.SearchRequest){
		pi.client.Search.WithContext(context.Background()),
		pi.client.Search.WithIndex(pi.index),
		pi.client.Search.WithTrackTotalHits(false),
		pi.client.Search.WithBody(&buf),
	}

	var res map[string]any
	err := pi.searchWithOpts(opts, func(r io.ReadCloser) error {
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return fmt.Errorf("publicationindex.getScopedFacetValues: failed to parse es6 response body: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("publicationindex.getScopedFacetValues: %w", err)
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

func buildPublicationUserQuery(args *models.SearchArgs) M {
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
					"title^40",
					"contributor.phrase_ngram^0.05",
					"all",
				},
				"lenient":                             true,
				"analyze_wildcard":                    false,
				"default_operator":                    "AND",
				"minimum_should_match":                "100%",
				"flags":                               "PHRASE|WHITESPACE",
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
					"minimum_should_match": "0",
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

	// query.bool.filter: search without score
	if args.Filters != nil {
		for field, terms := range args.Filters {

			if qf := getRegularPublicationFilter(field, terms); qf != nil {
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

type publicationResEnvelope struct {
	// ScrollID string `json:"_scroll_id"`
	Hits struct {
		Total int
		Hits  []struct {
			ID string `json:"_id"`
			// Source json.RawMessage `json:"_source"`
			// Highlight json.RawMessage
		}
	}
	Aggregations struct {
		Facets M
	}
}

func decodePublicationRes(r *publicationResEnvelope, facets []string) (*models.SearchHits, error) {
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

			//boolean returned 0 and 1, so not to be distinguished from integers
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
