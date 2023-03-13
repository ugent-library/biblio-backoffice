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

type TransferPublicationsSuite struct {
	suite.Suite
}

func (s *TransferPublicationsSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)

	t := s.T()

	json := `{
		"id": "00000000000000000000000003",
		"title": "title",
		"type": "dissertation",
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
		"user": {
			"id": "00000000-0000-0000-0000-000000000001",
			"name": "full name"
		},
		"author": [
			{
				"credit_role": [
					"first_author"
				],
				"first_name": "first name",
				"full_name": "full name",
				"id": "00000000-0000-0000-0000-000000000001",
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
		],
		"defense_date": "2023-01-01",
		"defense_place": "place of defense",
		"defense_time": "12:00",
		"supervisor": [
			{
			  "first_name": "first name",
			  "full_name": "full name",
			  "id": "00000000-0000-0000-0000-000000000001",
			  "last_name": "last name",
			  "orcid": "0000-0000-0000-0000",
			  "ugent_id": [
				"000000000000",
				"000000000001"
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

	json, err := toJSONL([]byte(json))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = addPublication(json)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *TransferPublicationsSuite) TestTransfer() {
	t := s.T()

	stdOut, _, err := transferPublications("00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002", "00000000000000000000000003")
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `p: 00000000000000000000000003: s: .* ::: author: 00000000-0000-0000-0000-000000000001 -> 00000000-0000-0000-0000-000000000002\np: 00000000000000000000000003: s: .* ::: supervisor: 00000000-0000-0000-0000-000000000001 -> 00000000-0000-0000-0000-000000000002`, string(stdOut))

	stdOut, _, err = getPublication("00000000000000000000000003")
	if err != nil {
		t.Fatal(err)
	}

	p := &models.Publication{}
	err = json.Unmarshal([]byte(stdOut), p)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, p.Author[0].ID, "00000000-0000-0000-0000-000000000002")
	assert.Equal(t, p.Supervisor[0].ID, "00000000-0000-0000-0000-000000000002")

	// TODO assert p.User. Implement this after import_publications has tests.
}

func TestTransferPublicationsSuite(t *testing.T) {
	suite.Run(t, new(TransferPublicationsSuite))
}

func transferPublications(src string, dst string, id string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "transfer", src, dst, id})

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
