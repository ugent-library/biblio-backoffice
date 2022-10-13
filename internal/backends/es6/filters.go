package es6

import (
	"regexp"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
	internal_time "github.com/ugent-library/biblio-backend/internal/time"
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

func isYear(v string) bool {
	return regexYear.MatchString(v)
}

func isDate(v string) bool {
	return regexDate.MatchString(v)
}

func parseTimeSince(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))

	if v == "today" {
		t := time.Now().UTC().Truncate(time.Hour * 24)
		return internal_time.FormatTimeUTC(&t)
	} else if v == "yesterday" {
		t := time.Now().UTC().Add(time.Hour * (-24)).Truncate(time.Hour * 24)
		return internal_time.FormatTimeUTC(&t)
	} else if isYear(v) {
		return v + "-01-01T00:00:00Z"
	} else if isDate(v) {
		return v + "T00:00:00Z"
	}

	//invalid time: search for time in the future in order to return 0 results
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
