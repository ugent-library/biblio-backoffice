package engine

import (
	"fmt"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/forms"
)

func (e *Engine) UserDatasets(userID string, args *SearchArgs) (*models.DatasetHits, error) {
	qp, err := forms.Encode(args)
	if err != nil {
		return nil, err
	}
	hits := &models.DatasetHits{}
	if _, err := e.get(fmt.Sprintf("/user/%s/dataset", userID), qp, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

// TODO: set constraint to research_data
func (e *Engine) GetDataset(id string) (*models.Dataset, error) {
	dataset := &models.Dataset{}
	if _, err := e.get(fmt.Sprintf("/publication/%s", id), nil, dataset); err != nil {
		return nil, err
	}
	return dataset, nil
}

// TODO: set constraint to research_data
func (e *Engine) UpdateDataset(dataset *models.Dataset) (*models.Dataset, error) {
	resDataset := &models.Dataset{}
	if _, err := e.put(fmt.Sprintf("/publication/%s", dataset.ID), dataset, resDataset); err != nil {
		return nil, err
	}
	return resDataset, nil
}
