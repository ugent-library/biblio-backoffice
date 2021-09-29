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
	Config Config
	client *http.Client
}

func New(c Config) (*Engine, error) {
	e := &Engine{
		Config: c,
		client: http.DefaultClient,
	}

	return e, nil
}
