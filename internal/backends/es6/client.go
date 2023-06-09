package es6

import (
	"github.com/elastic/go-elasticsearch/v6"
)

type Config struct {
	ClientConfig elasticsearch.Config
}

type Client struct {
	es *elasticsearch.Client
}

func New(c Config) (*Client, error) {
	client, err := elasticsearch.NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{es: client}, nil
}
