package engine

import (
	"encoding/json"
	"fmt"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
)

func (e *Engine) UserPublications(userID string, args *SearchArgs) (*models.PublicationHits, error) {
	qp, err := forms.Encode(args)
	if err != nil {
		return nil, err
	}
	hits := &models.PublicationHits{}
	if _, err := e.get(fmt.Sprintf("/user/%s/publication", userID), qp, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Engine) GetPublication(id string) (*models.Publication, error) {
	pub := &models.Publication{}
	if _, err := e.get(fmt.Sprintf("/publication/%s", id), nil, pub); err != nil {
		return nil, err
	}
	return pub, nil
}

func (e *Engine) UpdatePublication(pub *models.Publication) (*models.Publication, jsonapi.Errors) {
	resPub := &models.Publication{}
	if _, err := e.put(fmt.Sprintf("/publication/%s", pub.ID), pub, resPub); err != nil {
		var errors jsonapi.Errors
		_ = json.Unmarshal([]byte(err.Error()), &errors)

		return nil, errors
	}
	return resPub, nil
}
