package engine

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/forms"
)

func (e *Engine) UserPublications(userID string, args *SearchArgs) (*models.PublicationHits, error) {
	qp, err := forms.Encode(args)
	if err != nil {
		return nil, err
	}
	hits := &models.PublicationHits{}
	if _, err := e.get(fmt.Sprintf("/user/%s/publication", userID), qp, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Engine) GetPublication(id string) (*models.Publication, error) {
	pub := &models.Publication{}
	if _, err := e.get(fmt.Sprintf("/publication/%s", id), nil, pub); err != nil {
		return nil, err
	}
	return pub, nil
}

func (e *Engine) UpdatePublication(pub *models.Publication) (*models.Publication, error) {
	resPub := &models.Publication{}
	if _, err := e.put(fmt.Sprintf("/publication/%s", pub.ID), pub, resPub); err != nil {
		return nil, err
	}
	return resPub, nil
}

func (e *Engine) GetRelatedPublications(id string) ([]*models.RelatedPublication, error) {
	related := make([]*models.RelatedPublication, 0)
	if _, err := e.get(fmt.Sprintf("/publication/%s/related", id), nil, &related); err != nil {
		return nil, err
	}
	return related, nil
}

func (e *Engine) AddRelatedPublication(id, relatedID string) error {
	reqBody := struct {
		RelatedPublicationID string `json:"related_publication_id"`
	}{relatedID}
	_, err := e.post(fmt.Sprintf("/publication/%s/related", id), &reqBody, nil)
	return err
}

func (e *Engine) AddPublicationFile(id string, pubFile models.PublicationFile, b []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", pubFile.Filename)
	if err != nil {
		return err
	}
	part.Write(b)

	if err = writer.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/publication/%s/file", e.Config.URL, id), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(e.Config.Username, e.Config.Password)

	_, err = e.doRequest(req, nil)
	return err
}
