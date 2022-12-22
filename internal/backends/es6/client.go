package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	// "github.com/elastic/go-elasticsearch/v6/esutil"
)

type AliasStatus struct {
	Indexes []string
	Name    string
}

var PathNotFoundErr error = errors.New("path not found")
var AliasNotFoundErr error = errors.New("alias not found")
var IndexAtAliasErr error = errors.New("index at alias location")
var AliasAlreadyInitializedError error = errors.New("alias already initialized")

type Config struct {
	ClientConfig elasticsearch.Config
	Index        string
	Settings     string
}

type Client struct {
	Config
	es *elasticsearch.Client
}

type M map[string]any

func New(c Config) (*Client, error) {
	client, err := elasticsearch.NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{Config: c, es: client}, nil
}

func (c *Client) CreateIndex() error {
	return c.createIndex(c.Index, c.Settings)
}

func (c *Client) DeleteIndex() error {
	return c.deleteIndex(c.Index)
}

func (c *Client) createIndex(name string, settings string) error {
	r := strings.NewReader(settings)
	res, err := c.es.Indices.Create(name, c.es.Indices.Create.WithBody(r))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("unexpected es6 error: %s", res)
	}
	return nil
}

func (c *Client) hasIndex(name string) (bool, error) {
	res, err := esapi.IndicesExistsRequest{
		Index: []string{name},
	}.Do(
		context.Background(),
		c.es)
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

func (c *Client) hasAlias(name string) (bool, error) {
	res, err := esapi.IndicesExistsAliasRequest{
		Index: []string{name},
	}.Do(
		context.Background(),
		c.es)
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

func (c *Client) getAliasStatus() (*AliasStatus, error) {
	/*
		Index at alias location?
		TODO: find a better way to check if a target is a regular index
	*/
	{
		hasIndex, iErr := c.hasIndex(c.Index)
		if iErr != nil {
			return nil, iErr
		}
		if !hasIndex {
			return nil, PathNotFoundErr
		}
		hasAlias, aErr := c.hasAlias(c.Index)
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
			Name: []string{c.Index},
		}.Do(
			context.Background(),
			c.es)

		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == 404 {

			return nil, AliasNotFoundErr

		} else if res.StatusCode == 200 {

			var aliasGet map[string]interface{} = make(map[string]interface{})
			if e := json.NewDecoder(res.Body).Decode(&aliasGet); e != nil {
				return nil, e
			}

			indexNames := make([]string, 0)
			for k, _ := range aliasGet {
				indexNames = append(indexNames, k)
			}

			return &AliasStatus{
				Name:    c.Index,
				Indexes: indexNames,
			}, nil
		}

		return nil, fmt.Errorf("unexpected es6 error: %s", res)
	}
}

func (c *Client) copyIndex(from string, to string) error {
	var p map[string]interface{} = map[string]interface{}{
		"source": map[string]string{
			"index": from,
		},
		"dest": map[string]string{
			"index": to,
		},
	}
	payload, payloadErr := json.Marshal(p)
	if payloadErr != nil {
		return payloadErr
	}
	res, err := esapi.ReindexRequest{
		Body: bytes.NewReader(payload),
	}.Do(
		context.Background(),
		c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("unexpected es6 error: %s", res)
	}
	return nil
}

func (c *Client) deleteIndex(name string) error {
	res, err := c.es.Indices.Delete([]string{name})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("unexpected es6 error: %s", res)
	}
	return nil
}

func (c *Client) switchAlias(aliasName string, oldIndexName string, newIndexName string) error {
	var p map[string]interface{} = map[string]interface{}{
		"actions": []interface{}{
			map[string]interface{}{
				"add": map[string]string{
					"alias": aliasName,
					"index": newIndexName,
				},
			},
			map[string]interface{}{
				"remove": map[string]string{
					"alias": aliasName,
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
	}.Do(context.Background(), c.es)

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
func (c *Client) Init() error {

	_, aliasStatusErr := c.getAliasStatus()
	if aliasStatusErr == nil {
		return AliasAlreadyInitializedError
	} else if aliasStatusErr != AliasNotFoundErr && aliasStatusErr != PathNotFoundErr {
		return aliasStatusErr
	}

	now := time.Now().UTC()
	newIndexName := fmt.Sprintf(
		"%s_version_%s",
		c.Index,
		now.Format("20060102150405"),
	)

	if err := c.createIndex(newIndexName, c.Settings); err != nil {
		return err
	}

	var p map[string]interface{} = map[string]interface{}{
		"actions": []interface{}{
			map[string]interface{}{
				"add": map[string]string{
					"alias": c.Index,
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
	}.Do(context.Background(), c.es)

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
	1. create new index with timestamp
	2. copies existing index to new index
	3. switches alias from old index to nex index
	4. removes old index
*/
func (c *Client) Reindex() error {

	// new index name: <aliasName>_version_<yyyymmddHHMMSS>
	now := time.Now().UTC()
	newIndexName := fmt.Sprintf(
		"%s_version_%s",
		c.Index,
		now.Format("20060102150405"),
	)

	//index is alias or does not exist
	/*
		{ "alias": "aliasName", "indexes": ["idx1", "idx2"]}
	*/
	aliasStatus, aliasStatusErr := c.getAliasStatus()
	if aliasStatusErr != nil {
		return aliasStatusErr
	}

	if len(aliasStatus.Indexes) > 1 {
		return fmt.Errorf("alias %s points to more than one index", aliasStatus.Name)
	}

	oldIndexName := aliasStatus.Indexes[0]

	// create new index (without alias)
	if err := c.createIndex(newIndexName, c.Settings); err != nil {
		return err
	}

	// copy old index to new index
	if err := c.copyIndex(oldIndexName, newIndexName); err != nil {
		return err
	}

	// switch alias from old index to new index
	if err := c.switchAlias(c.Index, oldIndexName, newIndexName); err != nil {
		return err
	}

	// remove old index
	if err := c.deleteIndex(oldIndexName); err != nil {
		return err
	}

	return nil
}
