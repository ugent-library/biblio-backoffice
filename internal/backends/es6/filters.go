package es6

import (
	"regexp"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backoffice/internal/models"
	internal_time "github.com/ugent-library/biblio-backoffice/internal/time"
)

// regular field filter: accepts syntax in the filter value
type FieldFilter struct {
	models.BaseFilter
}

func (ff *FieldFilter) ToQuery() map[string]interface{} {
	return ParseScope(ff.Name, ff.Values...)
}

func (ff *FieldFilter) GetType() string {
	return "field"
}

// date filter
type DateSinceFilter struct {
	models.BaseFilter
}

func (dbf *DateSinceFilter) ToQuery() map[string]interface{} {
	return map[string]interface{}{
		"range": map[string]interface{}{
			dbf.Field: map[string]interface{}{
				"gte": parseTimeSince(dbf.Values[0]),
			},
		},
	}
}

func (ff *DateSinceFilter) GetType() string {
	return "date_since"
}

var regexYear = regexp.MustCompile(`^\d{4}$`)
var regexDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var regexDatestamp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

func parseTimeSince(v string) string {
	v = strings.TrimSpace(v)

	if v == "today" {
		t := time.Now().UTC().Truncate(time.Hour * 24)
		return internal_time.FormatTimeUTC(&t)
	} else if v == "yesterday" {
		t := time.Now().UTC().Add(time.Hour * (-24)).Truncate(time.Hour * 24)
		return internal_time.FormatTimeUTC(&t)
	} else if regexYear.MatchString(v) {
		return v + "-01-01T00:00:00Z"
	} else if regexDate.MatchString(v) {
		return v + "T00:00:00Z"
	} else if regexDatestamp.MatchString(v) {
		return v
	}

	// invalid time: search for time in the future in order to return 0 results
	t := time.Now().UTC().AddDate(100, 0, 0).Truncate(time.Hour * 24)
	return internal_time.FormatTimeUTC(&t)
}

func ToTypeFilter(t string, name string, field string, values ...string) models.Filterable {
	if t == "date_since" {
		f := &DateSinceFilter{}
		f.Name = name
		f.Field = field
		f.Values = values
		return f
	} else if t == "field" {
		f := &FieldFilter{}
		f.Name = name
		f.Field = field
		f.Values = values
		return f
	}

	return nil
}

// filter without facet values
var RegularPublicationFilters = []map[string]string{
	{
		"name":  "created_since",
		"field": "date_created",
		"type":  "date_since",
	},
	{
		"name":  "updated_since",
		"field": "date_updated",
		"type":  "date_since",
	},
}

var RegularDatasetFilters = []map[string]string{
	{
		"name":  "created_since",
		"field": "date_created",
		"type":  "date_since",
	},
	{
		"name":  "updated_since",
		"field": "date_updated",
		"type":  "date_since",
	},
}

func getRegularPublicationFilter(name string, values ...string) models.Filterable {
	for _, cf := range RegularPublicationFilters {
		if cf["name"] == name {
			return ToTypeFilter(cf["type"], cf["name"], cf["field"], values...)
		}
	}
	return nil
}
func getRegularDatasetFilter(name string, values ...string) models.Filterable {
	for _, cf := range RegularPublicationFilters {
		if cf["name"] == name {
			return ToTypeFilter(cf["type"], cf["name"], cf["field"], values...)
		}
	}
	return nil
}
