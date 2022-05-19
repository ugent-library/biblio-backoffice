package es6

import "strings"

func ParseScope(field string, terms ...string) M {
	orFields := strings.Split(field, "|")
	if len(orFields) > 1 {
		orFilters := make([]M, 0, len(orFields))
		for _, orField := range orFields {
			orFilters = append(orFilters, M{
				"terms": M{orField: terms},
			})
		}
		return M{
			"bool": M{
				"should":               orFilters,
				"minimum_should_match": "1",
			},
		}
	} else if strings.HasPrefix(field, "!") {
		return M{
			"bool": M{
				"must_not": M{"terms": M{field[1:]: terms}},
			},
		}
	} else {
		return M{"terms": M{field: terms}}
	}
}
