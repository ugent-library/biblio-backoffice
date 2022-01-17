package spdx

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type license struct {
	LicenceID string `json:"licenseId"`
	Name      string `json:"name"`
}

type licenseEnv struct {
	Licenses []license `json:"licenses"`
}

type localSuggester struct {
	index bleve.Index
}

func New() *localSuggester {
	l := &licenseEnv{}
	file, _ := ioutil.ReadFile("etc/spdx-licenses.json")
	_ = json.Unmarshal([]byte(file), l)

	indexMapping := bleve.NewIndexMapping()
	docMapping := bleve.NewDocumentMapping()
	textFieldMapping := bleve.NewTextFieldMapping()
	docMapping.AddFieldMappingsAt("licenseId", textFieldMapping)
	docMapping.AddFieldMappingsAt("name", textFieldMapping)
	indexMapping.AddDocumentMapping("_default", docMapping)
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		log.Fatal(err)
	}

	for _, license := range l.Licenses {
		if err := index.Index(license.LicenceID, license); err != nil {
			log.Fatal(err)
		}
	}

	return &localSuggester{index: index}
}

// simulate matchphraseprefix query
// https://github.com/blevesearch/bleve/issues/377
func (c *localSuggester) SuggestLicenses(q string) ([]models.Completion, error) {
	if q == "" {
		return nil, nil
	}

	var searchResults *bleve.SearchResult
	var err error
	words := strings.Fields(q)

	if len(words) == 1 {
		bq := bleve.NewDisjunctionQuery(
			bleve.NewPrefixQuery(q),
			bleve.NewMatchQuery(q),
		)
		search := bleve.NewSearchRequest(bq)
		search.Fields = []string{"name"}
		searchResults, err = c.index.Search(search)
		if err != nil {
			return nil, err
		}
	} else {
		phrase := ""
		k := 0
		for k != len(words)-1 {
			phrase += words[k] + " "
			k++
		}

		phrase = phrase[0 : len(phrase)-1]
		prefix := words[len(words)-1]

		bq := bleve.NewConjunctionQuery(
			bleve.NewMatchPhraseQuery(phrase),
			bleve.NewPrefixQuery(prefix),
		)
		search := bleve.NewSearchRequest(bq)
		search.Fields = []string{"name"}
		searchResults, err = c.index.Search(search)
		if err != nil {
			return nil, err
		}

		if len(searchResults.Hits) == 0 {
			search := bleve.NewSearchRequest(bleve.NewMatchQuery(prefix))
			search.Fields = []string{"name"}
			searchResults, err = c.index.Search(search)
			if err != nil {
				return nil, err
			}
		}
	}

	hits := make([]models.Completion, 0, len(searchResults.Hits))
	for _, hit := range searchResults.Hits {
		hits = append(hits, models.Completion{
			ID:          hit.ID,
			Heading:     hit.ID,
			Description: hit.Fields["name"].(string),
		})
	}
	return hits, nil
}
