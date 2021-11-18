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

func (e *Engine) GetDataset(id string) (*models.Dataset, error) {
	dataset := &models.Dataset{}
	if _, err := e.get(fmt.Sprintf("/dataset/%s", id), nil, dataset); err != nil {
		return nil, err
	}
	return dataset, nil
}

func (e *Engine) ImportUserDatasetByIdentifier(userID, source, identifier string) (*models.Dataset, error) {
	reqData := struct {
		Source     string `json:"source"`
		Identifier string `json:"identifier"`
	}{
		source,
		identifier,
	}
	dataset := &models.Dataset{}
	if _, err := e.post(fmt.Sprintf("/user/%s/publication/import", userID), &reqData, dataset); err != nil {
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

func (e *Engine) PublishDataset(dataset *models.Dataset) (*models.Dataset, error) {
	dataset.Status = "public"
	return e.UpdateDataset(dataset)
}

func (e *Engine) GetDatasetPublications(id string) ([]*models.Publication, error) {
	pubs := make([]*models.Publication, 0)
	if _, err := e.get(fmt.Sprintf("/dataset/%s/publication", id), nil, &pubs); err != nil {
		return nil, err
	}
	return pubs, nil
}
