package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type GetAllPublicationsSuite struct {
	suite.Suite
}

func (s *GetAllPublicationsSuite) SetupTest() {
	t := s.T()

	viper.Set("port", "3999")
	viper.Set("insecure", true)

	json := `{
		"title": "title",
		"type": "journal_article",
		"status": "public",
		"year": "2023",
		"author": [
			{
				"credit_role": [
					"first_author"
				],
				"person_id": "00000000-0000-0000-0000-000000000000"
			}
		]
	}`

	for i := 1; i <= 5; i += 1 {
		rec, _ := addKey(json, "id", fmt.Sprintf("0000000000000000000000002%d", i))
		jsonl, err := toJSONL([]byte(rec))
		if err != nil {
			t.Fatal(err)
		}

		_, _, err = addPublication(jsonl)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Test empty input
func (s *GetAllPublicationsSuite) TestGetAllPublications() {
	t := s.T()

	stdOut, _, err := getAllPublications()
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(strings.NewReader(stdOut))

	r, _ := regexp.Compile("0000000000000000000000002.")
	counter := 0
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatal("could not read line from input: %w", err)
		}

		p := &models.Publication{}
		err = json.Unmarshal([]byte(line), p)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, p.ID, "publication ID is nil")

		if r.MatchString(p.ID) {
			counter++
		}
	}

	assert.Equal(t, 5, counter, "number of publications isn't 5, got %s", counter)
}

func (s *GetAllPublicationsSuite) TearTest() {
	t := s.T()

	for i := 1; i <= 5; i += 1 {
		d, e, l := purgePublication(fmt.Sprintf("0000000000000000000000002%d", i))

		t.Log(d)
		t.Log(e)
		t.Log(l)
	}
}

func TestGetAllPublicationSuite(t *testing.T) {
	suite.Run(t, new(GetAllPublicationsSuite))
}

func getAllPublications() (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "get-all"})

	err := rootCmd.Execute()
	if err != nil {
		return "", "", err
	}

	cmdOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	cmdErr, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	return string(cmdOut), string(cmdErr), nil
}
