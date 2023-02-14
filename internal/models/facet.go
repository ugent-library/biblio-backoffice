package models

type FacetValues []Facet

type Facet struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

func (fv FacetValues) HasMatches() bool {
	for _, v := range fv {
		if v.Count > 0 {
			return true
		}
	}
	return false
}

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
	GetType() string
	ToQuery() map[string]interface{}
}

type Filter struct {
	Name   string
	Field  string
	Values []string
}

type BaseFilter struct {
	Filter
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
