package es6

import (
	"regexp"
	"strings"
	"time"

	internal_time "github.com/ugent-library/biblio-backoffice/time"
)

/*
Name: 		public field name (e.g. url query parameter)
Field: 		internal field name (e.g. elasticsearch field)
Values:		array of string values
Type:		type of filter. To distinguish from other filters
ToQuery:	convert and return search engine specific filter
*/
type Filterable interface {
	GetName() string
	GetField() string
	GetValues() []string
	ToQuery() map[string]any
}

type BaseFilter struct {
	Name   string
	Field  string
	Values []string
}

func (bf *BaseFilter) GetName() string {
	return bf.Name
}

func (bf *BaseFilter) GetField() string {
	return bf.Field
}

func (bf *BaseFilter) GetValues() []string {
	return bf.Values
}

// regular field filter: accepts syntax in the filter value
type FieldFilter struct {
	BaseFilter
}

func (ff *FieldFilter) ToQuery() map[string]any {
	return ParseScope(ff.Name, ff.Values...)
}

// date filter
type DateSinceFilter struct {
	BaseFilter
}

func (dbf *DateSinceFilter) ToQuery() map[string]any {
	var regexYear = regexp.MustCompile(`^\d{4}$`)
	var regexDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	var regexDatestamp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

	t := strings.TrimSpace(dbf.Values[0])
	fdt := ""

	if t == "today" {
		dt := time.Now().UTC().Truncate(time.Hour * 24)
		fdt = internal_time.FormatTimeUTC(&dt)
	} else if t == "yesterday" {
		dt := time.Now().UTC().Add(time.Hour * (-24)).Truncate(time.Hour * 24)
		fdt = internal_time.FormatTimeUTC(&dt)
	} else if regexYear.MatchString(t) {
		fdt = t + "-01-01T00:00:00Z"
	} else if regexDate.MatchString(t) {
		fdt = t + "T00:00:00Z"
	} else if regexDatestamp.MatchString(t) {
		fdt = t
	} else {
		// invalid time: search for time in the future in order to return 0 results
		dt := time.Now().UTC().AddDate(100, 0, 0).Truncate(time.Hour * 24)
		fdt = internal_time.FormatTimeUTC(&dt)
	}

	return map[string]any{
		"range": map[string]any{
			dbf.Field: map[string]any{
				"gte": fdt,
			},
		},
	}
}

func ToTypeFilter(t string, name string, field string, values []string) Filterable {
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

func getRegularPublicationFilter(name string, values []string) Filterable {
	for _, cf := range RegularPublicationFilters {
		if cf["name"] == name {
			return ToTypeFilter(cf["type"], cf["name"], cf["field"], values)
		}
	}
	return nil
}
func getRegularDatasetFilter(name string, values []string) Filterable {
	for _, cf := range RegularPublicationFilters {
		if cf["name"] == name {
			return ToTypeFilter(cf["type"], cf["name"], cf["field"], values)
		}
	}
	return nil
}
