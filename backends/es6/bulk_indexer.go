package es6

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	"github.com/ugent-library/biblio-backoffice/backends"
)

type bulkIndexer[T any] struct {
	bi         esutil.BulkIndexer
	docFn      func(T) (string, []byte, error)
	indexErrFn func(string, error)
}

func newBulkIndexer[T any](client *elasticsearch.Client, index string, docFn func(T) (string, []byte, error), config backends.BulkIndexerConfig) (*bulkIndexer[T], error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        client,
		Index:         index,
		FlushInterval: 1 * time.Second,
		Refresh:       "true",
		OnError: func(ctx context.Context, err error) {
			// TODO wrap error
			config.OnError(err)
		},
	})

	if err != nil {
		return nil, fmt.Errorf("es6.newBulkIndexer: %w", err)
	}

	return &bulkIndexer[T]{
		bi:         bi,
		docFn:      docFn,
		indexErrFn: config.OnIndexError,
	}, nil
}

func (b *bulkIndexer[T]) Index(ctx context.Context, t T) error {
	id, doc, err := b.docFn(t)
	if err != nil {
		return fmt.Errorf("bulkindexer.Index: %w", err)
	}

	err = b.bi.Add(
		ctx,
		esutil.BulkIndexerItem{
			Action:       "index",
			DocumentID:   id,
			DocumentType: "_doc",
			Body:         bytes.NewReader(doc),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				// TODO wrap error
				if err == nil {
					err = fmt.Errorf("%+v", res.Error)
				}
				b.indexErrFn(item.DocumentID, err)
			},
		},
	)

	if err != nil {
		return fmt.Errorf("bulkindexer.Index: %w", err)
	}
	return nil
}

func (b *bulkIndexer[T]) Close(ctx context.Context) error {
	if err := b.bi.Close(ctx); err != nil {
		return fmt.Errorf("bulkindexer.Close: %w", err)
	}
	return nil
}
