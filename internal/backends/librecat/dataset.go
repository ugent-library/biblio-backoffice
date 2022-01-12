package librecat

import (
	"fmt"

	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/forms"
)

func (c *Client) UserDatasets(userID string, args *engine.SearchArgs) (*models.DatasetHits, error) {
	qp, err := forms.Encode(args)
	if err != nil {
		return nil, err
	}
	hits := &models.DatasetHits{}
	if _, err := c.get(fmt.Sprintf("/user/%s/dataset", userID), qp, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (c *Client) GetDataset(id string) (*models.Dataset, error) {
	dataset := &models.Dataset{}
	if _, err := c.get(fmt.Sprintf("/dataset/%s", id), nil, dataset); err != nil {
		return nil, err
	}
	return dataset, nil
}

func (c *Client) ImportUserDatasetByIdentifier(userID, source, identifier string) (*models.Dataset, error) {
	reqData := struct {
		Source          string `json:"source"`
		Identifier      string `json:"identifier"`
		CreationContext string `json:"creation_context"`
	}{
		source,
		identifier,
		"biblio-backend",
	}
	dataset := &models.Dataset{}
	if _, err := c.post(fmt.Sprintf("/user/%s/publication/import", userID), &reqData, dataset); err != nil {
		return nil, err
	}
	return dataset, nil
}

// TODO: set constraint to research_data
func (c *Client) UpdateDataset(dataset *models.Dataset) (*models.Dataset, error) {
	resDataset := &models.Dataset{}
	if _, err := c.put(fmt.Sprintf("/publication/%s", dataset.ID), dataset, resDataset); err != nil {
		return nil, err
	}
	return resDataset, nil
}

func (c *Client) PublishDataset(dataset *models.Dataset) (*models.Dataset, error) {
	oldStatus := dataset.Status
	dataset.Status = "public"
	updatedDataset, err := c.UpdateDataset(dataset)
	dataset.Status = oldStatus
	return updatedDataset, err
}

func (c *Client) GetDatasetPublications(id string) ([]*models.Publication, error) {
	pubs := make([]*models.Publication, 0)
	if _, err := c.get(fmt.Sprintf("/dataset/%s/publication", id), nil, &pubs); err != nil {
		return nil, err
	}
	return pubs, nil
}

func (c *Client) DeleteDataset(id string) error {
	_, err := c.delete(fmt.Sprintf("/publication/%s", id), nil, nil)
	return err
}
