package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Publications struct {
	Client
	scopes []M
}

func NewPublications(c Client) *Publications {
	return &Publications{Client: c}
}

func (publications *Publications) Search(args *models.SearchArgs) (*models.PublicationHits, error) {
	// BUILD QUERY AND FILTERS FROM USER INPUT
	query := publications.buildUserQuery(args)

	queryFilters := query["query"].(M)["bool"].(M)["filter"].([]M)
	queryMust := query["query"].(M)["bool"].(M)["must"].(M)
	query["size"] = args.Limit()
	query["from"] = args.Offset()

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

			filters := make([]M, 0, len(publications.scopes)+1)

			// add all internal filters
			filters = append(filters, queryMust)
			filters = append(filters, publications.scopes...)
			// filters = append(filters, internalFilters...)

			// add external filters only if not matching
			for _, filter := range queryFilters {
				terms := filter["terms"]
				if terms == nil {
					continue
				}
				if _, found := terms.(M)[field]; found {
					continue
				} else {
					filters = append(filters, filter)
				}
			}

			// TODO make configurable for each facet
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
	}

	// ADD QUERY FILTERS
	queryFilters = append(queryFilters, publications.scopes...)
	// queryFilters = append(queryFilters, internalFilters...)
	query["query"].(M)["bool"].(M)["filter"] = queryFilters

	// ADD SORTS
	sorts := []string{"date_updated:desc", "year:desc"}
	if len(args.Sort) > 0 {
		switch args.Sort[0] {
		case "date-updated-desc":
			// sorts = []string{"date_updated:desc", "year:desc"}
		case "date-created-desc":
			sorts = []string{"date_created:desc", "year:desc"}
		case "year-desc":
			sorts = []string{"year:desc"}
		}
	}

	// SEND QUERY TO ES
	opts := []func(*esapi.SearchRequest){
		publications.Client.es.Search.WithContext(context.Background()),
		publications.Client.es.Search.WithIndex(publications.Client.Index),
		publications.Client.es.Search.WithTrackTotalHits(true),
		publications.Client.es.Search.WithSort(sorts...),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	opts = append(opts, publications.Client.es.Search.WithBody(&buf))

	res, err := publications.Client.es.Search(opts...)
	if err != nil {
		return nil, err
	}

	// READ RESPONSE FROM ES
	hits, err := decodePublicationRes(res, args.Facets)
	if err != nil {
		return nil, err
	}

	hits.Limit = args.Limit()
	hits.Offset = args.Offset()

	return hits, nil
}

func (publications *Publications) buildUserQuery(args *models.SearchArgs) M {
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
					"title^40",
					"all",
					"author.full_name.phrase_ngram^0.05",
					"author.full_name.ngram^0.01",
				},
				"lenient":                             true,
				"analyze_wildcard":                    false,
				"default_operator":                    "OR",
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
			queryFilters = append(queryFilters, ParseScope(field, terms...))
		}
		query["query"].(M)["bool"].(M)["filter"] = queryFilters
	}

	query["size"] = 20
	query["from"] = (args.Page - 1) * 20

	return query
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

func decodePublicationRes(res *esapi.Response, facets []string) (*models.PublicationHits, error) {
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return nil, err
		}
		return nil, errors.New("Es6 error response: " + buf.String())
	}

	var r publicationResEnvelope
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error parsing the response body")
	}

	hits := models.PublicationHits{}
	hits.Total = r.Hits.Total

	hits.Facets = make(map[string][]models.Facet)
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
	for facetName, facets := range hits.Facets {
		hits.Facets[facetName] = reorderFacets(facetName, facets)
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

func (publications *Publications) Index(p *models.Publication) error {
	doc := NewIndexedPublication(p)
	body := M{
		"doc":           doc,
		"doc_as_upsert": true,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := esapi.UpdateRequest{
		Index: publications.Client.Index,
		// DocumentID: d.SnapshotID,
		DocumentID: p.ID,
		Body:       bytes.NewReader(payload),
	}.Do(ctx, publications.Client.es)
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

// TODO error chan
func (publications *Publications) IndexMultiple(inCh <-chan *models.Publication) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  publications.Client.Index,
		Client: publications.Client.es,
		OnError: func(c context.Context, e error) {
			log.Fatalf("ERROR: %s", e)
		},
		/*
			TODO: appropriate place for this?
			without this a controller may search too soon,
			and see no results
		*/
		Refresh: "wait_for",
	})
	if err != nil {
		log.Fatal(err)
	}

	for p := range inCh {
		doc := NewIndexedPublication(p)

		payload, err := json.Marshal(doc)
		if err != nil {
			log.Panic(err)
		}

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",
				// DocumentID:   doc.SnapshotID,
				DocumentID:   doc.ID,
				DocumentType: "_doc",
				Body:         bytes.NewReader(payload),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Panicf("ERROR: %s", err)
					} else {
						log.Panicf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Panicf("Unexpected error: %s", err)
		}
	}

	// Close the indexer
	if err := bi.Close(context.Background()); err != nil {
		log.Panicf("Unexpected error: %s", err)
	}
}

func (publications *Publications) WithScope(field string, terms ...string) backends.PublicationSearchService {
	p := publications.Clone()
	p.scopes = append(p.scopes, ParseScope(field, terms...))
	return p
}

func (publications *Publications) Clone() *Publications {
	newScopes := make([]M, 0, len(publications.scopes))
	newScopes = append(newScopes, publications.scopes...)
	return &Publications{
		Client: publications.Client,
		scopes: newScopes,
	}
}
