package es6

import (
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	// "github.com/elastic/go-elasticsearch/v6/esutil"
)

type Config struct {
	ClientConfig    elasticsearch.Config
	DatasetIndex    string
	DatasetSettings string
}

type Client struct {
	Config
	es *elasticsearch.Client
}

type M map[string]interface{}

func New(c Config) (*Client, error) {
	client, err := elasticsearch.NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{c, client}, nil
}

func (c *Client) CreateDatasetIndex() error {
	r := strings.NewReader(c.DatasetSettings)
	res, err := c.es.Indices.Create(c.DatasetIndex, c.es.Indices.Create.WithBody(r))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (c *Client) DeleteDatasetIndex() error {
	res, err := c.es.Indices.Delete([]string{c.DatasetIndex})
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}
