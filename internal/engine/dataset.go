package engine

import (
	"errors"
	"log"
	"sync"

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

	d, err = e.StorageService.SaveDataset(d)
	if err != nil {
		return nil, err
	}

	if err := e.DatasetSearchService.IndexDataset(d); err != nil {
		log.Printf("error indexing dataset %+v", err)
		return nil, err
	}

	return d, nil
}

func (e *Engine) UpdateDataset(d *models.Dataset) (*models.Dataset, error) {
	d.Vacuum()

	if err := d.Validate(); err != nil {
		log.Printf("%#v", err)
		return nil, err
	}

	d, err := e.StorageService.SaveDataset(d)
	if err != nil {
		return nil, err
	}

	if err := e.DatasetSearchService.IndexDataset(d); err != nil {
		log.Printf("error indexing dataset %+v", err)
		return nil, err
	}

	return d, nil
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

func (e *Engine) GetDatasetPublications(d *models.Dataset) ([]*models.Publication, error) {
	publicationIds := make([]string, len(d.RelatedPublication))
	for _, rp := range d.RelatedPublication {
		publicationIds = append(publicationIds, rp.ID)
	}
	return e.StorageService.GetPublications(publicationIds)
}

func (e *Engine) IndexAllDatasets() (err error) {
	var indexWG sync.WaitGroup

	// indexing channel
	indexC := make(chan *models.Dataset)

	go func() {
		indexWG.Add(1)
		defer indexWG.Done()
		e.DatasetSearchService.IndexDatasets(indexC)
	}()

	// send recs to indexer
	e.StorageService.EachDataset(func(d *models.Dataset) bool {
		indexC <- d
		return true
	})

	close(indexC)

	// wait for indexing to finish
	indexWG.Wait()

	return
}
