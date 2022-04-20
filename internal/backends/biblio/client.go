package biblio

import (
	"net/http"
)

type Config struct {
	URL string
}

type Client struct {
	config Config
	http   *http.Client
}

func New(c Config) *Client {
	return &Client{
		config: c,
		http:   http.DefaultClient,
	}
}
