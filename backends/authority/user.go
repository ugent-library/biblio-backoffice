package authority

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ugent-library/biblio-backoffice/models"
)

func (c *Client) GetUser(id string) (*models.User, error) {
	var record bson.M
	err := c.mongo.Database("authority").Collection("person").FindOne(
		context.Background(),
		bson.M{"_id": id, "active": 1}).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during document retrieval")
	}
	return c.recordToUser(record)
}

func (c *Client) GetUserByUsername(username string) (*models.User, error) {
	var record bson.M
	err := c.mongo.Database("authority").Collection("person").FindOne(
		context.Background(),
		bson.M{"ugent_username": username, "active": 1}).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during document retrieval")
	}
	return c.recordToUser(record)
}

func (c *Client) SuggestUsers(q string) ([]*models.Person, error) {
	limit := 25
	persons := make([]*models.Person, 0, limit)

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
				"filter": M{
					"term": M{
						"active": "true",
					},
				},
				"must": queryMust,
			},
		},
		"size": limit,
	}

	var responseBody personSearchEnvelope = personSearchEnvelope{}

	if e := c.search("biblio_person", requestBody, &responseBody); e != nil {
		return nil, e
	}

	for _, p := range responseBody.Hits.Hits {
		person := p.Source.Person
		person.ID = p.ID
		for _, d := range p.Source.Department {
			person.Affiliations = append(person.Affiliations, &models.Affiliation{OrganizationID: d.ID})
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (c *Client) recordToUser(record bson.M) (*models.User, error) {
	user := &models.User{}

	if v, ok := record["_id"]; ok {
		user.ID = v.(string)
	}
	if v, ok := record["email"]; ok {
		user.Email = strings.ToLower(v.(string))
	}
	if v, ok := record["ugent_username"]; ok {
		user.Username = v.(string)
	}
	if v, ok := record["active"]; ok {
		user.Active = v.(int32) == 1
	}
	if v, ok := record["orcid_token"]; ok {
		user.ORCIDToken = v.(string)
	}
	if v, ok := record["orcid_id"]; ok {
		user.ORCID = v.(string)
	}
	if v, ok := record["ugent_id"]; ok {
		for _, i := range v.(bson.A) {
			user.UGentID = append(user.UGentID, i.(string))
		}
	}
	if v, ok := record["roles"]; ok {
		for _, r := range v.(bson.A) {
			if r.(string) == "biblio-admin" {
				user.Role = "admin"
				break
			}
		}
	}
	if v, ok := record["ugent_department_id"]; ok {
		for _, i := range v.(bson.A) {
			user.Affiliations = append(user.Affiliations, &models.Affiliation{OrganizationID: i.(string)})
		}
	}
	if v, ok := record["preferred_first_name"]; ok {
		user.FirstName = v.(string)
	} else if v, ok := record["first_name"]; ok {
		user.FirstName = v.(string)
	}
	if v, ok := record["preferred_last_name"]; ok {
		user.LastName = v.(string)
	} else if v, ok := record["last_name"]; ok {
		user.LastName = v.(string)
	}

	// TODO: cleanup when authority database is synchronized with full_name
	if v, ok := record["full_name"]; ok {
		user.FullName = v.(string)
	}
	if user.FullName == "" {
		if user.FirstName != "" && user.LastName != "" {
			user.FullName = user.FirstName + " " + user.LastName
		} else if user.LastName != "" {
			user.FullName = user.LastName
		} else if user.FirstName != "" {
			user.FullName = user.FirstName
		}
	}

	if v, ok := record["date_created"]; ok {
		t, _ := time.Parse(time.RFC3339, v.(string))
		user.DateCreated = &t
	}
	if v, ok := record["date_updated"]; ok {
		t, _ := time.Parse(time.RFC3339, v.(string))
		user.DateUpdated = &t
	}

	return user, nil
}
