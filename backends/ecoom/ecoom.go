package ecoom

import (
	"context"
	"fmt"
	"strings"

	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewPublicationFixer(c *mongo.Client) func(context.Context, *models.Publication) error {
	return func(ctx context.Context, pub *models.Publication) error {
		for key := range pub.ExternalFields {
			if strings.HasPrefix(key, "ecoom-") {
				delete(pub.ExternalFields, key)
			}
		}

		if pub.WOSID == "" {
			return nil
		}

		cursor, err := c.Database("ecoom").
			Collection("publication").
			Find(
				ctx,
				bson.M{"locatie_nummer": pub.WOSID},
				options.Find().SetSort(bson.D{bson.E{Key: "year", Value: -1}}),
			)

		if err != nil {
			return err
		}

		var records []bson.M

		for cursor.Next(ctx) {
			var record bson.M
			if err := cursor.Decode(&record); err != nil {
				return err
			}
			records = append(records, record)
		}

		if err := cursor.Err(); err != nil {
			return err
		}

		if len(records) == 0 {
			return nil
		}

		newFields := models.Values{}

		for _, fund := range []string{"bof", "iof"} {
			var fundRecords []bson.M
			for _, rec := range records {
				if rec["fund"] == fund {
					fundRecords = append(fundRecords, rec)
				}
			}
			if len(fundRecords) == 0 {
				continue
			}
			if v, ok := fundRecords[0]["gewicht"]; ok {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "weight"), util.ParseString(v))
			}
			if v, ok := fundRecords[0]["css"]; ok {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "css"), util.ParseString(v))
			}
			if v, ok := fundRecords[0]["internationale_samenwerking"]; ok {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "international-collaboration"), util.ParseString(v))
			}
			if v, ok := fundRecords[0]["hoger_onderwijs"]; ok && util.ParseBoolean(v) {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "sector"), "higher-education")
			}
			if v, ok := fundRecords[0]["met_ziekenhuis"]; ok && util.ParseBoolean(v) {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "sector"), "hospital")
			}
			if v, ok := fundRecords[0]["met_publieke_instelling"]; ok && util.ParseBoolean(v) {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "sector"), "government")
			}
			if v, ok := fundRecords[0]["met_private_instelling"]; ok && util.ParseBoolean(v) {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "sector"), "private")
			}
			if v, ok := fundRecords[len(fundRecords)-1]["year"]; ok {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "validation"), util.ParseString(v))
			}
		}

		if len(newFields) > 0 {
			if pub.ExternalFields == nil {
				pub.ExternalFields = models.Values{}
			}
			for key, vals := range newFields {
				pub.ExternalFields.SetAll(key, vals...)
			}
		}

		return nil
	}
}
