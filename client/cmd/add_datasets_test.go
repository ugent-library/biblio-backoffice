package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AddDatasetsSuite struct {
	suite.Suite
}

func (s *AddDatasetsSuite) SetupTest() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)
}

// Test empty input
func (s *AddDatasetsSuite) TestAddEmptyInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"dataset", "add"})

	in := strings.NewReader("")
	rootCmd.SetIn(in)

	errAdd := rootCmd.Execute()
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	addCmdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", string(addCmdOut))
}

// Test non-JSONL input
func (s *AddDatasetsSuite) TestAddNonJSONLInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"dataset", "add"})

	in := strings.NewReader("invalid\n")
	rootCmd.SetIn(in)

	rootCmd.Execute()

	addCmdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not read json input: invalid character 'i' looking for beginning of value\n", string(addCmdOut))
}

// Test empty JSONL object
func (s *AddDatasetsSuite) TestAddEmptyJSONLInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"dataset", "add"})

	in := strings.NewReader("{}\n")
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `stored and indexed dataset .* at line 1`, string(stdOut))
}

// Test adding multiple records at once
func (s *AddDatasetsSuite) TestAddingMultipleDatasets() {
	t := s.T()

	json := `{
		"access_level": "info:eu-repo/semantics/openAccess",
		"author": [
		  {
			"first_name": "first name",
			"full_name": "full name",
			"id": "00000000-0000-0000-0000-000000000001",
			"last_name": "last name",
			"ugent_id": [
			  "000000000001"
			],
			"department": [
			  {
				"id": "CA20",
				"name": "Department of Research Affairs"
			  }
			]
		  }
		],
		"creator": {
		  "id": "00000000-0000-0000-0000-000000000001",
		  "name": "full name"
		},
		"department": [
		  {
			"id": "CA20",
			"tree": [
			  {
				"id": "UGent"
			  },
			  {
				"id": "CA"
			  },
			  {
				"id": "CA20"
			  }
			]
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
		"snapshot_id": "00000000000000000000000001",
		"status": "public",
		"title": "title",
		"url": "URL",
		"year": "2020"
	  }`

	input := ""
	for i := 1; i <= 5; i += 1 {
		jsonl, errJSONL := toJSONL([]byte(json))
		if errJSONL != nil {
			t.Fatal(errJSONL)
		}
		input = input + jsonl
	}

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"dataset", "add"})

	in := strings.NewReader(input)
	rootCmd.SetIn(in)

	rootCmd.Execute()

	addCmdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(addCmdOut), "\n")

	seq := 0
	for _, line := range lines {
		if line != "" {
			seq++
			assert.Regexp(t, fmt.Sprintf("stored and indexed dataset .* at line %d", seq), line)
		}
	}
}

// Test if all fields return properly
func (s *AddDatasetsSuite) TestAddAndGetCompleteDatasets() {
	t := s.T()

	file, errFile := os.ReadFile("../../etc/fixtures/complete.dataset.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a dataset
	addCmdOutFile, errJSONL := toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	_, _, errAdd := addDataset(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the dataset
	getCmdStdOut, _, errGetCmdOut := getDataset("00000000000000000000000001")

	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err := removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err := removeKey(getCmdStdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")
}

func TestAddDatasetsSuite(t *testing.T) {
	suite.Run(t, new(AddDatasetsSuite))
}

func addDataset(jsonl string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"dataset", "add"})

	in := strings.NewReader(jsonl)
	rootCmd.SetIn(in)

	err := rootCmd.Execute()
	if err != nil {
		return "", "", err
	}

	addCmdOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	addCmdErr, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	return string(addCmdOut), string(addCmdErr), nil
}
