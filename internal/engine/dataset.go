package engine

import (
	"context"

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
