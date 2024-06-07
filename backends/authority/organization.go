package authority

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/ugent-library/biblio-backoffice/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type organizationSearchEnvelope struct {
	Hits struct {
		Total int `json:"total"`
		Hits  []struct {
			ID     string              `json:"_id"`
			Source models.Organization `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type searchEnvelope struct {
	Hits struct {
		Total int `json:"total"`
		Hits  []struct {
			ID     string          `json:"_id"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

//go:embed organization_settings.json
var organizationSettings string

func (c *Client) EnsureOrganizationSeedIndexExists() error {
	res, err := esapi.IndicesExistsRequest{
		Index: []string{"biblio_organization"},
	}.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		res, err := c.es.Indices.Create("biblio_organization", c.es.Indices.Create.WithBody(strings.NewReader(organizationSettings)))
		if err != nil {
			return err
		}
		if res.IsError() {
			return fmt.Errorf("%+v", res)
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (c *Client) SeedOrganization(data []byte) error {
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	id := doc["id"].(string)
	doc["_id"] = id
	if _, err := c.mongo.Database("authority").Collection("organization").ReplaceOne(context.Background(), bson.M{"_id": id}, doc, options.Replace().SetUpsert(true)); err != nil {
		return err
	}
	res, err := c.es.Index(
		"biblio_organization",
		bytes.NewReader(data),
		c.es.Index.WithDocumentID(id),
		c.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("%+v", res)
	}
	return nil
}

func (c *Client) CountOrganizations() (int64, error) {
	return c.mongo.Database("authority").Collection("organization").CountDocuments(context.Background(), bson.D{}, options.Count().SetHint("_id_"))
}

func (c *Client) GetOrganization(id string) (*models.Organization, error) {
	requestBody := M{
		"query": M{
			"term": M{
				"_id": id,
			},
		},
		"size": 1,
	}
	responseBody := organizationSearchEnvelope{}

	if e := c.search("biblio_organization", requestBody, &responseBody); e != nil {
		return nil, e
	}

	if len(responseBody.Hits.Hits) == 0 {
		return nil, models.ErrNotFound
	}

	org := responseBody.Hits.Hits[0].Source
	org.ID = responseBody.Hits.Hits[0].ID

	return &org, nil
}

func (c *Client) SuggestOrganizations(q string) ([]models.Completion, error) {
	limit := 20
	completions := make([]models.Completion, 0, limit)

	// remove duplicate spaces
	q = regexMultipleSpaces.ReplaceAllString(q, " ")

	// trim
	q = strings.TrimSpace(q)

	qParts := strings.Split(q, " ")
	queryMust := make([]M, 0, len(qParts))

	for _, qp := range qParts {
		queryMust = append(queryMust, M{
			"query_string": M{
				"default_operator": "AND",
				"query":            fmt.Sprintf("%s*", qp),
				"default_field":    "all",
				"analyze_wildcard": "true",
			},
		})
	}

	requestBody := M{
		"query": M{
			"bool": M{
				"must": queryMust,
			},
		},
		"size": limit,
	}

	var responseBody searchEnvelope = searchEnvelope{}

	if e := c.search("biblio_organization", requestBody, &responseBody); e != nil {
		return nil, e
	}

	for _, h := range responseBody.Hits.Hits {
		var m map[string]any = make(map[string]any)
		if e := json.Unmarshal(h.Source, &m); e != nil {
			return nil, e
		}
		c := models.Completion{}
		c.ID = h.ID
		c.Heading = m["name"].(string)
		completions = append(completions, c)
	}

	return completions, nil
}
