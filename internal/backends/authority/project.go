package authority

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/backends/es6"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		return nil, backends.ErrNotFound
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

func (c *Client) SuggestProjects(q string) ([]models.Completion, error) {
	limit := 20
	completions := make([]models.Completion, 0, limit)

	var query es6.M = es6.M{
		"match_all": es6.M{},
	}

	q = strings.TrimSpace(q)

	if q != "" {
		dismaxQueries := make([]es6.M, 0, len(projectFieldsBoosts))
		for field, boost := range projectFieldsBoosts {
			dismaxQuery := es6.M{
				"match": es6.M{
					field: es6.M{
						"query":    q,
						"operator": "AND",
						"boost":    boost,
					},
				},
			}
			dismaxQueries = append(dismaxQueries, dismaxQuery)
		}
		query = es6.M{
			"dis_max": es6.M{
				"queries": dismaxQueries,
			},
		}
	}

	requestBody := es6.M{
		"query": query,
		"size":  limit,
		"sort":  []string{"_score:desc"},
	}

	var responseBody searchEnvelope = searchEnvelope{}

	if e := c.es.SearchWithBody("biblio_project", requestBody, &responseBody); e != nil {
		return nil, e
	}

	for _, h := range responseBody.Hits.Hits {
		var m map[string]any = make(map[string]any)
		if e := json.Unmarshal(h.Source, &m); e != nil {
			return nil, e
		}
		c := models.Completion{}
		c.ID = h.ID
		if v, ok := m["title"]; ok {
			c.Heading = v.(string)
		}
		if k, e := m["eu_acronym"]; e {
			c.Description = fmt.Sprintf("(%s)", k)
		}
		completions = append(completions, c)
	}

	return completions, nil
}
