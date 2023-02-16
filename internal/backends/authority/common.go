package authority

import (
	"regexp"

	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type personSearchEnvelope struct {
	Hits struct {
		Total int `json:"total"`
		Hits  []struct {
			ID     string        `json:"_id"`
			Source models.Person `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

var regexMultipleSpaces = regexp.MustCompile(`\s+`)
var regexNoBrackets = regexp.MustCompile(`[\[\]()\{\}]`)
