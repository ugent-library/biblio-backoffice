package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/ugent-library/biblio-backoffice/backends"
)

type indexSwitcher[T any] struct {
	client    *elasticsearch.Client
	alias     string
	index     string
	retention int
	bi        *bulkIndexer[T]
}

func newIndexSwitcher[T any](client *elasticsearch.Client, alias, settings string, retention int, docFn func(T) (string, []byte, error), config backends.BulkIndexerConfig) (*indexSwitcher[T], error) {
	// generate new index name
	index := fmt.Sprintf("%s_%s", alias, time.Now().UTC().Format("20060102150405"))

	// create new index
	body := strings.NewReader(settings)
	res, err := client.Indices.Create(index, client.Indices.Create.WithBody(body))
	if err != nil {
		return nil, fmt.Errorf("es6.newIndexSwitcher: failed to create es6 index %s: %w", index, err)
	}
	if res.IsError() {
		// TODO read res body
		return nil, fmt.Errorf("es6.newIndexSwitcher: failed to create es6 index %s: %+v", index, res)
	}

	// TODO the default settings of the bulk indexer are not very efficient for this use case
	bi, err := newBulkIndexer(client, index, docFn, config)
	if err != nil {
		return nil, fmt.Errorf("es6.newIndexSwitcher: failed to create new bulk indexer %s: %w", index, err)
	}

	return &indexSwitcher[T]{
		client:    client,
		alias:     alias,
		index:     index,
		retention: retention,
		bi:        bi,
	}, nil
}

func (is *indexSwitcher[T]) Index(ctx context.Context, t T) error {
	return is.bi.Index(ctx, t)
}

func (is *indexSwitcher[T]) Switch(ctx context.Context) error {
	if err := is.bi.Close(ctx); err != nil {
		return fmt.Errorf("indexswitcher.Switch: failed to close bulk indexer for index %s: %w", is.index, err)
	}

	actions := []map[string]any{
		{
			"add": map[string]string{
				"alias": is.alias,
				"index": is.index,
			},
		},
	}

	oldIndexes, err := is.oldIndexes(ctx)
	if err != nil {
		return fmt.Errorf("indexswitcher.Switch: %w", err)
	}

	for i, idx := range oldIndexes {
		if is.retention < 0 || i >= len(oldIndexes)-is.retention {
			actions = append(actions, map[string]any{
				"remove": map[string]string{
					"alias": is.alias,
					"index": idx,
				},
			})
		} else {
			actions = append(actions, map[string]any{
				"remove_index": map[string]string{
					"index": idx,
				},
			})
		}
	}

	body, err := json.Marshal(map[string]any{"actions": actions})
	if err != nil {
		return fmt.Errorf("indexswitcher.Switch: %w", err)
	}
	req := esapi.IndicesUpdateAliasesRequest{Body: bytes.NewReader(body)}
	res, err := req.Do(ctx, is.client)
	if err != nil {
		return fmt.Errorf("indexswitcher.Switch: es6 http error: %w", err)
	}
	if res.IsError() {
		// TODO read res body
		return fmt.Errorf("indexswitcher.Switch: es6 error: %+v", res)
	}

	return nil
}

func (is *indexSwitcher[T]) oldIndexes(ctx context.Context) ([]string, error) {
	req := esapi.CatIndicesRequest{
		Format: "json",
	}
	res, err := req.Do(ctx, is.client)
	if err != nil {
		return nil, fmt.Errorf("indexswitcher.oldIndexes: es6 http error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		// TODO read res body
		return nil, fmt.Errorf("indexswitcher.oldIndexes: es6 error: %+v", res)
	}

	indexes := []struct{ Index string }{}
	if err := json.NewDecoder(res.Body).Decode(&indexes); err != nil {
		return nil, fmt.Errorf("indexswitcher.oldIndexes: %w", err)
	}

	r := regexp.MustCompile(`^` + is.alias + `_[0-9]+$`)

	var oldIndexes []string
	for _, idx := range indexes {
		if r.MatchString(idx.Index) && idx.Index != is.index {
			oldIndexes = append(oldIndexes, idx.Index)
		}
	}

	sort.Strings(oldIndexes)

	return oldIndexes, nil
}
