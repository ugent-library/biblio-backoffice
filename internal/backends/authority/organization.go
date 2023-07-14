package authority

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
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
		return nil, backends.ErrNotFound
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
