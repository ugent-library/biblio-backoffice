package spdx

import (
	"github.com/ugent-library/biblio-backend/internal/models"
)

type localSuggester struct {
}

func New() *localSuggester {
	return &localSuggester{}
}

func (c *localSuggester) SuggestLicenses(q string) ([]models.Completion, error) {
	hits := make([]models.Completion, 0)
	return hits, nil
}
