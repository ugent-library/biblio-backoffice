package biblio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type requestPayload struct {
	Data any `json:"data"`
}

type responsePayload struct {
	Data json.RawMessage `json:"data"`
}

func (c *Client) get(path string, qp url.Values, responseData any) error {
	req, err := c.newRequest("GET", path, qp, nil)
	if err != nil {
		return fmt.Errorf("biblio frontend error: %w", err)
	}

	if strings.Contains(path, "/restricted/") {
		req.SetBasicAuth(c.config.Username, c.config.Password)
	}

	err = c.doRequest(req, responseData)
	if err != nil {
		return fmt.Errorf("biblio frontend error: %w", err)
	}

	return nil
}

func (c *Client) newRequest(method, path string, vals url.Values, requestData any) (*http.Request, error) {
	var buf io.ReadWriter
	if requestData != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(&requestPayload{Data: requestData})
		if err != nil {
			return nil, fmt.Errorf("json encoding error: %w", err)
		}
	}

	u := c.config.URL + path
	if vals != nil {
		u = u + "?" + vals.Encode()
	}

	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, fmt.Errorf("http request error: %w", err)
	}
	if requestData != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) doRequest(req *http.Request, responseData any) error {
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http client error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return fmt.Errorf("http response error: [status: %d]: %s", res.StatusCode, req.URL.String())
	}

	var p responsePayload
	if err = json.NewDecoder(res.Body).Decode(&p); err != nil {
		return fmt.Errorf("json decoding error: %w", err)
	}

	if responseData != nil {
		err = json.Unmarshal(p.Data, responseData)
		if err != nil {
			return fmt.Errorf("json decoding error: %w", err)
		}
	}
	return nil
}
