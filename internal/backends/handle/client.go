package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type Config struct {
	BaseURL         string
	FrontEndBaseURL string
	Prefix          string
	Username        string
	Password        string
}

type Client struct {
	config Config
	http   *http.Client
}

func NewClient(c Config) *Client {
	return &Client{
		config: c,
		http:   http.DefaultClient,
	}
}

// func (c *Client) get(path string, qp url.Values, responseData any) (*http.Response, error) {
// 	req, err := c.newRequest("GET", path, qp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.SetBasicAuth(c.config.Username, c.config.Password)
// 	return c.doRequest(req, responseData)
// }

func (c *Client) put(path string, qp url.Values, responseData any) (*http.Response, error) {
	req, err := c.newRequest("PUT", path, qp)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.config.Username, c.config.Password)
	return c.doRequest(req, responseData)
}

// func (c *Client) delete(path string, qp url.Values, responseData any) (*http.Response, error) {
// 	req, err := c.newRequest("DELETE", path, qp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.SetBasicAuth(c.config.Username, c.config.Password)
// 	return c.doRequest(req, responseData)
// }

func (c *Client) newRequest(method, path string, vals url.Values) (*http.Request, error) {
	url := c.config.BaseURL + path
	if vals != nil {
		url = url + "?" + vals.Encode()
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) doRequest(req *http.Request, responseData any) (*http.Response, error) {
	res, err := c.http.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(responseData); err != nil {
		return res, err
	}

	return res, nil
}

// func (client *Client) GetHandle(localId string) (*models.Handle, error) {
// 	h := &models.Handle{}
// 	_, err := client.get(
// 		fmt.Sprintf("/%s/LU-%s", client.config.Prefix, localId),
// 		nil,
// 		h,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return h, nil
// }

func (client *Client) UpsertHandle(localId string) (*models.Handle, error) {
	h := &models.Handle{}
	qp := url.Values{}
	qp.Add("url", fmt.Sprintf("%s/%s", client.config.FrontEndBaseURL, localId))
	_, err := client.put(
		fmt.Sprintf("/%s/LU-%s", client.config.Prefix, localId),
		qp,
		h,
	)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// func (client *Client) DeleteHandle(localId string) (*models.Handle, error) {
// 	h := &models.Handle{}
// 	_, err := client.delete(
// 		fmt.Sprintf("/%s/LU-%s", client.config.Prefix, localId),
// 		nil,
// 		h,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return h, nil
// }
