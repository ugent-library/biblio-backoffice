package jcr

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/ugent-library/biblio-backoffice/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type jcrRaw struct {
	Year              *int                   `bson:"year"`
	Eigenfactor       interface{}            `bson:"eigenfactor,omitempty"`
	ImmediacyIndex    interface{}            `bson:"immediacy_index,omitempty"` // float or "N/A"
	ImpactFactor      interface{}            `bson:"impact_factor,omitempty"`
	ImpactFactor5Yr   interface{}            `bson:"impact_factor_5yr,omitempty"`
	TotalCites        any                    `bson:"total_cites,omitempty"`
	CategoryRank      map[string]string      `bson:"category_rank,omitempty"`
	CategoryQuartile  map[string]interface{} `bson:"category_quartile,omitempty"`
	CategoryDecile    map[string]interface{} `bson:"category_decile,omitempty"`
	CategoryVigintile map[string]interface{} `bson:"category_vigintile,omitempty"`
}

type rank struct {
	fracture  string
	quotient  float64
	quartile  *int
	decile    *int
	vigintile *int
}

type jcr struct {
	year            *int
	eigenfactor     *float64
	immediacyIndex  *float64
	impactFactor    *float64
	impactFactor5Yr *float64
	totalCites      *int
	categoryRank    map[string]rank
}

