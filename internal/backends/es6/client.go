package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

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

type M map[string]any

func New(c Config) (*Client, error) {
	client, err := elasticsearch.NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{Config: c, es: client}, nil
}

func (c *Client) CreateIndex() error {
	return c.createIndex(c.Index, c.Settings)
}

func (c *Client) DeleteIndex() error {
	return c.deleteIndex(c.Index)
}

func (c *Client) createIndex(name string, settings string) error {
	r := strings.NewReader(settings)
	res, err := c.es.Indices.Create(name, c.es.Indices.Create.WithBody(r))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("unexpected es6 error: %s", res)
	}
	return nil
}

func (c *Client) deleteIndex(name string) error {
	res, err := c.es.Indices.Delete([]string{name})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("unexpected es6 error: %s", res)
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

func (c *Client) SearchWithBody(index string, requestBody M, responseBody any) error {

	sorts := make([]string, 0)

	if v, exists := requestBody["sort"]; exists {
		sorts = append(sorts, v.([]string)...)
		delete(requestBody, "sort")
	}

	opts := []func(*esapi.SearchRequest){
		c.es.Search.WithContext(context.Background()),
		c.es.Search.WithIndex(index),
		c.es.Search.WithTrackTotalHits(true),
		c.es.Search.WithSort(sorts...),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(requestBody); err != nil {
		return err
	}
	opts = append(opts, c.es.Search.WithBody(&buf))

	return c.searchWithOpts(opts, responseBody)
}
