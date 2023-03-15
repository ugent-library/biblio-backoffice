package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PurgeAllDatasetsSuite struct {
	suite.Suite
}

func (s *PurgeAllDatasetsSuite) SetupSuite() {
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

func (s *PurgeAllDatasetsSuite) TestPurgeAll() {
	t := s.T()

	stdOut, _, err := purgeAllDatasets(false)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "no confirmation flag set. you need to set the --yes flag", stdOut)

	// Artificial sleep. PurgeAll triggers an async ES6 task. Multiple tasks registered with
	// ES6 cause conflicts. A timeout gives ES a chance to resolve tasks sequentially.
	time.Sleep(5 * time.Second)

	stdOut, _, err = purgeAllDatasets(true)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "purged all datasets", stdOut)

	// Retrieve the dataset
	stdOut, _, err = getDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not find dataset with id 00000000000000000000000001", stdOut)
}

func TestPurgeAllDatasetsSuite(t *testing.T) {
	suite.Run(t, new(PurgeAllDatasetsSuite))
}

func purgeAllDatasets(confirm bool) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	if confirm {
		rootCmd.SetArgs([]string{"dataset", "purge-all", "--yes"})
	} else {
		rootCmd.SetArgs([]string{"dataset", "purge-all"})
	}

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
