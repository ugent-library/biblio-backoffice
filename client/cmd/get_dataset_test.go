package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetDatasetSuite struct {
	suite.Suite
}

func (s *GetDatasetSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)
}

func (s *GetDatasetSuite) TestGetDataset() {
	t := s.T()

	file, errFile := os.ReadFile("../../etc/fixtures/complete.dataset.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a dataset
	jsonl, err := toJSONL(file)
	if err != nil {
		t.Fatal(err)
	}

	_, _, errAdd := addDataset(jsonl)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the dataset
	stdOut, _, err := getDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err := removeKey(jsonl, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err := removeKey(stdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")
}

func (s *GetDatasetSuite) TestDatasetNotFound() {
	t := s.T()

	stdOut, _, err := getDataset("notfound")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not find dataset with id notfound", stdOut)
}

func (s *GetDatasetSuite) TearDownSuite() {
	t := s.T()

	_, _, err := purgeDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDatasetSuite(t *testing.T) {
	suite.Run(t, new(GetDatasetSuite))
}

func getDataset(id string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"dataset", "get", id})

	errGet := rootCmd.Execute()
	if errGet != nil {
		return "", "", errGet
	}

	getCmdOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	getCmdErr, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	return string(getCmdOut), string(getCmdErr), nil
}
