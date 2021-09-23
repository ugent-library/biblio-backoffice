package engine

import (
	"fmt"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/forms"
)

func (e *Engine) UserDatasets(userID string, args *SearchArgs) (*models.PublicationHits, error) {
	qp, err := forms.Encode(args)
	if err != nil {
		return nil, err
	}
	hits := &models.PublicationHits{}
	if _, err := e.get(fmt.Sprintf("/user/%s/dataset", userID), qp, hits); err != nil {
		return nil, err
	}
	return hits, nil
}
