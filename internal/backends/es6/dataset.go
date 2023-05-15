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

type Datasets struct {
	Client
	scopes []M
}

func NewDatasets(c Client) *Datasets {
	return &Datasets{Client: c}
}

func (datasets *Datasets) Search(args *models.SearchArgs) (*models.DatasetHits, error) {
	// BUILD QUERY AND FILTERS FROM USER INPUT
	query := buildDatasetUserQuery(args)

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
	if args.Facets != nil {
		query["aggs"] = M{
			"facets": M{
				"global": M{},
				"aggs":   M{},
			},
		}

		// facet filter contains all query and all filters except itself
		for _, field := range args.Facets {

			filters := make([]M, 0, len(datasets.scopes)+1)

			// add all internal filters
			filters = append(filters, queryMust)
			filters = append(filters, datasets.scopes...)
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

	var res datasetResEnvelope = datasetResEnvelope{}
	err := datasets.Client.SearchWithOpts(opts, &res)
	if err != nil {
		return nil, err
	}

	// READ RESPONSE FROM ES
	hits, err := decodeDatasetRes(&res, args.Facets)
	if err != nil {
		return nil, err
	}

	hits.Limit = args.Limit()
	hits.Offset = args.Offset()

	return hits, nil
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
					"identifier_values^50",
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

			if qf := getRegularDatasetFilter(field, terms...); qf != nil {
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

func decodeDatasetRes(r *datasetResEnvelope, facets []string) (*models.DatasetHits, error) {

	hits := models.DatasetHits{}
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
		var hit models.Dataset

		if err := json.Unmarshal(h.Source, &hit); err != nil {
			return nil, err
		}

		hits.Hits = append(hits.Hits, &hit)
	}

	return &hits, nil
}

func (publications *Datasets) Index(d *models.Dataset) error {
	payload, err := json.Marshal(NewIndexedDataset(d))
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := esapi.IndexRequest{
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

func (datasets *Datasets) Delete(id string) error {
	ctx := context.Background()
	res, err := esapi.DeleteRequest{
		Index:      datasets.Client.Index,
		DocumentID: id,
	}.Do(ctx, datasets.Client.es)
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

func (datasets *Datasets) DeleteAll() error {
	ctx := context.Background()
	req := esapi.DeleteByQueryRequest{
		Index: []string{datasets.Client.Index},
		Body: strings.NewReader(`{
			"query" : { 
				"match_all" : {}
			}
		}`),
	}
	res, err := req.Do(ctx, datasets.Client.es)
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

func (datasets *Datasets) WithScope(field string, terms ...string) backends.DatasetSearchService {
	d := datasets.Clone()
	d.scopes = append(d.scopes, ParseScope(field, terms...))
	return d
}

func (datasets *Datasets) Clone() *Datasets {
	newScopes := make([]M, 0, len(datasets.scopes))
	newScopes = append(newScopes, datasets.scopes...)
	return &Datasets{
		Client: datasets.Client,
		scopes: newScopes,
	}
}

func (dataset *Datasets) NewBulkIndexer(config backends.BulkIndexerConfig) (backends.BulkIndexer[*models.Dataset], error) {
	docFn := func(d *models.Dataset) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedDataset(d))
		return d.ID, doc, err
	}
	return newBulkIndexer(dataset.Client.es, dataset.Client.Index, docFn, config)
}

func (datasets *Datasets) NewIndexSwitcher(config backends.BulkIndexerConfig) (backends.IndexSwitcher[*models.Dataset], error) {
	docFn := func(d *models.Dataset) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedDataset(d))
		return d.ID, doc, err
	}
	return newIndexSwitcher(datasets.Client.es, datasets.Client.Index, datasets.Client.Settings, datasets.Client.IndexRetention, docFn, config)
}
