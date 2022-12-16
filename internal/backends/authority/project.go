package authority

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Client) GetProject(id string) (*models.Project, error) {
	var pr map[string]string = make(map[string]string)
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
		}).Decode(&pr)
	if err == mongo.ErrNoDocuments {
		return nil, backends.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during document retrieval")
	}
	return &models.Project{
		ID:        pr["_id"],
		Title:     pr["title"],
		StartDate: pr["start_date"],
		EndDate:   pr["end_date"],
	}, nil
}
