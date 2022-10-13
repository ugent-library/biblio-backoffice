package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type Config struct {
	URL      string
	Prefix   string
	Username string
	Password string
}

type Client struct {
	config Config
	http   *http.Client
}

/*
	copy from handle-server-api
*/
type HandleData struct {
	Url    string `json:"url"`
	Format string `json:"format"`
}

type HandleValue struct {
	Timestamp string      `json:"timestamp"`
	Type      string      `json:"type"`
	Index     int         `json:"index"`
	Ttl       int         `json:"ttl"`
	Data      *HandleData `json:"data"`
}

type Handle struct {
	Handle       string         `json:"handle"`
	ResponseCode int            `json:"responseCode"`
	Values       []*HandleValue `json:"values,omitempty"`
	Message      string         `json:"message,omitempty"`
}

func New(c Config) *Client {
	return &Client{
		config: c,
		http:   http.DefaultClient,
	}
}

type requestPayload struct {
	Data any `json:"data"`
}

type responsePayload struct {
	Data json.RawMessage `json:"data"`
}

func (c *Client) get(path string, qp url.Values, responseData any) (*http.Response, error) {
	req, err := c.newRequest("GET", path, qp, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.config.Username, c.config.Password)
	return c.doRequest(req, responseData)
}

func (c *Client) put(path string, qp url.Values, requestData any, responseData any) (*http.Response, error) {
	req, err := c.newRequest("PUT", path, qp, requestData)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.config.Username, c.config.Password)
	return c.doRequest(req, responseData)
}

func (c *Client) delete(path string, qp url.Values, responseData any) (*http.Response, error) {
	req, err := c.newRequest("DELETE", path, qp, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.config.Username, c.config.Password)
	return c.doRequest(req, responseData)
}

func (c *Client) newRequest(method, path string, vals url.Values, requestData any) (*http.Request, error) {
	var buf io.ReadWriter
	if requestData != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(&requestPayload{Data: requestData})
		if err != nil {
			return nil, err
		}
	}
	u := c.config.URL + path
	if vals != nil {
		u = u + "?" + vals.Encode()
	}
	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}
	if requestData != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) doRequest(req *http.Request, responseData any) (*http.Response, error) {
	res, err := c.http.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	var p responsePayload
	if err = json.NewDecoder(res.Body).Decode(&p); err != nil {
		return res, err
	}

	if responseData != nil {
		return res, json.Unmarshal(p.Data, responseData)
	}
	return res, nil
}

func (client *Client) GetByPublication(publication *models.Publication) (*Handle, error) {
	return client.Get(
		fmt.Sprintf("LU-%s", publication.ID),
	)
}

func (client *Client) Get(localId string) (*Handle, error) {
	h := &Handle{}
	_, err := client.get(
		fmt.Sprintf("/%s/%s", client.config.Prefix, localId),
		nil,
		h,
	)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (client *Client) UpsertByPublication(publication *models.Publication) (*Handle, error) {
	return client.Upsert(
		fmt.Sprintf("LU-%s", publication.ID),
	)
}

func (client *Client) Upsert(localId string) (*Handle, error) {
	return nil, nil
}

func (client *Client) DeleteByPublication(publication *models.Publication) (*Handle, error) {
	return client.Delete(
		fmt.Sprintf("LU-%s", publication.ID),
	)
}

func (client *Client) Delete(localId string) (*Handle, error) {
	return nil, nil
}
