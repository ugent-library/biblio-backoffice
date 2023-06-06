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
