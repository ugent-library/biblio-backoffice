package authority

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"slices"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ugent-library/biblio-backoffice/models"
)

//go:embed person_settings.json
var personSettings string

func (c *Client) EnsurePersonIndexExists() error {
	res, err := esapi.IndicesExistsRequest{
		Index: []string{"biblio_person"},
	}.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		res, err := c.es.Indices.Create("biblio_person", c.es.Indices.Create.WithBody(strings.NewReader(personSettings)))
		if err != nil {
			return err
		}
		if res.IsError() {
			return fmt.Errorf("%+v", res)
		}
	}
	return nil
}

func (c *Client) SeedPerson(data []byte) error {
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	id := doc["ids"].([]any)[0].(string)
	active := doc["active"].(bool)
	doc["_id"] = id
	if active {
		doc["active"] = 1
	} else {
		doc["active"] = 0
	}
	if _, err := c.mongo.Database("authority").Collection("person").InsertOne(context.Background(), doc); err != nil {
		return err
	}
	res, err := c.es.Index(
		"biblio_person",
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

func (c *Client) CountPeople() (int64, error) {
	return c.mongo.Database("authority").Collection("person").CountDocuments(context.Background(), bson.D{}, options.Count().SetHint("_id_"))
}

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
	person := &models.Person{}

	if v, ok := record["_id"]; ok {
		person.ID = v.(string)
	}
	if v, ok := record["email"]; ok {
		person.Email = strings.ToLower(v.(string))
	}
	if v, ok := record["ugent_username"]; ok {
		person.Username = v.(string)
	}
	if v, ok := record["active"]; ok {
		person.Active = v.(int32) == 1
	}
	if v, ok := record["orcid_token"]; ok {
		person.ORCIDToken = v.(string)
	}
	if v, ok := record["orcid_id"]; ok {
		person.ORCID = v.(string)
	}
	if v, ok := record["ugent_id"]; ok {
		for _, i := range v.(bson.A) {
			person.UGentID = append(person.UGentID, i.(string))
		}
	}
	if v, ok := record["roles"]; ok {
		for _, r := range v.(bson.A) {
			if r.(string) == "biblio-admin" {
				person.Role = "admin"
				break
			}
		}
	}
	if v, ok := record["ugent_department_id"]; ok {
		for _, i := range v.(bson.A) {
			person.Affiliations = append(person.Affiliations, &models.Affiliation{OrganizationID: i.(string)})
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
	if v, ok := record["preferred_first_name"]; ok {
		person.FirstName = v.(string)
	} else if v, ok := record["first_name"]; ok {
		person.FirstName = v.(string)
	}
	if v, ok := record["preferred_last_name"]; ok {
		person.LastName = v.(string)
	} else if v, ok := record["last_name"]; ok {
		person.LastName = v.(string)
	}

	// TODO: cleanup when authority database is synchronized with full_name
	if v, ok := record["full_name"]; ok {
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

	if v, ok := record["date_created"]; ok {
		t, _ := time.Parse(time.RFC3339, v.(string))
		person.DateCreated = &t
	}
	if v, ok := record["date_updated"]; ok {
		t, _ := time.Parse(time.RFC3339, v.(string))
		person.DateUpdated = &t
	}

	return person, nil
}
