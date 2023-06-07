package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type PublicationIndex struct {
	Client
	scopes []M
}

func newPublicationIndex(c Client) *PublicationIndex {
	return &PublicationIndex{Client: c}
}

func (pi *PublicationIndex) Search(args *models.SearchArgs) (*models.PublicationHits, error) {
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
		case "date-created-asc":
			sorts = []string{"date_created:asc", "year:asc"}
		case "date-created-desc":
			sorts = []string{"date_created:desc", "year:desc"}
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
		pi.Client.es.Search.WithContext(context.Background()),
		pi.Client.es.Search.WithIndex(pi.Client.Index),
		pi.Client.es.Search.WithTrackTotalHits(true),
		pi.Client.es.Search.WithSort(sorts...),
		pi.Client.es.Search.WithBody(&buf),
	}

	var res publicationResEnvelope

	err := pi.searchWithOpts(opts, func(r io.ReadCloser) error {
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return fmt.Errorf("error parsing the response body")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// READ RESPONSE FROM ES
	hits, err := decodePublicationRes(&res, args.Facets)
	if err != nil {
		return nil, err
	}

	hits.Limit = args.Limit()
	hits.Offset = args.Offset()

	return hits, nil
}

func (pi *PublicationIndex) Each(searchArgs *models.SearchArgs, maxSize int, cb func(*models.Publication)) error {
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
			return err
		}

		opts := []func(*esapi.SearchRequest){
			pi.Client.es.Search.WithContext(context.Background()),
			pi.Client.es.Search.WithIndex(pi.Client.Index),
			pi.Client.es.Search.WithTrackTotalHits(true),
			pi.Client.es.Search.WithSort("id:asc"),
			pi.Client.es.Search.WithBody(&buf),
		}

		var res publicationResEnvelope

		err := pi.searchWithOpts(opts, func(r io.ReadCloser) error {
			if err := json.NewDecoder(r).Decode(&res); err != nil {
				return fmt.Errorf("error parsing the response body")
			}

			return nil
		})

		if err != nil {
			return err
		}

		hits, err := decodePublicationRes(&res, []string{})
		if err != nil {
			return err
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
			sortValue = hits.Hits[len(hits.Hits)-1].ID
		}

		if len(hits.Hits) < limit {
			return nil
		}
	}
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
					"id^100",
					"doi^50",
					"isbn^50",
					"eisbn^50",
					"issn^50",
					"eissn^50",
					"wos_id^50",
					"title^40",
					"department.tree.id^50",
					"all",
					"author.full_name.phrase_ngram^0.05",
					"author.full_name.ngram^0.01",
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
					"author.full_name": M{
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

func (pi *PublicationIndex) Delete(id string) error {
	ctx := context.Background()
	res, err := esapi.DeleteRequest{
		Index:      pi.Client.Index,
		DocumentID: id,
	}.Do(ctx, pi.Client.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return err
		}
		return errors.New("Es6 error response: " + buf.String())
	}

	return nil
}

func (pi *PublicationIndex) DeleteAll() error {
	ctx := context.Background()
	req := esapi.DeleteByQueryRequest{
		Index: []string{pi.Client.Index},
		Body: strings.NewReader(`{
			"query" : {
				"match_all" : {}
			}
		}`),
	}
	res, err := req.Do(ctx, pi.Client.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return err
		}
		return errors.New("Es6 error response: " + buf.String())
	}

	return nil
}

func (pi *PublicationIndex) WithScope(field string, terms ...string) backends.PublicationIndex {
	newScopes := make([]M, 0, len(pi.scopes))

	// Copy existing scopes
	newScopes = append(newScopes, pi.scopes...)

	// Add new scopes
	newScopes = append(newScopes, ParseScope(field, terms...))

	return &PublicationIndex{
		Client: pi.Client,
		scopes: newScopes,
	}
}

func (pi *PublicationIndex) searchWithOpts(opts []func(*esapi.SearchRequest), fn func(r io.ReadCloser) error) error {
	res, err := pi.es.Search(opts...)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return err
		}
		return errors.New("Es6 error response: " + buf.String())
	}

	return fn(res.Body)
}

type publicationResEnvelope struct {
	// ScrollID string `json:"_scroll_id"`
	Hits struct {
		Total int
		Hits  []struct {
			Source    json.RawMessage `json:"_source"`
			Highlight json.RawMessage
		}
	}
	Aggregations struct {
		Facets M
	}
}

func decodePublicationRes(r *publicationResEnvelope, facets []string) (*models.PublicationHits, error) {

	hits := models.PublicationHits{}
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
		var hit models.Publication

		if err := json.Unmarshal(h.Source, &hit); err != nil {
			return nil, err
		}

		hits.Hits = append(hits.Hits, &hit)
	}

	return &hits, nil
}
