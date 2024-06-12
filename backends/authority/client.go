package authority

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type M map[string]any

type Config struct {
	MongoDBURI string
	ESURI      []string
}

type Client struct {
	mongo *mongo.Client
	es    *elasticsearch.Client
}

func New(config Config) (*Client, error) {
	m, e := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(config.MongoDBURI))

	if e != nil {
		return nil, errors.Wrap(e, "unable to initialize connection to mongodb")
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: config.ESURI,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize connection to frontend elasticsearch")
	}

	return &Client{
		mongo: m,
		es:    es,
	}, nil
}

func (c *Client) search(index string, requestBody M, responseBody any) error {
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
