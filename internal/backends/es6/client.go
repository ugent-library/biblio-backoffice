package es6

import (
	"bytes"
	"io"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
)

type Config struct {
	ClientConfig   elasticsearch.Config
	Index          string
	Settings       string
	IndexRetention int // -1: keep all old indexes, >=0: keep x old indexes
}

type Client struct {
	Config
	es *elasticsearch.Client
}

func New(c Config) (*Client, error) {
	client, err := elasticsearch.NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{Config: c, es: client}, nil
}

func (c *Client) SearchWithOpts(opts []func(*esapi.SearchRequest), fn func(r io.ReadCloser) error) error {
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

	return fn(res.Body)
}
