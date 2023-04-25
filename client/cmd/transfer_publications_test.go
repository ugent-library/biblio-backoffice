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
		"id": "00000000000000000000000071",
		"snapshot_id": "00000000000000000000000071",
		"title": "title",
		"type": "dissertation",
		"status": "public",
		"year": "2023",
		"classification": "D1",
		"extern": false,
		"has_been_public": false,
		"legacy": false,
		"locked": false,
		"vabb_approved": false,
		"date_created": "2023-02-07T11:17:35.203928951+01:00",
		"date_updated": "2023-02-07T11:33:49.661899119+01:00",
		"date_from": "2023-02-07T11:33:49.663142+01:00",
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

	_, _, err = importPublication(json)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *TransferPublicationsSuite) TestTransfer() {
	t := s.T()

	stdOut, _, err := transferPublications("00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002", "00000000000000000000000071")
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `p: 00000000000000000000000071: s: .* ::: user: 00000000-0000-0000-0000-000000000001 -> 00000000-0000-0000-0000-000000000002\np: 00000000000000000000000071: s: .* ::: author: 00000000-0000-0000-0000-000000000001 -> 00000000-0000-0000-0000-000000000002\np: 00000000000000000000000071: s: .* ::: supervisor: 00000000-0000-0000-0000-000000000001 -> 00000000-0000-0000-0000-000000000002`, string(stdOut))

	stdOut, _, err = getPublication("00000000000000000000000071")
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
	assert.Equal(t, p.User.ID, "00000000-0000-0000-0000-000000000002")
}

func (s *TransferPublicationsSuite) TearDownSuite() {
	t := s.T()
	_, _, err := purgePublication("00000000000000000000000071")

	if err != nil {
		t.Fatal(err)
	}
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
