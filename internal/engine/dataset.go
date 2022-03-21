package engine

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/workers/dataset"
	"go.temporal.io/sdk/client"
)

func (e *Engine) StoreDataset(d *models.Dataset) error {
	ctx := context.Background()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "store-dataset-" + uuid.New().String(),
		TaskQueue: "store-dataset",
	}

	workflowRun, err := e.Temporal.ExecuteWorkflow(ctx, workflowOptions, dataset.StoreDatasetWorkflow, d)
	if err != nil {
		return err
	}

	// wait for workflow to finish
	return workflowRun.Get(ctx, nil)
}

func (e *Engine) ImportUserDatasetByIdentifier(userID, source, identifier string) (*models.Dataset, error) {
	s, ok := e.DatasetSources[source]
	if !ok {
		return nil, errors.New("unknown dataset source")
	}

	d, err := s.GetDataset(identifier)
	if err != nil {
		return nil, err
	}

	d.ID = uuid.NewString()
	d.UserID = userID
	d.Status = "private"

	return e.CreateDataset(d)
}

func (e *Engine) PublishDataset(d *models.Dataset) (*models.Dataset, error) {
	d.Status = "public"
	return e.UpdateDataset(d)
}
