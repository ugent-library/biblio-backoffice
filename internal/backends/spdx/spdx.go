package spdx

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
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
	file, _ := ioutil.ReadFile("etc/licenses.json")
	_ = json.Unmarshal([]byte(file), l)

	mapping := bleve.NewIndexMapping()
	licenseMapping := bleve.NewDocumentMapping()
	idFieldMapping := bleve.NewTextFieldMapping()
	idFieldMapping.Analyzer = keyword.Name
	licenseMapping.AddFieldMappingsAt("ID", idFieldMapping)
	nameFieldMapping := bleve.NewTextFieldMapping()
	nameFieldMapping.Analyzer = "en"
	licenseMapping.AddFieldMappingsAt("Name", nameFieldMapping)
	mapping.AddDocumentMapping("license", licenseMapping)
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		log.Fatal(err)
	}

	return &localSuggester{index: index}
}

func (c *localSuggester) SuggestLicenses(q string) ([]models.Completion, error) {
	hits := make([]models.Completion, 0)
	return hits, nil
}
