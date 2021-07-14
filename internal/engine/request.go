package engine

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type message struct {
	Data   json.RawMessage `json:"data,omitempty"`
	Errors []string        `json:"errors,omitempty"`
}

func (e *Engine) get(path string, qp url.Values, v interface{}) (*http.Response, error) {
	req, err := e.newRequest("GET", path, qp, nil)
	if err != nil {
		return nil, err
	}
	return e.doRequest(req, v)
}

func (e *Engine) newRequest(method, path string, vals url.Values, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	u := e.config.URL + path
	if vals != nil {
		u = u + "?" + vals.Encode()
	}

	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(e.config.Username, e.config.Password)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (e *Engine) doRequest(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := e.client.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()
	var msg message
	if err = json.NewDecoder(res.Body).Decode(&msg); err != nil {
		return res, err
	}
	return res, json.Unmarshal(msg.Data, v)
}
