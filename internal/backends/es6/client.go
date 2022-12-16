package es6

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	// "github.com/elastic/go-elasticsearch/v6/esutil"
)

type Config struct {
	ClientConfig elasticsearch.Config
	Index        string
	Settings     string
}

type Client struct {
	Config
	es *elasticsearch.Client
}

type M map[string]any

func New(c Config) (*Client, error) {
	client, err := elasticsearch.NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{Config: c, es: client}, nil
}

func (c *Client) CreateIndex() error {
	r := strings.NewReader(c.Settings)
	res, err := c.es.Indices.Create(c.Index, c.es.Indices.Create.WithBody(r))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (c *Client) DeleteIndex() error {
	res, err := c.es.Indices.Delete([]string{c.Index})
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (c *Client) searchWithOpts(opts []func(*esapi.SearchRequest), responseBody any) error {

	res, err := c.es.Search(opts...)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, res.Body); err != nil {
			return err
		}
		return errors.New("Es6 error response: " + buf.String())
	}

	if err := json.NewDecoder(res.Body).Decode(responseBody); err != nil {
		return errors.Wrap(err, "Error parsing the response body")
	}

	return nil
}
