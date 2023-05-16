package es6

import (
	"encoding/json"

	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type DatasetSearchService struct {
	Client
}

func NewDatasetSearchService(c Client) *DatasetSearchService {
	return &DatasetSearchService{Client: c}
}

func (ds *DatasetSearchService) NewIndex() backends.DatasetIndex {
	return newDatasetIndex(ds.Client)
}

func (ds *DatasetSearchService) NewBulkIndexer(config backends.BulkIndexerConfig) (backends.BulkIndexer[*models.Dataset], error) {
	docFn := func(d *models.Dataset) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedDataset(d))
		return d.ID, doc, err
	}
	return newBulkIndexer(ds.Client.es, ds.Client.Index, docFn, config)
}

func (ds *DatasetSearchService) NewIndexSwitcher(config backends.BulkIndexerConfig) (backends.IndexSwitcher[*models.Dataset], error) {
	docFn := func(d *models.Dataset) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedDataset(d))
		return d.ID, doc, err
	}
	return newIndexSwitcher(ds.Client.es, ds.Client.Index, ds.Client.Settings, ds.Client.IndexRetention, docFn, config)
}
