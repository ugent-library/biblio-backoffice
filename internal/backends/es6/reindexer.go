package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
)

type reindexer struct {
	Client
	alias string
}

type AliasStatus struct {
	Indexes []string
	Name    string
}

var PathNotFoundErr error = errors.New("path not found")
var AliasNotFoundErr error = errors.New("alias not found")
var IndexAtAliasErr error = errors.New("index at alias location")
var AliasAlreadyInitializedError error = errors.New("alias already initialized")

func (r *reindexer) hasIndex(name string) (bool, error) {
	res, err := esapi.IndicesExistsRequest{
		Index: []string{name},
	}.Do(
		context.Background(),
		r.Client.es)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return false, nil
	}

	if res.IsError() {
		return false, fmt.Errorf("unexpected es6 error: %s", res)
	}

	return res.StatusCode == 200, nil
}

func (r *reindexer) getAliasStatus(alias string) (*AliasStatus, error) {
	/*
		Index at alias location?
		TODO: find a better way to check if a target is a regular index
	*/
	{
		hasIndex, iErr := r.hasIndex(alias)
		if iErr != nil {
			return nil, iErr
		}
		if !hasIndex {
			return nil, PathNotFoundErr
		}
		hasAlias, aErr := r.hasAlias(alias)
		if aErr != nil {
			return nil, aErr
		}
		if hasIndex && !hasAlias {
			return nil, IndexAtAliasErr
		}
	}

	// Find out what the alias is pointing to
	{
		res, err := esapi.IndicesGetAliasRequest{
			Name: []string{alias},
		}.Do(
			context.Background(),
			r.Client.es)

		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == 404 {

			return nil, AliasNotFoundErr

		} else if res.StatusCode != 200 {
			return nil, fmt.Errorf("unexpected es6 error: %s", res)
		}

		var aliasGet map[string]interface{} = make(map[string]interface{})
		if e := json.NewDecoder(res.Body).Decode(&aliasGet); e != nil {
			return nil, e
		}

		indexNames := make([]string, 0)
		for k, _ := range aliasGet {
			indexNames = append(indexNames, k)
		}

		return &AliasStatus{
			Name:    alias,
			Indexes: indexNames,
		}, nil

	}
}

func (r *reindexer) hasAlias(alias string) (bool, error) {
	res, err := esapi.IndicesExistsAliasRequest{
		Index: []string{alias},
	}.Do(
		context.Background(),
		r.Client.es)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return false, nil
	}

	if res.IsError() {
		return false, fmt.Errorf("unexpected es6 error: %s", res)
	}

	return res.StatusCode == 200, nil
}

func (r *reindexer) switchAlias(alias string, oldIndexName string, newIndexName string) error {
	var p map[string]interface{} = map[string]interface{}{
		"actions": []interface{}{
			map[string]interface{}{
				"add": map[string]string{
					"alias": alias,
					"index": newIndexName,
				},
			},
			map[string]interface{}{
				"remove": map[string]string{
					"alias": alias,
					"index": oldIndexName,
				},
			},
		},
	}
	payload, payloadErr := json.Marshal(p)
	if payloadErr != nil {
		return payloadErr
	}
	res, resErr := esapi.IndicesUpdateAliasesRequest{
		Body: bytes.NewReader(payload),
	}.Do(context.Background(), r.Client.es)

	if resErr != nil {
		return resErr
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}

	return nil
}

/*
	Initializes index:

	1. creates new index with timestamp
	2. sets alias to new index

	Returns error when:

	* Something (index or alias) already exists at the alias location
	* unexpected es6 error occurs
*/
func (r *reindexer) InitAlias() error {
	_, aliasStatusErr := r.getAliasStatus(r.alias)
	if aliasStatusErr == nil {
		return AliasAlreadyInitializedError
	} else if aliasStatusErr != AliasNotFoundErr && aliasStatusErr != PathNotFoundErr {
		return aliasStatusErr
	}

	now := time.Now().UTC()
	newIndexName := fmt.Sprintf(
		"%s_version_%s",
		r.alias,
		now.Format("20060102150405"),
	)

	if err := r.Client.createIndex(newIndexName, r.Client.Settings); err != nil {
		return err
	}

	var p map[string]interface{} = map[string]interface{}{
		"actions": []interface{}{
			map[string]interface{}{
				"add": map[string]string{
					"alias": r.alias,
					"index": newIndexName,
				},
			},
		},
	}
	payload, payloadErr := json.Marshal(p)
	if payloadErr != nil {
		return payloadErr
	}
	res, err := esapi.IndicesUpdateAliasesRequest{
		Body: bytes.NewReader(payload),
	}.Do(context.Background(), r.Client.es)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("unexpected es6 error: %s", res)
	}

	return nil
}

/*
	Returns all indexes that have the current alias as prefix,
	and shows if they are currently linked to that alias

	Technicall returns array of objects
	each object contains keys "index" and "active":

		object["index"] is the name of an index that is version of the current alias
		object["active"] contains either "true" or "false"
*/
func (r *reindexer) ListIndexes() ([]map[string]string, error) {

	aliasStatus, aliasStatusErr := r.getAliasStatus(r.alias)
	if aliasStatusErr != nil {
		return nil, aliasStatusErr
	}

	if len(aliasStatus.Indexes) > 1 {
		return nil, fmt.Errorf(
			"alias %s more than index",
			aliasStatus.Name,
		)
	}

	res, resErr := esapi.CatIndicesRequest{
		Format: "json",
	}.Do(context.Background(), r.Client.es)

	if resErr != nil {
		return nil, resErr
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error: %s", res)
	}

	allIndexes := make([]map[string]string, 0)
	if e := json.NewDecoder(res.Body).Decode(&allIndexes); e != nil {
		return nil, e
	}

	liveIndex := aliasStatus.Indexes[0]
	prefix := r.alias + "_version_"

	indexes := make([]map[string]string, 0, len(allIndexes))

	for _, idx := range allIndexes {
		if !strings.HasPrefix(idx["index"], prefix) {
			continue
		}
		var isLive string = "false"
		if idx["index"] == liveIndex {
			isLive = "true"
		}
		var m map[string]string = map[string]string{
			"index":  idx["index"],
			"active": isLive,
		}
		indexes = append(indexes, m)
	}

	return indexes, nil
}

func (r *reindexer) RemoveOldIndexes(keep int) error {

	indexes, e := r.ListIndexes()
	if e != nil {
		return e
	}

	oldIndexes := make([]string, 0, len(indexes))
	for _, idx := range indexes {
		if idx["active"] == "false" {
			oldIndexes = append(oldIndexes, idx["index"])
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(oldIndexes)))

	nLeft := len(oldIndexes)
	for _, idx := range oldIndexes {
		if nLeft <= keep {
			break
		}
		if e := r.deleteIndex(idx); e != nil {
			return e
		}
		nLeft--
		fmt.Fprintf(os.Stderr, "removed idx %s\n", idx)
	}

	return nil
}
