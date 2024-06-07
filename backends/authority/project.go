package authority

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backoffice/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:embed project_settings.json
var projectSettings string

func (c *Client) EnsureProjectIndexExists() error {
	res, err := esapi.IndicesExistsRequest{
		Index: []string{"biblio_project"},
	}.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		res, err := c.es.Indices.Create("biblio_project", c.es.Indices.Create.WithBody(strings.NewReader(projectSettings)))
		if err != nil {
			return err
		}
		if res.IsError() {
			return fmt.Errorf("%+v", res)
		}
	}
	return nil
}

func (c *Client) SeedProject(data []byte) error {
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	id := doc["iweto_id"].(string)
	doc["_id"] = id
	if _, err := c.mongo.Database("authority").Collection("project").InsertOne(context.Background(), doc); err != nil {
		return err
	}
	res, err := c.es.Index(
		"biblio_project",
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

func (c *Client) CountProjects() (int64, error) {
	return c.mongo.Database("authority").Collection("project").CountDocuments(context.Background(), bson.D{}, options.Count().SetHint("_id_"))
}

func (c *Client) GetProject(id string) (*models.Project, error) {
	var rec bson.M
	err := c.mongo.Database("authority").Collection("project").FindOne(
		context.Background(),
		bson.M{
			"$or": bson.A{
				bson.M{"_id": id},
				bson.M{"eu_id": id},
				bson.M{
					"eu_acronym": bson.M{
						"$in": bson.A{
							id,
							strings.ToLower(id),
							strings.ToUpper(id),
						},
					},
				},
			},
		}).Decode(&rec)
	if err == mongo.ErrNoDocuments {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during document retrieval")
	}

	p := &models.Project{
		ID: rec["_id"].(string),
	}
	if v, ok := rec["title"]; ok {
		p.Title = v.(string)
	}
	if v, ok := rec["start_date"]; ok {
		p.StartDate = v.(string)
	}
	if v, ok := rec["end_date"]; ok {
		p.EndDate = v.(string)
	}
	if v, ok := rec["eu_id"]; ok {
		if p.EUProject == nil {
			p.EUProject = &models.EUProject{}
		}
		p.EUProject.ID = v.(string)
	}
	if v, ok := rec["eu_call_id"]; ok {
		if p.EUProject == nil {
			p.EUProject = &models.EUProject{}
		}
		p.EUProject.CallID = v.(string)
	}
	if v, ok := rec["eu_acronym"]; ok {
		if p.EUProject == nil {
			p.EUProject = &models.EUProject{}
		}
		p.EUProject.Acronym = v.(string)
	}
	if v, ok := rec["eu_framework_programme"]; ok {
		if p.EUProject == nil {
			p.EUProject = &models.EUProject{}
		}
		p.EUProject.FrameworkProgramme = v.(string)
	}

	if v, ok := rec["gismo_id"]; ok {
		p.GISMOID = v.(string)
	}

	if v, ok := rec["iweto_id"]; ok {
		p.IWETOID = v.(string)
	}
	// iweto_id not filled in everywhere, but should be same as id for now
	if p.IWETOID == "" {
		p.IWETOID = p.ID
	}

	return p, nil
}

var projectFieldsBoosts = map[string]string{
	// field: boost
	"_id":                    "100",
	"iweto_id":               "80",
	"gismo_id":               "80",
	"eu_call_id":             "80",
	"eu_id":                  "80",
	"eu_framework_programme": "80",
	"eu_acronym":             "70",
	"all":                    "0.1",
	"phrase_ngram":           "0.05",
	"ngram":                  "0.01",
}

func (c *Client) SuggestProjects(q string) ([]*models.Project, error) {
	limit := 20
	projects := make([]*models.Project, 0, limit)

	var query M = M{
		"match_all": M{},
	}

	q = strings.TrimSpace(q)

	if q != "" {
		dismaxQueries := make([]M, 0, len(projectFieldsBoosts))
		for field, boost := range projectFieldsBoosts {
			dismaxQuery := M{
				"match": M{
					field: M{
						"query":    q,
						"operator": "AND",
						"boost":    boost,
					},
				},
			}
			dismaxQueries = append(dismaxQueries, dismaxQuery)
		}
		query = M{
			"dis_max": M{
				"queries": dismaxQueries,
			},
		}
	}

	requestBody := M{
		"query": query,
		"size":  limit,
		"sort":  []string{"_score:desc"},
	}

	var responseBody searchEnvelope = searchEnvelope{}

	if e := c.search("biblio_project", requestBody, &responseBody); e != nil {
		return nil, e
	}

	for _, h := range responseBody.Hits.Hits {
		m := make(map[string]any)
		if e := json.Unmarshal(h.Source, &m); e != nil {
			return nil, e
		}
		p := models.Project{}
		p.ID = h.ID
		if v, ok := m["title"]; ok {
			p.Title = v.(string)
		}
		if v, ok := m["start_date"]; ok {
			p.StartDate = v.(string)
		}
		if v, ok := m["end_date"]; ok {
			p.EndDate = v.(string)
		}
		if v, ok := m["eu_id"]; ok {
			if p.EUProject == nil {
				p.EUProject = &models.EUProject{}
			}
			p.EUProject.ID = v.(string)
		}
		if v, ok := m["eu_call_id"]; ok {
			if p.EUProject == nil {
				p.EUProject = &models.EUProject{}
			}
			p.EUProject.CallID = v.(string)
		}
		if v, ok := m["eu_acronym"]; ok {
			if p.EUProject == nil {
				p.EUProject = &models.EUProject{}
			}
			p.EUProject.Acronym = v.(string)
		}
		if v, ok := m["eu_framework_programme"]; ok {
			if p.EUProject == nil {
				p.EUProject = &models.EUProject{}
			}
			p.EUProject.FrameworkProgramme = v.(string)
		}

		if v, ok := m["gismo_id"]; ok {
			p.GISMOID = v.(string)
		}

		if v, ok := m["iweto_id"]; ok {
			p.IWETOID = v.(string)
		}
		// iweto_id not filled in everywhere, but should be same as id for now
		if p.IWETOID == "" {
			p.IWETOID = p.ID
		}

		projects = append(projects, &p)
	}

	return projects, nil
}
