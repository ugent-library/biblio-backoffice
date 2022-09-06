package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Datasets struct {
	Client
	scopes        []M
	includeFacets bool
}

func NewDatasets(c Client) *Datasets {
	return &Datasets{Client: c}
}

func (datasets *Datasets) Search(args *models.SearchArgs) (*models.DatasetHits, error) {
	// BUILD QUERY AND FILTERS FROM USER INPUT
	query := datasets.buildUserQuery(args)

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

	// FACETS
	// 	create global bucket so that not all buckets are influenced by query and filters
	// 	name "facets" is not important
	if datasets.includeFacets {
		query["aggs"] = M{
			"facets": M{
				"global": M{},
				"aggs":   M{},
			},
		}

		// facet filter contains all query and all filters except itself
		for _, field := range datasetFacetFields {

			filters := make([]M, 0, len(datasets.scopes)+1)

			// add all internal filters
			filters = append(filters, queryMust)
			filters = append(filters, datasets.scopes...)
			// filters = append(filters, internalFilters...)

			// add external filter only if not matching
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
	queryFilters = append(queryFilters, datasets.scopes...)
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
		datasets.Client.es.Search.WithContext(context.Background()),
		datasets.Client.es.Search.WithIndex(datasets.Client.Index),
		datasets.Client.es.Search.WithTrackTotalHits(true),
		datasets.Client.es.Search.WithSort(sorts...),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	opts = append(opts, datasets.Client.es.Search.WithBody(&buf))

	res, err := datasets.Client.es.Search(opts...)
	if err != nil {
		return nil, err
	}

	// READ RESPONSE FROM ES
	hits, err := decodeDatasetRes(res)
	if err != nil {
		return nil, err
	}

	hits.Limit = args.Limit()
	hits.Offset = args.Offset()

	return hits, nil
}

func (datasets *Datasets) buildUserQuery(args *models.SearchArgs) M {
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
			queryFilters = append(queryFilters, ParseScope(field, terms...))
		}
		query["query"].(M)["bool"].(M)["filter"] = queryFilters
	}

	query["size"] = 20
	query["from"] = (args.Page - 1) * 20

	return query
}

type datasetResEnvelope struct {
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

/*
type resFacet struct {
	DocCount int
	Key      string
}*/

func decodeDatasetRes(res *esapi.Response) (*models.DatasetHits, error) {
	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return nil, err
		}
		return nil, errors.New("Es6 error response: " + buf.String())
	}

	var r datasetResEnvelope
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error parsing the response body")
	}

	hits := models.DatasetHits{}
	hits.Total = r.Hits.Total

	hits.Facets = make(map[string][]models.Facet)
	for _, facet := range datasetFacetFields {
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
	for facetName, facets := range hits.Facets {
		hits.Facets[facetName] = reorderFacets(facetName, facets)
	}

	for _, h := range r.Hits.Hits {
		var hit models.Dataset

		if err := json.Unmarshal(h.Source, &hit); err != nil {
			return nil, err
		}

		hits.Hits = append(hits.Hits, &hit)
	}

	return &hits, nil
}

func (publications *Datasets) Index(d *models.Dataset) error {
	body := M{
		// not needed anymore in es7 with date nano type
		"doc": struct {
			*models.Dataset
			DateCreated string `json:"date_created"`
			DateUpdated string `json:"date_updated"`
		}{
			Dataset:     d,
			DateCreated: d.DateCreated.Format(time.RFC3339),
			DateUpdated: d.DateUpdated.Format(time.RFC3339),
		},
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
		DocumentID: d.ID,
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

func (datasets *Datasets) IndexMultiple(inCh <-chan *models.Dataset) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  datasets.Client.Index,
		Client: datasets.Client.es,
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
		// not needed anymore in es7 with date nano type
		doc := struct {
			*models.Dataset
			DateCreated string `json:"date_created"`
			DateUpdated string `json:"date_updated"`
		}{
			Dataset:     p,
			DateCreated: p.DateCreated.Format(time.RFC3339),
			DateUpdated: p.DateUpdated.Format(time.RFC3339),
		}

		payload, err := json.Marshal(&doc)
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

func (datasets *Datasets) WithScope(field string, terms ...string) backends.DatasetSearchService {
	d := datasets.Clone()
	d.scopes = append(d.scopes, ParseScope(field, terms...))
	return d
}

func (datasets *Datasets) IncludeFacets(includeFacets bool) backends.DatasetSearchService {
	d := datasets.Clone()
	d.includeFacets = includeFacets
	return d
}

func (datasets *Datasets) Clone() *Datasets {
	newScopes := make([]M, 0, len(datasets.scopes))
	newScopes = append(newScopes, datasets.scopes...)
	return &Datasets{
		Client:        datasets.Client,
		scopes:        newScopes,
		includeFacets: datasets.includeFacets,
	}
}
