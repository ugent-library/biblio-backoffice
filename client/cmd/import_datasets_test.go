package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ImportDatasetsSuite struct {
	suite.Suite
}

func (s *ImportDatasetsSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)
}

func (s *ImportDatasetsSuite) TestValidImport() {
	t := s.T()

	json := `{
		"access_level": "info:eu-repo/semantics/openAccess",
		"author": [
		  {
			"person_id": "00000000-0000-0000-0000-000000000001"
		  }
		],
		"creator_id": "00000000-0000-0000-0000-000000000001",
		"related_organizations": [
		  {
			"organization_id": "CA20",
		  }
		],
		"doi": "doi",
		"identifiers": {
		  "DOI": [
		    "doi"
		  ]
		},
		"format": [
		  "text/csv"
		],
		"publisher": "publisher",
		"has_been_public": true,
		"license": "CC-BY-4.0",
		"locked": false,
		"status": "public",
		"title": "title",
		"url": "URL",
		"year": "2020",
		"id": "00000000000000000000000011",
		"snapshot_id": "00000000000000000000000011",
		"date_created": "2023-02-07T11:17:35.203928951+01:00",
		"date_updated": "2023-02-07T11:33:49.661899119+01:00",
		"date_from": "2023-02-07T11:33:49.663142+01:00"
	  }`

	json, err := toJSONL([]byte(json))
	if err != nil {
		t.Fatal(err)
	}

	stdOut, _, err := importDataset(json)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `stored and indexed dataset .* at line .*`, string(stdOut))

	// Retrieve the dataset
	stdOut, _, err = getDataset("00000000000000000000000011")

	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, json, stdOut, "the imported and retrieved records aren't equal.")
}

func (s *ImportDatasetsSuite) TestImportNonJSONLInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"dataset", "import"})

	in := strings.NewReader("invalid\n")
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not read json input: invalid character 'i' looking for beginning of value\n", string(stdOut))
}

func (s *ImportDatasetsSuite) TestImportEmptyJSONLInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"dataset", "import"})

	in := strings.NewReader("{}\n")
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `failed to validate dataset at line .: dataset.id.required\[\/id\], dataset.status.required\[\/status\]`, string(stdOut))
}

func (s *ImportDatasetsSuite) TearDownSuite() {
	t := s.T()
	_, _, err := purgeDataset("00000000000000000000000011")

	if err != nil {
		t.Fatal(err)
	}
}

func TestImportDatasetsSuite(t *testing.T) {
	suite.Run(t, new(ImportDatasetsSuite))
}

func importDataset(jsonl string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"dataset", "import"})

	in := strings.NewReader(jsonl)
	rootCmd.SetIn(in)

	err := rootCmd.Execute()
	if err != nil {
		return "", "", err
	}

	impCmdOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	impCmdErr, err := ioutil.ReadAll(stdErr)
	if err != nil {
		return "", "", err
	}

	return string(impCmdOut), string(impCmdErr), nil
}
