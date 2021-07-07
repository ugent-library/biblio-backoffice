package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	URL      string
	Username string
	Password string
}

type Engine struct {
	config Config
	client *http.Client
}

type message struct {
	Data   json.RawMessage `json:"data,omitempty"`
	Errors []string        `json:"errors,omitempty"`
}

func New(c Config) (*Engine, error) {
	e := &Engine{
		config: c,
		client: http.DefaultClient,
	}

	return e, nil
}

func (e *Engine) UserPublications(userID string) (*PublicationHits, error) {
	path := fmt.Sprintf("/users/%s/publications", userID)
	hits := &PublicationHits{}
	if _, err := e.get(path, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Engine) UserPublicationsContributed(userID string) (*PublicationHits, error) {
	path := fmt.Sprintf("/users/%s/publications-contributed", userID)
	hits := &PublicationHits{}
	if _, err := e.get(path, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Engine) UserPublicationsCreated(userID string) (*PublicationHits, error) {
	path := fmt.Sprintf("/users/%s/publications-created", userID)
	hits := &PublicationHits{}
	if _, err := e.get(path, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Engine) get(path string, v interface{}) (*http.Response, error) {
	req, err := e.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return e.doRequest(req, v)
}

func (e *Engine) newRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, e.config.URL+path, buf)
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
