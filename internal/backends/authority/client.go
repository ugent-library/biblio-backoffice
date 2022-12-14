package authority

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI string
}

type Client struct {
	mongo *mongo.Client
}

func New(config Config) (*Client, error) {
	m, e := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(config.MongoURI))

	if e != nil {
		return nil, errors.Wrap(e, "unable to initialize connection to")
	}

	return &Client{
		mongo: m,
	}, nil
}
