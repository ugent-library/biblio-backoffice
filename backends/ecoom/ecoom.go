package ecoom

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/ugent-library/biblio-backoffice/models"
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
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "weight"), parseString(v))
			}
			if v, ok := fundRecords[0]["css"]; ok {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "css"), parseString(v))
			}
			if v, ok := fundRecords[0]["internationale_samenwerking"]; ok {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "international-collaboration"), lo.Ternary(parseBoolean(v), "true", "false"))
			}
			if v, ok := fundRecords[0]["hoger_onderwijs"]; ok && parseBoolean(v) {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "sector"), "higher-education")
			}
			if v, ok := fundRecords[0]["met_ziekenhuis"]; ok && parseBoolean(v) {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "sector"), "hospital")
			}
			if v, ok := fundRecords[0]["met_publieke_instelling"]; ok && parseBoolean(v) {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "sector"), "government")
			}
			if v, ok := fundRecords[0]["met_private_instelling"]; ok && parseBoolean(v) {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "sector"), "private")
			}
			if v, ok := fundRecords[len(fundRecords)-1]["year"]; ok {
				newFields.Add(fmt.Sprintf("ecoom-%s-%s", fund, "validation"), parseString(v))
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

func parseBoolean(v any) bool {
	switch b := v.(type) {
	case int32:
		return b == 1
	case int64:
		return b == 1
	case string:
		return b == "true" || b == "1"
	case bool:
		return b
	}
	return false
}

func parseString(v any) string {
	switch s := v.(type) {
	case int:
		return fmt.Sprintf("%d", s)
	case int32:
		return fmt.Sprintf("%d", s)
	case int64:
		return fmt.Sprintf("%d", s)
	case float32:
		return fmt.Sprintf("%g", s)
	case float64:
		return fmt.Sprintf("%g", s)
	case string:
		return s
	case bool:
		if s {
			return "true"
		} else {
			return "false"
		}
	}
	return ""
}
