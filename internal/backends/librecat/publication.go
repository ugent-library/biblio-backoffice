package librecat

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/forms"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func (e *Client) Publications(args *models.SearchArgs) (*models.PublicationHits, error) {
	qp, err := forms.Encode(args)
	if err != nil {
		return nil, err
	}
	hits := &models.PublicationHits{}
	if _, err := e.get("/publication", qp, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Client) UserPublications(userID string, args *models.SearchArgs) (*models.PublicationHits, error) {
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

func (e *Client) GetPublication(id string) (*models.Publication, error) {
	pub := &models.Publication{}
	if _, err := e.get(fmt.Sprintf("/publication/%s", id), nil, pub); err != nil {
		return nil, err
	}
	return pub, nil
}

// quick and dirty reverse proxy
func (c *Client) ServePublicationFile(fileURL string, w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(fileURL)
	baseURL, _ := url.Parse(c.config.URL)
	proxy := httputil.NewSingleHostReverseProxy(baseURL)
	// update the headers to allow for SSL redirection
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.URL.Path = strings.Replace(u.Path, baseURL.Path, "", 1)
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Header.Del("Cookie")
	r.Host = u.Host
	r.SetBasicAuth(c.config.Username, c.config.Password)
	proxy.ServeHTTP(w, r)
}

func (e *Client) CreateUserPublication(userID, pubType string) (*models.Publication, error) {
	pub := &models.Publication{Type: pubType, Status: "private", CreationContext: "biblio-backend"}
	resPub := &models.Publication{}
	if _, err := e.post(fmt.Sprintf("/user/%s/publication", userID), pub, resPub); err != nil {
		return nil, err
	}
	return resPub, nil
}

func (e *Client) ImportUserPublicationByIdentifier(userID, source, identifier string) (*models.Publication, error) {
	reqData := struct {
		Source          string `json:"source"`
		Identifier      string `json:"identifier"`
		CreationContext string `json:"creation_context"`
	}{
		source,
		identifier,
		"biblio-backend",
	}
	pub := &models.Publication{}
	if _, err := e.post(fmt.Sprintf("/user/%s/publication/import", userID), &reqData, pub); err != nil {
		return nil, err
	}
	return pub, nil
}

func (e *Client) ImportUserPublications(userID, source string, file io.Reader) (string, error) {
	pipedReader, pipedWriter := io.Pipe()
	multiPartWriter := multipart.NewWriter(pipedWriter)

	go func() {
		var err error
		var part io.Writer

		defer func() {
			if err != nil {
				pipedWriter.CloseWithError(err)
			} else {
				pipedWriter.Close()
			}
		}()
		defer multiPartWriter.Close()

		if err := multiPartWriter.WriteField("source", source); err != nil {
			return
		}

		if err := multiPartWriter.WriteField("creation_context", "biblio-backend"); err != nil {
			return
		}

		part, err = multiPartWriter.CreateFormFile("file", "DUMMY")
		if err != nil {
			return
		}

		if _, err = io.Copy(part, file); err != nil {
			return
		}

	}()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%s/publication/import-from-file", e.config.URL, userID), pipedReader)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())
	req.SetBasicAuth(e.config.Username, e.config.Password)

	resData := struct {
		BatchID string `json:"batch_id"`
	}{}

	if _, err = e.doRequest(req, &resData); err != nil {
		return "", err
	}

	return resData.BatchID, nil
}

func (e *Client) UpdatePublication(pub *models.Publication) (*models.Publication, error) {
	resPub := &models.Publication{}
	if _, err := e.put(fmt.Sprintf("/publication/%s", pub.ID), pub, resPub); err != nil {
		return nil, err
	}
	return resPub, nil
}

func (e *Client) PublishPublication(pub *models.Publication) (*models.Publication, error) {
	oldStatus := pub.Status
	pub.Status = "public"
	updatedPub, err := e.UpdatePublication(pub)
	pub.Status = oldStatus
	return updatedPub, err
}

func (e *Client) BatchPublishPublications(userID string, args *models.SearchArgs) (err error) {
	var hits *models.PublicationHits
	for {
		hits, err = e.UserPublications(userID, args)
		for _, pub := range hits.Hits {
			pub.Status = "public"
			if _, err = e.UpdatePublication(pub); err != nil {
				break
			}
		}
		if !hits.NextPage {
			break
		}
		args.Page = args.Page + 1
	}
	return
}

func (e *Client) GetPublicationDatasets(id string) ([]*models.Dataset, error) {
	datasets := make([]*models.Dataset, 0)
	if _, err := e.get(fmt.Sprintf("/publication/%s/dataset", id), nil, &datasets); err != nil {
		return nil, err
	}
	return datasets, nil
}

func (e *Client) AddPublicationDataset(id, datasetID string) error {
	reqBody := struct {
		RelatedPublicationID string `json:"related_publication_id"`
	}{datasetID}
	_, err := e.post(fmt.Sprintf("/publication/%s/related", id), &reqBody, nil)
	return err
}

func (e *Client) RemovePublicationDataset(id, datasetID string) error {
	_, err := e.delete(fmt.Sprintf("/publication/%s/dataset/%s", id, datasetID), nil, nil)
	return err
}

func (e *Client) AddPublicationFile(id string, pubFile *models.PublicationFile, file io.Reader) error {
	pipedReader, pipedWriter := io.Pipe()
	multiPartWriter := multipart.NewWriter(pipedWriter)

	go func() {
		var err error
		var part io.Writer

		defer func() {
			if err != nil {
				pipedWriter.CloseWithError(err)
			} else {
				pipedWriter.Close()
			}
		}()
		defer multiPartWriter.Close()

		part, err = multiPartWriter.CreateFormFile("file", pubFile.Filename)
		if err != nil {
			return
		}

		if _, err = io.Copy(part, file); err != nil {
			return
		}

	}()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/publication/%s/file", e.config.URL, id), pipedReader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())
	req.SetBasicAuth(e.config.Username, e.config.Password)

	_, err = e.doRequest(req, nil)

	return err
}

func (e *Client) UpdatePublicationFile(id string, pubFile *models.PublicationFile) error {
	_, err := e.put(fmt.Sprintf("/publication/%s/file/%s", id, pubFile.ID), pubFile, nil)
	return err
}

func (e *Client) RemovePublicationFile(id, fileID string) error {
	_, err := e.delete(fmt.Sprintf("/publication/%s/file/%s", id, fileID), nil, nil)
	return err
}

func (e *Client) DeletePublication(id string) error {
	_, err := e.delete(fmt.Sprintf("/publication/%s", id), nil, nil)
	return err
}
