package es6

import (
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
)

type facetDefinition struct {
	config M
}

var facetDefinitions = map[string]facetDefinition{
	"reviewer_tags": {
		config: M{
			"terms": M{
				"field":         "reviewer_tags",
				"order":         M{"_key": "asc"},
				"size":          999,
				"min_doc_count": 0,
			},
		},
	},
	"year": {
		config: M{
			"terms": M{
				"field":         "year",
				"order":         M{"_key": "desc"},
				"size":          999,
				"min_doc_count": 0,
			},
		},
	},
	"wos_type": {
		config: M{
			"terms": M{
				"field":         "wos_type",
				"order":         M{"_key": "asc"},
				"size":          999,
				"min_doc_count": 0,
			},
		},
	},
}

// TODO remove this when all facets have a static definition above
func defaultFacetDefinition(field string) facetDefinition {
	return facetDefinition{
		config: M{
			"terms": M{
				"field":         field,
				"order":         M{"_key": "asc"},
				"size":          100,
				"min_doc_count": 0,
			},
		},
	}
}

// facets that also work as filters
var fixedFacetValues = map[string][]string{
	//"publication_statuses" includes "deleted"
	"classification":     vocabularies.Map["publication_classifications"],
	"extern":             {"true", "false"},
	"faculty_id":         append([]string{backends.MissingValue}, vocabularies.Map["faculties"]...),
	"file_relation":      vocabularies.Map["publication_file_relations"],
	"has_message":        {"true", "false"},
	"has_files":          {"true", "false"},
	"legacy":             {"true", "false"},
	"locked":             {"true", "false"},
	"publication_status": append([]string{backends.MissingValue}, vocabularies.Map["publication_publishing_statuses"]...),
	"status":             vocabularies.Map["visible_publication_statuses"],
	"type":               vocabularies.Map["publication_types"],
	"miscellaneous_type": vocabularies.Map["miscellaneous_types"],
	"vabb_type":          vocabularies.Map["publication_vabb_types"],
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
