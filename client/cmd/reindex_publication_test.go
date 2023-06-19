package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PublicationReindexSuite struct {
	suite.Suite
}

func (s *PublicationReindexSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)
}

func (s *PublicationReindexSuite) TestReindex() {
	t := s.T()

	stdOut, _, err := reindexPublication()
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `^Indexing to a new index\n`, stdOut)
	assert.Regexp(t, `\nDone\n$`, stdOut)
}

func TestPublicationReindexSuite(t *testing.T) {
	suite.Run(t, new(PublicationReindexSuite))
}

func reindexPublication() (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "reindex"})

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
