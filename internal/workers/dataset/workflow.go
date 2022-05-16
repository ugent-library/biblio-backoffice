package dataset

import (
	"time"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	DatasetService backends.Store
}

func StoreDatasetWorkflow(ctx workflow.Context, dataset *models.Dataset) (err error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Second * 10,
			MaximumAttempts:    3,
		},
	})

	var a *Activities

	err = workflow.ExecuteLocalActivity(ctx, a.StoreDatasetInRepository, dataset).Get(ctx, nil)

	return
}

func (a *Activities) StoreDatasetInRepository(dataset *models.Dataset) error {
	return a.DatasetService.UpdateDataset(dataset)
}
