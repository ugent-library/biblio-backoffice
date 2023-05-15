package authority

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Client) GetUser(id string) (*models.User, error) {
	var record bson.M
	err := c.mongo.Database("authority").Collection("person").FindOne(
		context.Background(),
		bson.M{"_id": id, "active": 1}).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return nil, backends.ErrNotFound
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
		return nil, backends.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during document retrieval")
	}
	return c.recordToUser(record)
}

func (c *Client) SuggestUsers(q string) ([]models.Person, error) {
	limit := 25
	persons := make([]models.Person, 0, limit)

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
		person := p.Source
		person.ID = p.ID
		persons = append(persons, person)
	}

	return persons, nil
}

func (c *Client) recordToUser(record bson.M) (*models.User, error) {

	var user *models.User = &models.User{}

	if v, e := record["_id"]; e {
		user.ID = v.(string)
	}
	if v, e := record["email"]; e {
		user.Email = strings.ToLower(v.(string))
	}
	if v, e := record["ugent_username"]; e {
		user.Username = v.(string)
	}
	if v, e := record["active"]; e {
		user.Active = v.(int32) == 1
	}
	if v, e := record["orcid_token"]; e {
		user.ORCIDToken = v.(string)
	}
	if v, e := record["orcid_id"]; e {
		user.ORCID = v.(string)
	}
	if v, e := record["ugent_id"]; e {
		for _, i := range v.(bson.A) {
			user.UGentID = append(user.UGentID, i.(string))
		}
	}
	if v, e := record["roles"]; e {
		for _, r := range v.(bson.A) {
			if r.(string) == "biblio-admin" {
				user.Role = "admin"
				break
			}
		}
	}
	if v, e := record["ugent_department_id"]; e {
		for _, i := range v.(bson.A) {
			user.Department = append(user.Department, models.UserDepartment{ID: i.(string)})
		}
	}
	if v, e := record["preferred_first_name"]; e {
		user.FirstName = v.(string)
	} else if v, e := record["first_name"]; e {
		user.FirstName = v.(string)
	}
	if v, e := record["preferred_last_name"]; e {
		user.LastName = v.(string)
	} else if v, e := record["last_name"]; e {
		user.LastName = v.(string)
	}

	// TODO: cleanup when authority database is synchronized with full_name
	if v, e := record["full_name"]; e {
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

	if v, e := record["date_created"]; e {
		t, _ := time.Parse(time.RFC3339, v.(string))
		user.DateCreated = &t
	}
	if v, e := record["date_updated"]; e {
		t, _ := time.Parse(time.RFC3339, v.(string))
		user.DateUpdated = &t
	}

	return user, nil
}
