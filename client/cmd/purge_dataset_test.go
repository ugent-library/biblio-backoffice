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

type PurgeDatasetSuite struct {
	suite.Suite
}

func (s *PurgeDatasetSuite) SetupSuite() {
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

func (s *PurgeDatasetSuite) TestPurgeSingle() {
	t := s.T()

	stdOut, _, err := purgeDataset("notexists")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not find dataset with id notexists", stdOut)

	stdOut, _, err = purgeDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "purged dataset 00000000000000000000000001", stdOut)
}

func TestPurgeDatasetSuite(t *testing.T) {
	suite.Run(t, new(PurgeDatasetSuite))
}

func purgeDataset(id string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"dataset", "purge", id})

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
