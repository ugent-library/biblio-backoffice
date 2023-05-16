package es6

import (
	"encoding/json"

	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type PublicationSearchService struct {
	Client
	scopes []M
}

func NewPublicationSearchService(c Client) *PublicationSearchService {
	return &PublicationSearchService{Client: c}
}

func (ps *PublicationSearchService) NewIndex() backends.PublicationIndex {
	return newPublicationIndex(ps.Client)
}

func (ps *PublicationSearchService) NewBulkIndexer(config backends.BulkIndexerConfig) (backends.BulkIndexer[*models.Publication], error) {
	docFn := func(p *models.Publication) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedPublication(p))
		return p.ID, doc, err
	}
	return newBulkIndexer(ps.Client.es, ps.Client.Index, docFn, config)
}

func (ps *PublicationSearchService) NewIndexSwitcher(config backends.BulkIndexerConfig) (backends.IndexSwitcher[*models.Publication], error) {
	docFn := func(p *models.Publication) (string, []byte, error) {
		doc, err := json.Marshal(NewIndexedPublication(p))
		return p.ID, doc, err
	}
	return newIndexSwitcher(ps.Client.es, ps.Client.Index,
		ps.Client.Settings, ps.Client.IndexRetention, docFn, config)
}
