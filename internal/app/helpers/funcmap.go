package helpers

import (
	"fmt"
	"html/template"

	"github.com/rvflash/elapsed"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/friendly"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"searchArgs":      models.NewSearchArgs,
		"timeElapsed":     elapsed.LocalTime,
		"formatRange":     FormatRange,
		"formatBool":      FormatBool,
		"formatBytes":     friendly.Bytes,
		"resolveDOI":      ResolveDOI,
		"resolveWOSID":    ResolveWOSID,
		"resolvePubMedID": ResolvePubMedID,
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

func ResolveDOI(id string) string {
	return "https://doi.org/" + id
}

func ResolveWOSID(id string) string {
	return "https://www.webofscience.com/wos/woscc/full-record/" + id
}

func ResolvePubMedID(id string) string {
	return "https://www.ncbi.nlm.nih.gov/pubmed/" + id
}
