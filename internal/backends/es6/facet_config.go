package es6

import "github.com/ugent-library/biblio-backend/internal/models"

var publicationFacetFields []string = []string{"status", "type", "faculty"}
var datasetFacetFields []string = []string{"status", "faculty"}
var fixedFacetValues = map[string][]string{
	"status": {
		"new",
		"private",
		"public",
		"returned",
	},
	"type": {
		"book",
		"book_chapter",
		"book_editor",
		"conference",
		"dissertation",
		"issue_editor",
		"journal_article",
		"miscellaneous",
	},
	"faculty": {
		"CA",
		"DS",
		"DI",
		"EB",
		"FW",
		"GE",
		"LA",
		"LW",
		"PS",
		"PP",
		"RE",
		"TW",
		"WE",
		"GUK",
		"UZGent",
		"HOART",
		"HOGENT",
		"HOWEST",
		"IBBT",
		"IMEC",
		"VIB",
	},
}

func reorderFacets(t string, facets []models.Facet) []models.Facet {
	fixedValues, e := fixedFacetValues[t]

	//no fixed order defined
	if !e {
		return facets
	}

	//fixed order is defined
	newFacets := make([]models.Facet, 0, len(facets))

	for _, fixedVal := range fixedValues {
		foundFacet := false
		for _, facet := range facets {
			if fixedVal == facet.Value {
				newFacets = append(newFacets, facet)
				foundFacet = true
				break
			}
		}
		/*
			min_doc_count: 0 does not ensure that all possible values
			are there, especially if some were never encountered before.
		*/
		if !foundFacet {
			newFacets = append(newFacets, models.Facet{
				Value: fixedVal,
				Count: 0,
			})
		}
	}

	return newFacets
}
