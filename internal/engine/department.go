package engine

import (
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (e *Engine) SuggestDepartments(q string) ([]models.Completion, error) {
	hits := make([]models.Completion, 0)
	qp := url.Values{}
	qp.Set("q", q)
	if _, err := e.get("/completion/organization", qp, &hits); err != nil {
		return nil, err
	}
	return hits, nil
}
