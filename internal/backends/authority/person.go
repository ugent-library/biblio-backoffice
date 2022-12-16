package authority

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Client) GetPerson(id string) (*models.Person, error) {
	var record bson.M
	err := c.mongo.Database("authority").Collection("person").FindOne(
		context.Background(),
		bson.M{"_id": id}).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return nil, backends.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during document retrieval")
	}
	return c.recordToPerson(record)
}

func (c *Client) recordToPerson(record bson.M) (*models.Person, error) {

	var person *models.Person = &models.Person{}

	if v, e := record["_id"]; e {
		person.ID = v.(string)
	}
	if v, e := record["active"]; e {
		person.Active = v.(int32) == 1
	}
	if v, e := record["orcid_id"]; e {
		person.ORCID = v.(string)
	}
	if v, e := record["ugent_id"]; e {
		for _, i := range v.(bson.A) {
			person.UGentID = append(person.UGentID, i.(string))
		}
	}
	if v, e := record["ugent_department_id"]; e {
		for _, i := range v.(bson.A) {
			person.Department = append(person.Department, models.PersonDepartment{ID: i.(string)})
		}
	}
	if v, e := record["preferred_first_name"]; e {
		person.FirstName = v.(string)
	} else if v, e := record["first_name"]; e {
		person.FirstName = v.(string)
	}
	if v, e := record["preferred_last_name"]; e {
		person.LastName = v.(string)
	} else if v, e := record["last_name"]; e {
		person.LastName = v.(string)
	}

	if person.FirstName != "" && person.LastName != "" {
		person.FullName = person.FirstName + " " + person.LastName
	} else if person.LastName != "" {
		person.FullName = person.LastName
	} else if person.FirstName != "" {
		person.FullName = person.FirstName
	}

	if v, e := record["date_created"]; e {
		t, _ := time.Parse(time.RFC3339, v.(string))
		person.DateCreated = &t
	}
	if v, e := record["date_updated"]; e {
		t, _ := time.Parse(time.RFC3339, v.(string))
		person.DateUpdated = &t
	}

	return person, nil
}
