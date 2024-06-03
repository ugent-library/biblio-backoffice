package es6

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
)

type SearchServiceConfig struct {
	Addresses        string
	DatasetIndex     string
	PublicationIndex string
	IndexRetention   int // -1: keep all old indexes, >=0: keep x old indexes
}

type SearchService struct {
	client           *elasticsearch.Client
	datasetIndex     string
	publicationIndex string
	indexRetention   int
}

func NewSearchService(c SearchServiceConfig) (backends.SearchService, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: strings.Split(c.Addresses, ","),
	})

	if err != nil {
		return nil, fmt.Errorf("es6.NewSearchService: %w", err)
	}

	return &SearchService{
		client:           client,
		datasetIndex:     c.DatasetIndex,
		publicationIndex: c.PublicationIndex,
		indexRetention:   c.IndexRetention,
	}, nil
}

func (s *SearchService) NewPublicationIndex(r *repositories.Repo) backends.PublicationIndex {
	e := newPublicationIndex(s.client, s.publicationIndex)
	return backends.NewPublicationIndex(e, r)
}

func (s *SearchService) NewPublicationBulkIndexer(config backends.BulkIndexerConfig) (backends.BulkIndexer[*models.Publication], error) {
	docFn := func(p *models.Publication) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedPublication(p))
		return p.ID, doc, err
	}

	return newBulkIndexer(s.client, s.publicationIndex, docFn, config)
}

func (s *SearchService) NewPublicationIndexSwitcher(config backends.BulkIndexerConfig) (backends.IndexSwitcher[*models.Publication], error) {
	settings, err := os.ReadFile("etc/es6/publication.json")
	if err != nil {
		return nil, fmt.Errorf("searchservice.NewPublicationIndexSwitcher: %w", err)
	}

	docFn := func(p *models.Publication) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedPublication(p))
		return p.ID, doc, err
	}

	return newIndexSwitcher(s.client, s.publicationIndex,
		string(settings), s.indexRetention, docFn, config)
}

func (s *SearchService) NewDatasetIndex(r *repositories.Repo) backends.DatasetIndex {
	e := newDatasetIndex(s.client, s.datasetIndex)
	return backends.NewDatasetIndex(e, r)
}

func (s *SearchService) NewDatasetBulkIndexer(config backends.BulkIndexerConfig) (backends.BulkIndexer[*models.Dataset], error) {
	docFn := func(d *models.Dataset) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedDataset(d))
		return d.ID, doc, err
	}

	return newBulkIndexer(s.client, s.datasetIndex, docFn, config)
}

func (s *SearchService) NewDatasetIndexSwitcher(config backends.BulkIndexerConfig) (backends.IndexSwitcher[*models.Dataset], error) {
	settings, err := os.ReadFile("etc/es6/dataset.json")
	if err != nil {
		return nil, fmt.Errorf("searchservice.NewDatasetIndexSwitcher: %w", err)
	}

	docFn := func(d *models.Dataset) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedDataset(d))
		return d.ID, doc, err
	}

	return newIndexSwitcher(s.client, s.datasetIndex,
		string(settings), s.indexRetention, docFn, config)
}
