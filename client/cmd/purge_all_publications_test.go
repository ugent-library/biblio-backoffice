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

type PurgeAllPublicationsSuite struct {
	suite.Suite
}

func (s *PurgeAllPublicationsSuite) SetupSuite() {
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

func (s *PurgeAllPublicationsSuite) TestPurgeAll() {
	t := s.T()

	stdOut, _, err := purgeAllPublications(false)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "no confirmation flag set. you need to set the --yes flag.", stdOut)

	stdOut, _, err = purgeAllPublications(true)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "purged all publications from biblio backoffice.", stdOut)
}

func TestPurgeAllPublicationsSuite(t *testing.T) {
	suite.Run(t, new(PurgeAllPublicationsSuite))
}

func purgeAllPublications(confirm bool) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	if confirm {
		rootCmd.SetArgs([]string{"publication", "purge-all", "--yes"})
	} else {
		rootCmd.SetArgs([]string{"publication", "purge-all"})
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
