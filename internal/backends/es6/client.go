package es6

import (
	"github.com/elastic/go-elasticsearch/v6"
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
