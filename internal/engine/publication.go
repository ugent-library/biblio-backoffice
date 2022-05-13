package engine

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

// TODO move to workflow
func (e *Engine) UpdatePublication(p *models.Publication) error {
	p.Vacuum()

	if err := p.Validate(); err != nil {
		log.Printf("%#v", err)
		return err
	}

	if err := e.Store.StorePublication(p); err != nil {
		return err
	}

	// if err := e.PublicationSearchService.IndexPublication(p); err != nil {
	// 	log.Printf("error indexing publication %+v", err)
	// 	return err
	// }

	return nil
}

// TODO make query dsl package
func (e *Engine) Publications(args *models.SearchArgs) (*models.PublicationHits, error) {
	args = args.Clone().WithFilter("status", "private", "public")
	return e.PublicationSearchService.SearchPublications(args)
}

// TODO make query dsl package
func (e *Engine) UserPublications(userID string, args *models.SearchArgs) (*models.PublicationHits, error) {
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
	return e.PublicationSearchService.SearchPublications(args)
}

// TODO make workflow
func (e *Engine) BatchPublishPublications(userID string, args *models.SearchArgs) (err error) {
	var hits *models.PublicationHits
	for {
		hits, err = e.UserPublications(userID, args)
		for _, pub := range hits.Hits {
			pub.Status = "public"
			if err = e.UpdatePublication(pub); err != nil {
				break
			}
		}
		if !hits.NextPage() {
			break
		}
		args.Page = args.Page + 1
	}
	return
}

// TODO make make controller helper or eliminate
func (e *Engine) GetPublicationDatasets(p *models.Publication) ([]*models.Dataset, error) {
	datasetIds := make([]string, len(p.RelatedDataset))
	for _, rd := range p.RelatedDataset {
		datasetIds = append(datasetIds, rd.ID)
	}
	return e.Store.GetDatasets(datasetIds)
}

// TODO make model helper method and move to controller
func (e *Engine) AddPublicationDataset(p *models.Publication, d *models.Dataset) error {
	return e.Store.Transaction(context.Background(), func(s backends.Store) error {
		if !p.HasRelatedDataset(d.ID) {
			p.RelatedDataset = append(p.RelatedDataset, models.RelatedDataset{ID: d.ID})
			if err := s.StorePublication(p); err != nil {
				return err
			}
		}
		if !d.HasRelatedPublication(p.ID) {
			d.RelatedPublication = append(d.RelatedPublication, models.RelatedPublication{ID: p.ID})
			if err := s.StoreDataset(d); err != nil {
				return err
			}
		}

		// TODO ensure consistency
		// if err := e.DatasetSearchService.IndexDataset(d); err != nil {
		// 	log.Printf("error indexing dataset: %v", err)
		// }
		// if err := e.PublicationSearchService.IndexPublication(p); err != nil {
		// 	log.Printf("error indexing publication: %v", err)
		// }

		return nil
	})
}

// TODO make model helper method and move to controller
func (e *Engine) RemovePublicationDataset(p *models.Publication, d *models.Dataset) error {
	return e.Store.Transaction(context.Background(), func(s backends.Store) error {
		if p.HasRelatedDataset(d.ID) {
			p.RemoveRelatedDataset(d.ID)
			if err := s.StorePublication(p); err != nil {
				return err
			}
		}
		if d.HasRelatedPublication(p.ID) {
			d.RemoveRelatedPublication(p.ID)
			if err := s.StoreDataset(d); err != nil {
				return err
			}
		}

		// TODO ensure consistency
		// if err := e.DatasetSearchService.IndexDataset(d); err != nil {
		// 	log.Printf("error indexing dataset: %v", err)
		// }
		// if err := e.PublicationSearchService.IndexPublication(p); err != nil {
		// 	log.Printf("error indexing publication: %v", err)
		// }

		return nil
	})
}

// TODO make workflow
func (e *Engine) ImportUserPublicationByIdentifier(userID, source, identifier string) (*models.Publication, error) {
	s, ok := e.PublicationSources[source]
	if !ok {
		return nil, errors.New("unknown dataset source")
	}
	p, err := s.GetPublication(identifier)
	if err != nil {
		return nil, err
	}

	p.ID = uuid.NewString()
	p.CreatorID = userID
	p.UserID = userID
	p.Status = "private"
	p.Classification = "U"

	if err := e.UpdatePublication(p); err != nil {
		return nil, err
	}

	return p, nil
}

// TODO make workflow
func (e *Engine) ImportUserPublications(userID, source string, file io.Reader) (string, error) {
	batchID := uuid.New().String()
	decFactory, ok := e.PublicationDecoders[source]
	if !ok {
		return "", errors.New("unknown publication source")
	}
	dec := decFactory(file)

	var indexWG sync.WaitGroup

	// indexing channel
	indexC := make(chan *models.Publication)

	// start bulk indexer
	go func() {
		indexWG.Add(1)
		defer indexWG.Done()
		e.PublicationSearchService.IndexPublications(indexC)
	}()

	var importErr error
	for {
		p := models.Publication{
			ID:             uuid.NewString(),
			BatchID:        batchID,
			Status:         "private",
			Classification: "U",
			CreatorID:      userID,
			UserID:         userID,
		}
		if err := dec.Decode(&p); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			importErr = err
			break
		}
		if err := p.Validate(); err != nil {
			importErr = err
			break
		}
		if err := e.Store.StorePublication(&p); err != nil {
			importErr = err
			break
		}

		indexC <- &p
	}

	// close indexing channel when all recs are stored
	close(indexC)
	// wait for indexing to finish
	indexWG.Wait()

	// TODO rollback if error
	if importErr != nil {
		return "", importErr
	}

	return batchID, nil
}

// TODO point to biblio frontend
func (c *Engine) ServePublicationThumbnail(fileURL string, w http.ResponseWriter, r *http.Request) {
	// panic("not implemented")
}

// TODO make workflow
func (e *Engine) IndexAllPublications() (err error) {
	var indexWG sync.WaitGroup

	// indexing channel
	indexC := make(chan *models.Publication)

	go func() {
		indexWG.Add(1)
		defer indexWG.Done()
		e.PublicationSearchService.IndexPublications(indexC)
	}()

	// send recs to indexer
	e.Store.EachPublication(func(p *models.Publication) bool {
		indexC <- p
		return true
	})

	close(indexC)

	// wait for indexing to finish
	indexWG.Wait()

	return
}
