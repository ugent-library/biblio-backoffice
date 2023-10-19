package authority

import (
	"context"
	"fmt"
	"strings"
	"time"

	"slices"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/util"
)

func (c *Client) GetPerson(id string) (*models.Person, error) {
	var record bson.M
	err := c.mongo.Database("authority").Collection("person").FindOne(
		context.Background(),
		bson.M{"ids": id}).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during document retrieval")
	}
	return c.recordToPerson(record)
}

func (c *Client) SuggestPeople(q string) ([]*models.Person, error) {
	limit := 20
	persons := make([]*models.Person, 0, limit)

	// remove duplicate spaces
	q = regexMultipleSpaces.ReplaceAllString(q, " ")

	// trim
	q = strings.TrimSpace(q)

	qParts := strings.Split(q, " ")
	queryMust := make([]M, 0, len(qParts))

	for _, qp := range qParts {

		// remove terms that contain brackets
		if regexNoBrackets.MatchString(qp) {
			continue
		}

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
				"should": M{
					"match_phrase_prefix": M{
						"full_name": M{
							"query": q,
							"boost": 100,
						},
					},
				},
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

func (c *Client) recordToPerson(record bson.M) (*models.Person, error) {
	var person *models.Person = &models.Person{}

	if v, e := record["_id"]; e {
		// _id might be stored as number, float or even "null"
		person.ID = util.ParseString(v)
	}
	if v, e := record["active"]; e {
		person.Active = util.ParseBoolean(v)
	}
	if v, e := record["orcid_id"]; e {
		// orcid might be stored as "null"
		person.ORCID = util.ParseString(v)
	}
	if v, e := record["ugent_id"]; e {
		for _, i := range v.(bson.A) {
			person.UGentID = append(person.UGentID, util.ParseString(i))
		}
	}
	if v, e := record["ugent_department_id"]; e {
		for _, i := range v.(bson.A) {
			person.Affiliations = append(person.Affiliations, &models.Affiliation{OrganizationID: util.ParseString(i)})
		}
	}
	if v, e := record["object_class"]; e {
		objectClass := []string{}
		for _, val := range v.(bson.A) {
			objectClass = append(objectClass, val.(string))
		}
		if slices.Contains(objectClass, "ugentFormerEmployee") && len(person.Affiliations) == 0 {
			person.Affiliations = append(person.Affiliations, &models.Affiliation{OrganizationID: "UGent"})
		}
		if slices.Contains(objectClass, "uzEmployee") {
			person.Affiliations = append(person.Affiliations, &models.Affiliation{OrganizationID: "UZGent"})
		}
	}
	if v, e := record["preferred_first_name"]; e {
		person.FirstName = util.ParseString(v)
	} else if v, e := record["first_name"]; e {
		person.FirstName = util.ParseString(v)
	}
	if v, e := record["preferred_last_name"]; e {
		person.LastName = util.ParseString(v)
	} else if v, e := record["last_name"]; e {
		person.LastName = util.ParseString(v)
	}

	// TODO: cleanup when authority database is synchronized with full_name
	if v, e := record["full_name"]; e {
		person.FullName = v.(string)
	}
	if person.FullName == "" {
		if person.FirstName != "" && person.LastName != "" {
			person.FullName = person.FirstName + " " + person.LastName
		} else if person.LastName != "" {
			person.FullName = person.LastName
		} else if person.FirstName != "" {
			person.FullName = person.FirstName
		}
	}

	if v, e := record["date_created"]; e {
		t, _ := time.Parse(time.RFC3339, util.ParseString(v))
		person.DateCreated = &t
	}
	if v, e := record["date_updated"]; e {
		t, _ := time.Parse(time.RFC3339, util.ParseString(v))
		person.DateUpdated = &t
	}

	return person, nil
}
