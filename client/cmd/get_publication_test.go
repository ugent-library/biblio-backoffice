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

type GetPublicationSuite struct {
	suite.Suite
}

func (s *GetPublicationSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)
}

func (s *GetPublicationSuite) TestGetPublication() {
	t := s.T()

	file, errFile := os.ReadFile("../../etc/fixtures/complete.book_chapter.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ := addKey(string(file), "id", "00000000000000000000000031")

	jsonl, err := toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err := getPublication("00000000000000000000000031")
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

func (s *GetPublicationSuite) TestPublicationNotFound() {
	t := s.T()

	stdOut, _, err := getPublication("notfound")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not find publication with id notfound", stdOut)
}

func (s *GetPublicationSuite) TearDownSuite() {
	t := s.T()

	_, _, err := purgePublication("00000000000000000000000031")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPublicationSuite(t *testing.T) {
	suite.Run(t, new(GetPublicationSuite))
}

func getPublication(id string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "get", id})

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
