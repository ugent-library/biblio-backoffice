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

type DatasetReindexSuite struct {
	suite.Suite
}

func (s *DatasetReindexSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)

	t := s.T()

	file, err := os.ReadFile("../../etc/fixtures/complete.dataset.json")
	if err != nil {
		t.Fatal(err)
	}

	addCmdOutFile, err := toJSONL(file)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = addDataset(addCmdOutFile)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *DatasetReindexSuite) TestReindex() {
	t := s.T()

	stdOut, _, err := reindexDataset()
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `^Indexing to a new index\n`, stdOut)
	assert.Regexp(t, `\nDone\n$`, stdOut)
}

func (s *DatasetReindexSuite) TearDownSuite() {
	t := s.T()

	_, _, err := purgeDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatasetReindexSuite(t *testing.T) {
	suite.Run(t, new(DatasetReindexSuite))
}

func reindexDataset() (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"dataset", "reindex"})

	errGet := rootCmd.Execute()
	if errGet != nil {
		return "", "", errGet
	}

	resOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	resErr, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	return string(resOut), string(resErr), nil
}
