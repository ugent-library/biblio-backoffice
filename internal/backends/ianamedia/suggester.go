package ianamedia

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type mediaType struct {
	Extensions []string `json:"extensions"`
}

type ianaEnv map[string]mediaType

type localSuggester struct {
	index bleve.Index
}

func New() *localSuggester {
	indexMapping := bleve.NewIndexMapping()
	docMapping := bleve.NewDocumentMapping()
	textFieldMapping := bleve.NewTextFieldMapping()
	docMapping.AddFieldMappingsAt("mediaType", textFieldMapping)
	docMapping.AddFieldMappingsAt("extensions", textFieldMapping)
	indexMapping.AddDocumentMapping("_default", docMapping)
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		log.Fatal(err)
	}

	return &localSuggester{index: index}
}

func (s *localSuggester) IndexAll() error {
	env := make(ianaEnv)

	file, err := ioutil.ReadFile("etc/iana-media-types.json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(file), &env); err != nil {
		return err
	}

	for k, mt := range env {
		doc := struct {
			MediaType  string   `json:"mediaType"`
			Extensions []string `json:"extensions"`
		}{
			k,
			mt.Extensions,
		}
		if err := s.index.Index(k, doc); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

// simulate matchphraseprefix query
// https://github.com/blevesearch/bleve/issues/377
func (s *localSuggester) SuggestMediaTypes(q string) ([]models.Completion, error) {
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
		search.Fields = []string{"extensions"}
		searchResults, err = s.index.Search(search)
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
		search.Fields = []string{"extensions"}
		searchResults, err = s.index.Search(search)
		if err != nil {
			return nil, err
		}

		if len(searchResults.Hits) == 0 {
			search := bleve.NewSearchRequest(bleve.NewMatchQuery(prefix))
			search.Fields = []string{"extensions"}
			searchResults, err = s.index.Search(search)
			if err != nil {
				return nil, err
			}
		}
	}
	// log.Printf("%+v", searchResults)
	hits := make([]models.Completion, 0, len(searchResults.Hits))
	for _, hit := range searchResults.Hits {
		log.Printf("%+v", hit.Fields)
		desc := ""
		if ext, ok := hit.Fields["extensions"].(string); ok {
			desc = "." + ext
		}
		hits = append(hits, models.Completion{
			ID:          hit.ID,
			Heading:     hit.ID,
			Description: desc,
		})
	}
	return hits, nil
}
