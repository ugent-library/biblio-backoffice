package es6

import (
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

var publicationFacetFields []string = []string{
	"status",
	"type",
	"faculty",
	"extern",
	"publication_status",
}
var datasetFacetFields []string = []string{"status", "faculty"}
var fixedFacetValues = map[string][]string{
	//"publication_statuses" includes "deleted"
	"status":             vocabularies.Map["visible_publication_statuses"],
	"type":               vocabularies.Map["publication_types"],
	"faculty":            vocabularies.Map["faculties"],
	"extern":             {"true", "false"},
	"publication_status": vocabularies.Map["publication_publishing_statuses"],
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
