package engine

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/ugent-library/biblio-backend/internal/models"
)

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
	d.ID = uuid.NewString()
	d.CreatorID = userID
	d.UserID = userID
	d.Status = "private"

	if err = e.Store.UpdateDataset(d); err != nil {
		return nil, err
	}

	// if err := e.DatasetSearchService.IndexDataset(d); err != nil {
	// 	log.Printf("error indexing dataset %+v", err)
	// 	return nil, err
	// }

	return d, nil
}

func (e *Engine) Datasets(args *models.SearchArgs) (*models.DatasetHits, error) {
	searcher := e.DatasetSearchService
	searcher = searcher.WithScope("status", "private", "public")
	return searcher.Search(args)
}

func (e *Engine) UserDatasets(userID string, args *models.SearchArgs) (*models.DatasetHits, error) {
	searcher := e.DatasetSearchService
	searcher = searcher.WithScope("status", "private", "public")
	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", userID)
	case "contributed":
		searcher = searcher.WithScope("author.id", userID)
	default:
		searcher = searcher.WithScope("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")
	return searcher.Search(args)
}

func (e *Engine) GetDatasetPublications(d *models.Dataset) ([]*models.Publication, error) {
	publicationIds := make([]string, len(d.RelatedPublication))
	for _, rp := range d.RelatedPublication {
		publicationIds = append(publicationIds, rp.ID)
	}
	return e.Store.GetPublications(publicationIds)
}

func (e *Engine) IndexAllDatasets() (err error) {
	var indexWG sync.WaitGroup

	// indexing channel
	indexC := make(chan *models.Dataset)

	go func() {
		indexWG.Add(1)
		defer indexWG.Done()
		e.DatasetSearchService.IndexMultiple(indexC)
	}()

	// send recs to indexer
	e.Store.EachDataset(func(d *models.Dataset) bool {
		indexC <- d
		return true
	})

	close(indexC)

	// wait for indexing to finish
	indexWG.Wait()

	return
}
