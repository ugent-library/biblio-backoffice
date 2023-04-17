package authority

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/backends/es6"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Client) GetPersons(ids []string) ([]*models.Person, error) {
	cursor, err := c.mongo.
		Database("authority").
		Collection("person").
		Find(context.Background(), bson.M{
			"ids": bson.M{
				"$in": ids,
			},
		})

	if err != nil {
		return nil, err
	}

	var records []bson.M = make([]bson.M, 0)
	if err := cursor.All(context.Background(), records); err != nil {
		if err == mongo.ErrNoDocuments {
			return []*models.Person{}, nil
		}
		return nil, err
	}

	var persons []*models.Person = make([]*models.Person, 0, len(records))
	for _, record := range records {
		person, personErr := c.recordToPerson(record)
		if personErr != nil {
			return nil, personErr
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (c *Client) GetPerson(id string) (*models.Person, error) {
	var record bson.M
	err := c.mongo.Database("authority").Collection("person").FindOne(
		context.Background(),
		bson.M{"ids": id}).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return nil, backends.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during document retrieval")
	}
	return c.recordToPerson(record)
}

func (c *Client) SuggestPeople(q string) ([]models.Person, error) {
	limit := 500
	persons := make([]models.Person, 0, limit)

	// remove duplicate spaces
	q = regexMultipleSpaces.ReplaceAllString(q, " ")

	// trim
	q = strings.TrimSpace(q)

	qParts := strings.Split(q, " ")
	queryMust := make([]es6.M, 0, len(qParts))

	for _, qp := range qParts {

		// remove terms that contain brackets
		if regexNoBrackets.MatchString(qp) {
			continue
		}

		queryMust = append(queryMust, es6.M{
			"query_string": es6.M{
				"default_operator": "AND",
				"query":            fmt.Sprintf("%s*", qp),
				"default_field":    "all",
				"analyze_wildcard": "true",
			},
		})
	}

	requestBody := es6.M{
		"query": es6.M{
			"bool": es6.M{
				"must": queryMust,
			},
		},
		"size": limit,
	}

	var responseBody personSearchEnvelope = personSearchEnvelope{}

	if e := c.es.SearchWithBody("biblio_person", requestBody, &responseBody); e != nil {
		return nil, e
	}

	for _, p := range responseBody.Hits.Hits {
		person := p.Source
		person.ID = p.ID
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
			person.Department = append(person.Department, models.PersonDepartment{ID: util.ParseString(i)})
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

	if person.FirstName != "" && person.LastName != "" {
		person.FullName = person.FirstName + " " + person.LastName
	} else if person.LastName != "" {
		person.FullName = person.LastName
	} else if person.FirstName != "" {
		person.FullName = person.FirstName
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
