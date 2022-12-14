package es6

import (
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
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
				"min_doc_count": 1,
			},
		},
	},
	"year": {
		config: M{
			"terms": M{
				"field":         "year",
				"order":         M{"_key": "desc"},
				"size":          999,
				"min_doc_count": 1,
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
				"size":          999,
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
	"faculty":            vocabularies.Map["faculties"],
	"file.relation":      vocabularies.Map["publication_file_relations"],
	"has_message":        {"true", "false"},
	"legacy":             {"true", "false"},
	"locked":             {"true", "false"},
	"publication_status": vocabularies.Map["publication_publishing_statuses"],
	"status":             vocabularies.Map["visible_publication_statuses"],
	"type":               vocabularies.Map["publication_types"],
	"vabb_type":          vocabularies.Map["publication_vabb_types"],
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
