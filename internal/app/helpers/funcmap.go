package helpers

import (
	"fmt"
	"html/template"

	"github.com/rvflash/elapsed"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"searchArgs":  models.NewSearchArgs,
		"timeElapsed": elapsed.LocalTime,
		"formatRange": FormatRange,
		"formatBool":  FormatBool,
		"filterBadge": FilterBadge,
	}
}

func FormatRange(start, end string) string {
	var v string
	if len(start) > 0 && len(end) > 0 && start == end {
		v = start
	} else if len(start) > 0 && len(end) > 0 {
		v = fmt.Sprintf("%s - %s", start, end)
	} else if len(start) > 0 {
		v = fmt.Sprintf("%s -", start)
	} else if len(end) > 0 {
		v = fmt.Sprintf("- %s", end)
	}

	return v
}

func FormatBool(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

func FilterBadge(searchArgs *models.SearchArgs, field string, fieldFacets []models.Facet) string {

	// any filter is selected
	if len(searchArgs.FiltersFor(field)) > 0 {
		return "badge-primary"
	}

	// valuable options available
	for _, facet := range fieldFacets {
		if facet.Count > 0 {
			return "badge-default"
		}
	}

	// no valuable options available (all options return 0 results)
	return "badge-secondary"
}
