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

type PurgePublicationSuite struct {
	suite.Suite
}

func (s *PurgePublicationSuite) SetupSuite() {
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

func (s *PurgePublicationSuite) TestPurgeAll() {
	t := s.T()

	stdOut, _, err := purgePublication("notexists")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not find publication with id notexists", stdOut)

	stdOut, _, err = purgePublication("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "purged publication 00000000000000000000000001", stdOut)
}

func TestPurgePublicationSuite(t *testing.T) {
	suite.Run(t, new(PurgePublicationSuite))
}

func purgePublication(id string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "purge", id})

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
