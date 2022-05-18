package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) SearchDatasets(args *models.SearchArgs) (*models.DatasetHits, error) {
	query := buildDatasetQuery(args)

	queryFilters := query["query"].(M)["bool"].(M)["filter"].([]M)
	queryMust := query["query"].(M)["bool"].(M)["must"].(M)
	query["size"] = args.Limit()
	query["from"] = args.Offset()

	// if args.Highlight {
	// 	query["highlight"] = M{
	// 		"require_field_match": false,
	// 		"pre_tags":            []string{"<mark>"},
	// 		"post_tags":           []string{"</mark>"},
	// 		"fields": M{
	// 			"metadata.title.ngram":       M{},
	// 			"metadata.author.name.ngram": M{},
	// 		},
	// 	}
	// }

	query["aggs"] = M{
		"facets": M{
			"global": M{},
			"aggs":   M{},
		},
	}

	// facet filter contains all query and all filters except itself
	for _, field := range []string{"status", "faculty"} {
		filters := []M{queryMust}

		for _, filter := range queryFilters {
			terms := filter["terms"]
			// non facet related filters (keep: always)
			// TODO: make difference between facet filter
			// and other filters more explicit
			if terms == nil {
				filters = append(filters, filter)
				continue
			}
			// facet related filters (keep: if not matching)
			if _, found := terms.(M)[field]; found {
				continue
			} else {
				filters = append(filters, filter)
			}
		}

		if field == "faculty" {

			query["aggs"].(M)["facets"].(M)["aggs"].(M)[field] = M{
				"filter": M{"bool": M{"must": filters}},
				"aggs": M{
					"facet": M{
						"terms": M{
							"field":   field,
							"order":   M{"_key": "asc"},
							"size":    200,
							"include": "^CA|DS|DI|EB|FW|GE|LA|LW|PS|PP|RE|TW|WE|GUK|UZGent|HOART|HOGENT|HOWEST|IBBT|IMEC|VIB$",
						},
					},
				},
			}

		} else {

			query["aggs"].(M)["facets"].(M)["aggs"].(M)[field] = M{
				"filter": M{"bool": M{"must": filters}},
				"aggs": M{
					"facet": M{
						"terms": M{
							"field": field,
							"order": M{"_key": "asc"},
							"size":  200,
						},
					},
				},
			}

		}
	}

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

	opts := []func(*esapi.SearchRequest){
		c.es.Search.WithContext(context.Background()),
		c.es.Search.WithIndex(c.DatasetIndex),
		c.es.Search.WithTrackTotalHits(true),
		c.es.Search.WithSort(sorts...),
	}

	// if args.Cursor {
	// 	opts = append(opts, s.client.Search.WithScroll(time.Minute))
	// }

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	opts = append(opts, c.es.Search.WithBody(&buf))

	res, err := c.es.Search(opts...)
	if err != nil {
		return nil, err
	}

	hits, err := decodeDatasetRes(res)
	if err != nil {
		return nil, err
	}

	hits.Limit = args.Limit()
	hits.Offset = args.Offset()

	return hits, nil
}

func buildDatasetQuery(args *models.SearchArgs) M {
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

	if args.Filters != nil {
		for field, terms := range args.Filters {
			orFields := strings.Split(field, "|")
			if len(orFields) > 1 {
				orFilters := []M{}
				for _, orField := range orFields {
					orFilters = append(orFilters, M{"terms": M{orField: terms}})
				}
				queryFilters = append(queryFilters, M{"bool": M{"should": orFilters, "minimum_should_match": "1"}})
			} else if strings.HasPrefix(field, "!") {
				queryFilters = append(queryFilters, M{"bool": M{"must_not": M{"terms": M{field[1:]: terms}}}})
			} else {
				queryFilters = append(queryFilters, M{"terms": M{field: terms}})
			}
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

type resFacet struct {
	DocCount int
	Key      string
}

func decodeDatasetRes(res *esapi.Response) (*models.DatasetHits, error) {
	defer res.Body.Close()

	if res.IsError() {
		buf := new(strings.Builder)
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
	for _, facet := range []string{"status"} {
		if _, found := r.Aggregations.Facets[facet]; !found {
			continue
		}

		for _, f := range r.Aggregations.Facets[facet].(map[string]interface{})["facet"].(map[string]interface{})["buckets"].([]interface{}) {
			fv := f.(map[string]interface{})
			value := ""
			switch v := fv["key"].(type) {
			case string:
				value = v
			case int:
				value = fmt.Sprintf("%d", v)
			case float64:
				value = fmt.Sprintf("%.2f", v)
			}
			hits.Facets[facet] = append(hits.Facets[facet], models.Facet{
				Value: value,
				Count: int(fv["doc_count"].(float64)),
			})
		}
	}

	for _, h := range r.Hits.Hits {
		var hit models.Dataset

		if err := json.Unmarshal(h.Source, &hit); err != nil {
			return nil, err
		}

		// if len(h.Highlight) > 0 {
		// 	hit.RawHighlight = h.Highlight
		// }

		hits.Hits = append(hits.Hits, &hit)
	}

	return &hits, nil
}

func (c *Client) IndexDataset(d *models.Dataset) error {
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
		Index:      c.DatasetIndex,
		DocumentID: d.ID,
		Body:       bytes.NewReader(payload),
	}.Do(ctx, c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		buf := new(strings.Builder)
		if _, err := io.Copy(buf, res.Body); err != nil {
			return err
		}
		return errors.New("Es6 error response: " + buf.String())
	}

	return nil
}

func (c *Client) IndexDatasets(inCh <-chan *models.Dataset) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  c.DatasetIndex,
		Client: c.es,
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
				Action:       "index",
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
