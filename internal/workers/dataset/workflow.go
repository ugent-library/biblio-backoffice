package dataset

import (
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"go.temporal.io/sdk/workflow"
)

type StoreActivities struct {
	DatasetService backends.DatasetService
}

func StoreWorkflow(ctx workflow.Context, dataset *models.Dataset) error {
	return nil
}
