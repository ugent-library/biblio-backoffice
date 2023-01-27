package authority

import (
	"context"

	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backoffice/internal/backends/es6"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoDBURI string
	ES6Config  es6.Config
}

type Client struct {
	mongo *mongo.Client
	es    *es6.Client
}

func New(config Config) (*Client, error) {
	m, e := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(config.MongoDBURI))

	if e != nil {
		return nil, errors.Wrap(e, "unable to initialize connection to mongodb")
	}

	es, esErr := es6.New(config.ES6Config)

	if esErr != nil {
		return nil, errors.Wrap(esErr, "unable to initialize connection to frontend elasticsearch")
	}

	return &Client{
		mongo: m,
		es:    es,
	}, nil
}
