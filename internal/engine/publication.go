package engine

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (e *Engine) GetPublication(id string) (*models.Publication, error) {
	return e.StorageService.GetPublication(id)
}

func (e *Engine) UpdatePublication(d *models.Publication) (*models.Publication, error) {
	if err := d.Validate(); err != nil {
		log.Printf("%#v", err)
		return nil, err
	}

	d, err := e.StorageService.UpdatePublication(d)
	if err != nil {
		return nil, err
	}

	if err := e.PublicationSearchService.IndexPublication(d); err != nil {
		log.Printf("error indexing publication %+v", err)
		return nil, err
	}

	return d, nil
}

func (e *Engine) Publications(args *models.SearchArgs) (*models.PublicationHits, error) {
	args = args.Clone().WithFilter("status", "private", "public")
	return e.PublicationSearchService.SearchPublications(args)
}

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

func (e *Engine) GetPublicationDatasets(id string) ([]*models.Dataset, error) {
	return nil, nil
}

func (e *Engine) AddPublicationDataset(id, datasetID string) error {
	return errors.New("not implemented")
}

func (e *Engine) RemovePublicationDataset(id, datasetID string) error {
	return errors.New("not implemented")
}

func (e *Engine) ImportUserPublicationByIdentifier(userID, source, identifier string) (*models.Publication, error) {
	return nil, errors.New("not implemented")
}

func (e *Engine) ImportUserPublications(userID, source string, file io.Reader) (string, error) {
	return "", errors.New("not implemented")
}

func (c *Engine) ServePublicationFile(fileURL string, w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (e *Engine) AddPublicationFile(id string, pubFile *models.PublicationFile, file io.Reader) error {
	return errors.New("not implemented")
}
func (e *Engine) UpdatePublicationFile(id string, pubFile *models.PublicationFile) error {
	return errors.New("not implemented")
}

func (e *Engine) RemovePublicationFile(id, fileID string) error {
	return errors.New("not implemented")
}
