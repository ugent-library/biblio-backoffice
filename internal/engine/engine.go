package engine

import (
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

func New(c Config) (*Engine, error) {
	e := &Engine{
		config: c,
		client: http.DefaultClient,
	}

	return e, nil
}