func NewPublicationFixer(c *mongo.Client) func(context.Context, *models.Publication) error {
	return func(ctx context.Context, pub *models.Publication) error {
		for key := range pub.ExternalFields {
			if strings.HasPrefix(key, "jcr-") {
				delete(pub.ExternalFields, key)
			}
		}

		if !((len(pub.ISSN) > 0 || len(pub.EISSN) > 0) && pub.Year != "") {
			return nil
		}

		if pub.Status != "public" {
			return nil
		}

		isxn := bson.A{}
		for _, i := range pub.ISSN {
			isxn = append(isxn, i)
		}
		for _, i := range pub.EISSN {
			isxn = append(isxn, i)
		}

		var pubYear int
		if v, err := strconv.ParseInt(pub.Year, 10, 32); err != nil {
			return err
		} else {
			pubYear = int(v)
		}

		cursor, err := c.Database("authority").
			Collection("jcr").
			Find(
				ctx,
				bson.M{
					"$or": bson.A{
						bson.M{
							"issn": bson.M{"$in": isxn},
						},
						bson.M{
							"eissn": bson.M{"$in": isxn},
						},
					},
					"year": bson.M{
						"$in": bson.A{pubYear, pubYear - 1},
					},
				},
			)

		if err != nil {
			return err
		}

		var records []jcr

		for cursor.Next(ctx) {
			rawRec := jcrRaw{}
			if err := cursor.Decode(&rawRec); err != nil {
				return err
			}

			rec := jcr{
				year: rawRec.Year,
			}

			if len(rawRec.CategoryRank) > 0 {
				rec.categoryRank = make(map[string]rank)
			}

			for category, val := range rawRec.CategoryRank {
				numbers := lo.Map(strings.Split(val, "/"), func(val string, idx int) int {
					i, _ := strconv.ParseInt(val, 10, 32)
					return int(i)
				})

				categoryRank := rank{
					fracture: val,
					quotient: float64(numbers[0]) / float64(numbers[1]),
				}

				if quartile, ok := parseInt(rawRec.CategoryQuartile[category]); ok {
					categoryRank.quartile = &quartile
				}
				if decile, ok := parseInt(rawRec.CategoryDecile[category]); ok {
					categoryRank.decile = &decile
				}
				if vigintile, ok := parseInt(rawRec.CategoryVigintile[category]); ok {
					categoryRank.vigintile = &vigintile
				}

				rec.categoryRank[category] = categoryRank
			}

			// fixes "N/A" and strings
			if v, ok := parseInt(rawRec.TotalCites); ok {
				rec.totalCites = &v
			}
			if v, ok := parseFloat(rawRec.Eigenfactor); ok {
				rec.eigenfactor = &v
			}
			if v, ok := parseFloat(rawRec.ImmediacyIndex); ok {
				rec.immediacyIndex = &v
			}
			if v, ok := parseFloat(rawRec.ImpactFactor); ok {
				rec.impactFactor = &v
			}
			if v, ok := parseFloat(rawRec.ImpactFactor5Yr); ok {
				rec.impactFactor5Yr = &v
			}

			records = append(records, rec)
		}

		if err := cursor.Err(); err != nil {
			return err
		}

		if len(records) == 0 {
			return nil
		}

		newFields := models.Values{}

		var recordsThisYear []jcr
		var recordsPrevYear []jcr

		for _, record := range records {
			if *record.year == pubYear {
				recordsThisYear = append(recordsThisYear, record)
			} else if *record.year == pubYear-1 {
				recordsPrevYear = append(recordsPrevYear, record)
			}
		}

		if len(recordsThisYear) > 0 {
			eigenfactor := lo.Map(lo.Filter(recordsThisYear, func(rec jcr, idx int) bool {
				return rec.eigenfactor != nil
			}), func(rec jcr, idx int) float64 {
				return *rec.eigenfactor
			})
			slices.Sort(eigenfactor)
			if len(eigenfactor) > 0 {
				newFields.Set("jcr-eigenfactor", ffloat(eigenfactor[len(eigenfactor)-1]))
			}

			immediacy_index := lo.Map(lo.Filter(recordsThisYear, func(rec jcr, idx int) bool {
				return rec.immediacyIndex != nil
			}), func(rec jcr, idx int) float64 {
				return *rec.immediacyIndex
			})
			slices.Sort(immediacy_index)
			if len(immediacy_index) > 0 {
				newFields.Set("jcr-immediacy_index", ffloat(immediacy_index[len(immediacy_index)-1]))
			}
			impact_factor := lo.Map(lo.Filter(recordsThisYear, func(rec jcr, idx int) bool {
				return rec.impactFactor != nil
			}), func(rec jcr, idx int) float64 {
				return *rec.impactFactor
			})
			slices.Sort(impact_factor)
			if len(impact_factor) > 0 {
				newFields.Set("jcr-impact_factor", ffloat(impact_factor[len(impact_factor)-1]))
			}

			impact_factor_5yr := lo.Map(lo.Filter(recordsThisYear, func(rec jcr, idx int) bool {
				return rec.impactFactor5Yr != nil
			}), func(rec jcr, idx int) float64 {
				return *rec.impactFactor5Yr
			})
			slices.Sort(impact_factor_5yr)
			if len(impact_factor_5yr) > 0 {
				newFields.Set("jcr-impact_factor_5yr", ffloat(impact_factor_5yr[len(impact_factor_5yr)-1]))
			}

			total_cites := lo.Map(lo.Filter(recordsThisYear, func(rec jcr, idx int) bool {
				return rec.totalCites != nil
			}), func(rec jcr, idx int) int {
				return *rec.totalCites
			})
			slices.Sort(total_cites)
			if len(total_cites) > 0 {
				newFields.Set("jcr-total_cites", fmt.Sprintf("%d", total_cites[len(total_cites)-1]))
			}

			bestCategory, bestCategoryRank := bestJCRCategoryRank(recordsThisYear)
			if bestCategory != nil {
				newFields.Set("jcr-category", *bestCategory)
				newFields.Set("jcr-category_rank", bestCategoryRank.fracture)
				if bestCategoryRank.quartile != nil {
					newFields.Set("jcr-category_quartile", fmt.Sprintf("%d", *bestCategoryRank.quartile))
				}
				if bestCategoryRank.decile != nil {
					newFields.Set("jcr-category_decile", fmt.Sprintf("%d", *bestCategoryRank.decile))
				}
				if bestCategoryRank.vigintile != nil {
					newFields.Set("jcr-category_vigintile", fmt.Sprintf("%d", *bestCategoryRank.vigintile))
				}
			}
		}

		if len(recordsPrevYear) > 0 {
			prev_impact_factor := lo.Map(lo.Filter(recordsPrevYear, func(rec jcr, idx int) bool {
				return rec.impactFactor != nil
			}), func(rec jcr, idx int) float64 {
				return *rec.impactFactor
			})
			slices.Sort(prev_impact_factor)
			if len(prev_impact_factor) > 0 {
				newFields.Set("jcr-prev_impact_factor", ffloat(prev_impact_factor[len(prev_impact_factor)-1]))
			}

			bestCategory, bestCategoryRank := bestJCRCategoryRank(recordsPrevYear)
			if bestCategory != nil {
				if bestCategoryRank.quartile != nil {
					newFields.Set("jcr-prev_category_quartile", fmt.Sprintf("%d", *bestCategoryRank.quartile))
				}
				if bestCategoryRank.decile != nil {
					newFields.Set("jcr-prev_category_decile", fmt.Sprintf("%d", *bestCategoryRank.decile))
				}
				if bestCategoryRank.vigintile != nil {
					newFields.Set("jcr-prev_category_vigintile", fmt.Sprintf("%d", *bestCategoryRank.vigintile))
				}
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

func bestJCRCategoryRank(jcrRecords []jcr) (*string, *rank) {
	var tmp map[string]rank = make(map[string]rank)
	var bestRank *rank
	var bestCategory *string

	for _, jcrRecord := range jcrRecords {
		for category, currRank := range jcrRecord.categoryRank {
			if tmpRank, ok := tmp[category]; ok && tmpRank.quotient < currRank.quotient {
				continue
			}
			tmp[category] = currRank
		}
	}

	for category, currRank := range tmp {
		if bestRank == nil || bestRank.quotient > currRank.quotient {
			copyCategory := category
			bestCategory = &copyCategory
			copyRank := currRank
			bestRank = &copyRank
		}
	}

	return bestCategory, bestRank
}

func parseInt(i interface{}) (int, bool) {
	switch v := i.(type) {
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case int:
		return v, true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		p, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, false
		}
		return int(p), true
	default:
		return 0, false
	}
}

func ffloat(v float64) string {
	// shortest floating point representation (no trailing zeros)
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func parseFloat(i interface{}) (float64, bool) {
	switch v := i.(type) {
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case string:
		p, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return p, true
	default:
		return 0, false
	}
}
