package authority

import (
	"regexp"

	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type personSearchEnvelope struct {
	Hits struct {
		Total int `json:"total"`
		Hits  []struct {
			ID     string `json:"_id"`
			Source struct {
				*models.Person
				Department []struct {
					ID string `json:"_id"`
				} `json:"department"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

var (
	regexMultipleSpaces = regexp.MustCompile(`\s+`)
	regexNoBrackets     = regexp.MustCompile(`[\[\]()\{\}]`)
)
