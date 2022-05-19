package es6

import (
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
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

type M map[string]interface{}

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
