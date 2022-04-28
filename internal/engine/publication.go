package engine

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/ugent-library/biblio-backend/internal/models"
)

// TODO move to workflow
func (e *Engine) UpdatePublication(p *models.Publication) (*models.Publication, error) {
	p.Vacuum()

	if err := p.Validate(); err != nil {
		log.Printf("%#v", err)
		return nil, err
	}

	p, err := e.StorageService.SavePublication(p)
	if err != nil {
		return nil, err
	}

	if err := e.PublicationSearchService.IndexPublication(p); err != nil {
		log.Printf("error indexing publication %+v", err)
		return nil, err
	}

	return p, nil
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
			if _, err = e.UpdatePublication(pub); err != nil {
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
	return e.StorageService.GetDatasets(datasetIds)
}

// TODO make model helper method and move to controller
func (e *Engine) AddPublicationDataset(p *models.Publication, d *models.Dataset) (*models.Publication, error) {
	tx, err := e.StorageService.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if !p.HasRelatedDataset(d.ID) {
		p.RelatedDataset = append(p.RelatedDataset, models.RelatedDataset{ID: d.ID})
		savedP, err := tx.SavePublication(p)
		if err != nil {
			return nil, err
		}
		p = savedP
	}
	if !d.HasRelatedPublication(p.ID) {
		d.RelatedPublication = append(d.RelatedPublication, models.RelatedPublication{ID: p.ID})
		if _, err := tx.SaveDataset(d); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return p, nil
}

// TODO make model helper method and move to controller
func (e *Engine) RemovePublicationDataset(p *models.Publication, d *models.Dataset) (*models.Publication, error) {
	tx, err := e.StorageService.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if p.HasRelatedDataset(d.ID) {
		var newRelatedDatasets []models.RelatedDataset
		for _, rd := range p.RelatedDataset {
			if rd.ID != d.ID {
				newRelatedDatasets = append(newRelatedDatasets, rd)
			}
		}
		p.RelatedDataset = newRelatedDatasets
		savedP, err := tx.SavePublication(p)
		if err != nil {
			return nil, err
		}
		p = savedP
	}
	if d.HasRelatedPublication(p.ID) {
		var newRelatedPublications []models.RelatedPublication
		for _, rd := range d.RelatedPublication {
			if rd.ID != d.ID {
				newRelatedPublications = append(newRelatedPublications, rd)
			}
		}
		d.RelatedPublication = newRelatedPublications
		if _, err := tx.SaveDataset(d); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return p, nil
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

	p.CreatorID = userID
	p.UserID = userID
	p.Status = "private"
	p.Classification = "U"

	return e.UpdatePublication(p)
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
			ID:             uuid.New().String(),
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
		savedP, err := e.StorageService.SavePublication(&p)
		if err != nil {
			importErr = err
			break
		}

		indexC <- savedP
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
	e.StorageService.EachPublication(func(p *models.Publication) bool {
		indexC <- p
		return true
	})

	close(indexC)

	// wait for indexing to finish
	indexWG.Wait()

	return
}
