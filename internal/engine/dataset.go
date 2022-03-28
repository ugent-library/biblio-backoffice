package engine

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ugent-library/biblio-backend/internal/models"
)

// func (e *Engine) StoreDataset(d *models.Dataset) error {
// 	ctx := context.Background()

// 	workflowOptions := client.StartWorkflowOptions{
// 		ID:        "store-dataset-" + uuid.New().String(),
// 		TaskQueue: "store-dataset",
// 	}

// 	workflowRun, err := e.Temporal.ExecuteWorkflow(ctx, workflowOptions, dataset.StoreDatasetWorkflow, d)
// 	if err != nil {
// 		return err
// 	}

// 	// wait for workflow to finish
// 	return workflowRun.Get(ctx, nil)
// }

func (e *Engine) ImportUserDatasetByIdentifier(userID, source, identifier string) (*models.Dataset, error) {
	s, ok := e.DatasetSources[source]
	if !ok {
		return nil, errors.New("unknown dataset source")
	}
	d, err := s.GetDataset(identifier)
	if err != nil {
		return nil, err
	}
	d.Vacuum()
	d.CreatorID = userID
	d.UserID = userID
	d.Status = "private"

	d, err = e.DatasetService.CreateDataset(d)
	if err != nil {
		return nil, err
	}

	if err := e.DatasetSearchService.IndexDataset(d); err != nil {
		log.Printf("error indexing dataset %+v", err)
		return nil, err
	}

	return d, nil
}

func (e *Engine) GetDataset(id string) (*models.Dataset, error) {
	return e.DatasetService.GetDataset(id)
}

func (e *Engine) UpdateDataset(d *models.Dataset) (*models.Dataset, error) {
	if err := e.ValidateDataset(d); err != nil {
		log.Printf("%#v", err)
		return nil, err
	}

	d, err := e.DatasetService.UpdateDataset(d)
	if err != nil {
		return nil, err
	}

	if err := e.DatasetSearchService.IndexDataset(d); err != nil {
		log.Printf("error indexing dataset %+v", err)
		return nil, err
	}

	return d, nil
}

func (e *Engine) PublishDataset(d *models.Dataset) (*models.Dataset, error) {
	d.Status = "public"
	return e.UpdateDataset(d)
}

func (e *Engine) DeleteDataset(d *models.Dataset) error {
	d.Status = "deleted"
	_, err := e.UpdateDataset(d)
	return err
}

func (e *Engine) Datasets(args *models.SearchArgs) (*models.DatasetHits, error) {
	args = args.Clone().WithFilter("status", "private", "public")
	return e.DatasetSearchService.SearchDatasets(args)
}

func (e *Engine) UserDatasets(userID string, args *models.SearchArgs) (*models.DatasetHits, error) {
	args = args.Clone().WithFilter("status", "private", "public")
	switch args.FilterFor("scope") {
	case "created":
		args.WithFilter("creator_id", userID)
	case "contributed":
		args.WithFilter("author.id", userID)
	default:
		args.WithFilter("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")
	return e.DatasetSearchService.SearchDatasets(args)
}

func (e *Engine) ValidateDataset(d *models.Dataset) error {
	m := make(map[string]interface{})
	dJSON, err := json.Marshal(d)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(dJSON, &m); err != nil {
		return err
	}

	return e.datasetSchema.Validate(m)
}
