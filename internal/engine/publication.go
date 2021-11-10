package engine

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sync"

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

// TODO change validation for new and draft publications in librecat
func (e *Engine) CreatePublication(pt string) (*models.Publication, error) {
	pub := &models.Publication{Type: pt, Status: "private", Title: "New publication"}
	resPub := &models.Publication{}
	if _, err := e.post("/publication", pub, resPub); err != nil {
		return nil, err
	}
	return resPub, nil
}

// TODO: set constraint to not research_data
func (e *Engine) ImportUserPublicationByIdentifier(userID, identifier, identifier_type string) (*models.Publication, error) {
	reqData := struct {
		Identifier     string `json:"identifier"`
		IdentifierType string `json:"identifier_type"`
	}{
		identifier,
		identifier_type,
	}
	pub := &models.Publication{}
	if _, err := e.post(fmt.Sprintf("/user/%s/publication/import", userID), &reqData, pub); err != nil {
		return nil, err
	}
	return pub, nil
}

func (e *Engine) ImportUserPublications(userID, source string, file io.Reader) (string, error) {
	b := &bytes.Buffer{}
	multiPartWriter := multipart.NewWriter(b)

	if err := multiPartWriter.WriteField("source", source); err != nil {
		return "", err
	}

	part, err := multiPartWriter.CreateFormFile("file", "DUMMY")
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	multiPartWriter.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%s/publication/import-from-file", e.Config.URL, userID), b)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())
	req.SetBasicAuth(e.Config.Username, e.Config.Password)

	if _, err = e.doRequest(req, nil); err != nil {
		return "", err
	}

	// IMPORTANT: wait for goroutine to finish
	// wg.Wait()

	return "", err

}

// func (e *Engine) ImportUserPublications(userID, source string, file io.Reader) (string, error) {
// 	pr, pw := io.Pipe()

// 	req, err := http.NewRequest(http.MethodPost, y, pr)
// 	if err != nil {
// 		return "", err
// 	}

// 	mw := multipart.NewWriter(pw)
// 	req.Header.Add("Content-Type", mw.FormDataContentType())
// 	req.SetBasicAuth(e.Config.Username, e.Config.Password)

// 	var wg sync.WaitGroup

// 	wg.Add(1)

// 	go func() {
// 		var err error

// 		defer func() {
// 			if err != nil {
// 				pw.CloseWithError(err)
// 			} else {
// 				pw.Close()
// 			}
// 		}()
// 		defer mw.Close()
// 		defer wg.Done()

// 		part, err := mw.CreateFormFile("file", "DUMMY")
// 		if err != nil {
// 			return
// 		}

// 		if _, err = io.Copy(part, file); err != nil {
// 			return
// 		}

// 		// if err = mw.WriteField("source", source); err != nil {
// 		// 	return
// 		// }
// 	}()

// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// pr, pw := io.Pipe()
// 	// defer pr.Close()
// 	// mw := multipart.NewWriter(pw)
// 	// errChan := make(chan error, 1)
// 	// go func() {

// 	// }()

// 	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%s/publication/import-from-file", e.Config.URL, userID), pr)
// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// req.Header.Add("Content-Type", mw.FormDataContentType())
// 	// req.SetBasicAuth(e.Config.Username, e.Config.Password)

// 	resData := &struct {
// 		BatchID string `json:"batch_id"`
// 	}{}

// 	if _, err = e.doRequest(req, resData); err != nil {
// 		return "", err
// 	}

// 	// IMPORTANT: wait for goroutine to finish
// 	wg.Wait()

// 	return resData.BatchID, nil
// }

func (e *Engine) UpdatePublication(pub *models.Publication) (*models.Publication, error) {
	resPub := &models.Publication{}
	if _, err := e.put(fmt.Sprintf("/publication/%s", pub.ID), pub, resPub); err != nil {
		return nil, err
	}
	return resPub, nil
}

func (e *Engine) PublishPublication(pub *models.Publication) (*models.Publication, error) {
	pub.Status = "public"
	return e.UpdatePublication(pub)
}

func (e *Engine) GetPublicationDatasets(id string) ([]*models.Dataset, error) {
	datasets := make([]*models.Dataset, 0)
	if _, err := e.get(fmt.Sprintf("/publication/%s/dataset", id), nil, &datasets); err != nil {
		return nil, err
	}
	return datasets, nil
}

func (e *Engine) AddPublicationDataset(id, datasetID string) error {
	reqBody := struct {
		RelatedPublicationID string `json:"related_publication_id"`
	}{datasetID}
	_, err := e.post(fmt.Sprintf("/publication/%s/related", id), &reqBody, nil)
	return err
}

func (e *Engine) RemovePublicationDataset(id, datasetID string) error {
	_, err := e.delete(fmt.Sprintf("/publication/%s/dataset/%s", id, datasetID), nil, nil)
	return err
}

func (e *Engine) AddPublicationFile(id string, pubFile models.PublicationFile, file io.Reader) error {
	pipedReader, pipedWriter := io.Pipe()
	multiPartWriter := multipart.NewWriter(pipedWriter)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {

		defer wg.Done()
		defer pipedWriter.Close()
		defer multiPartWriter.Close()

		part, err := multiPartWriter.CreateFormFile("file", pubFile.Filename)
		if err != nil {
			return
		}

		if _, err = io.Copy(part, file); err != nil {
			return
		}

	}()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/publication/%s/file", e.Config.URL, id), pipedReader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())
	req.SetBasicAuth(e.Config.Username, e.Config.Password)

	if _, err = e.doRequest(req, nil); err != nil {
		return err
	}

	// IMPORTANT: wait for goroutine to finish
	wg.Wait()

	return err
}

func (e *Engine) UpdatePublicationFile(id string, pubFile *models.PublicationFile) error {
	_, err := e.put(fmt.Sprintf("/publication/%s/file/%s", id, pubFile.ID), pubFile, nil)
	return err
}

func (e *Engine) RemovePublicationFile(id, fileID string) error {
	_, err := e.delete(fmt.Sprintf("/publication/%s/file/%s", id, fileID), nil, nil)
	return err
}
