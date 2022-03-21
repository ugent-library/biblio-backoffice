package fc6

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (c *Client) GetDataset(id string) (*models.Dataset, error) {
	log.Printf("get dataset %s", id)

	url := c.config.URL + "/biblio-objects/" + id + "/metadata.json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	req.SetBasicAuth("fedoraAdmin", "fedoraAdmin")
	res, err := c.http.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer res.Body.Close()

	var d = models.Dataset{}
	if err = json.NewDecoder(res.Body).Decode(&d); err != nil {
		log.Print(err)
		return nil, err
	}

	return &d, nil
}

func (c *Client) GetDatasetPublications(id string) ([]*models.Publication, error) {
	return nil, nil
}

func (c *Client) CreateDataset(d *models.Dataset) (*models.Dataset, error) {
	url := c.config.URL + "/biblio-objects/" + d.ID
	body := strings.NewReader(`
		@prefix pcdm: <http://pcdm.org/models#>
		<> a pcdm:Object .
		`)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("fedoraAdmin", "fedoraAdmin")
	req.Header.Set("Content-Type", "text/turtle")

	_, err = c.http.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if _, err = c.UpdateDataset(d); err != nil {
		log.Print(err)
		return nil, err
	}

	if err = c.markResourceAsFile("/biblio-objects/" + d.ID + "/metadata.json"); err != nil {
		return nil, err
	}

	return d, nil
}

func (c *Client) UpdateDataset(d *models.Dataset) (*models.Dataset, error) {
	now := time.Now()
	if d.DateCreated == nil {
		d.DateCreated = &now
	}
	d.DateUpdated = &now

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(d); err != nil {
		return nil, err
	}

	url := c.config.URL + "/biblio-objects/" + d.ID + "/metadata.json"
	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("fedoraAdmin", "fedoraAdmin")
	req.Header.Set("Content-Type", "application/json")

	_, err = c.http.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return d, nil
}

func (c *Client) DeleteDataset(id string) error {
	url := c.config.URL + "/biblio-objects/" + id
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("fedoraAdmin", "fedoraAdmin")
	req.Header.Set("Content-Type", "text/turtle")

	_, err = c.http.Do(req)
	return err
}
