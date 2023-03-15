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

type PublicationReindexSuite struct {
	suite.Suite
}

func (s *PublicationReindexSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)

	t := s.T()

	file, err := os.ReadFile("../../etc/fixtures/complete.book_chapter.json")
	if err != nil {
		t.Fatal(err)
	}

	addCmdOutFile, err := toJSONL(file)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = addPublication(addCmdOutFile)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *PublicationReindexSuite) TestReindex() {
	t := s.T()

	stdOut, _, err := reindexPublication()
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `Indexing to a new index\nIndexed .* publications...\nSwitching to new index...\nIndexing changes since start of reindex...\nDone`, stdOut)
}

func (s *PublicationReindexSuite) TearDownSuite() {
	t := s.T()

	_, _, err := purgeDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
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
