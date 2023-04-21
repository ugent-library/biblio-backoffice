package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type CleanupPublicationsSuite struct {
	suite.Suite
}

func (s *CleanupPublicationsSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)

	t := s.T()

	json := `{
		"id": "00000000000000000000000011",
		"title": "title",
		"type": "journal_article",
		"status": "public",
		"year": "2023",
		"department": [
			{
			  "id": "CA20"
			}
		],
		"keyword": [
			"keyword",
			"",
			"  keyword2  "
		],
		"author": [
			{
				"credit_role": [
					"first_author"
				],
				"first_name": "first name",
				"full_name": "full name",
				"id": "00000000-0000-0000-0000-000000000000",
				"last_name": "last name",
				"ugent_id": [
					"000000000000"
				],
				"department": [
				{
					"id": "AA00",
					"name": "department name"
				}
				]
			}
		]
	}`

	addCmdOutFile, err := toJSONL([]byte(json))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = addPublication(addCmdOutFile)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *CleanupPublicationsSuite) TestCleanup() {
	t := s.T()

	stdOut, _, err := cleanupPublications()
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `fixed publication\[snapshot_id: .*, id: 00000000000000000000000011\]\ndone. cleaned 1 publications.`, string(stdOut))

	// Get the publication and assert that keyword and department tree were added

	stdOut, _, err = getPublication("00000000000000000000000011")
	if err != nil {
		t.Fatal(err)
	}

	p := &models.Publication{}
	err = json.Unmarshal([]byte(stdOut), p)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, p.Keyword, 2, "number of keywords isn't 2, got:", len(p.Keyword))
	assert.Equal(t, p.Keyword[0], "keyword")
	assert.Equal(t, p.Keyword[1], "keyword2")
	assert.NotNil(t, p.Department[0].Tree, "department tree is missing, got:", p.Department)
}

func (s *CleanupPublicationsSuite) TearDownSuite() {
	t := s.T()
	_, _, err := purgePublication("00000000000000000000000011")

	if err != nil {
		t.Fatal(err)
	}
}

func TestCleanupPublicationsSuite(t *testing.T) {
	suite.Run(t, new(CleanupPublicationsSuite))
}

func cleanupPublications() (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "cleanup"})

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
