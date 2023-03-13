package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ImportPublicationsSuite struct {
	suite.Suite
}

func (s *ImportPublicationsSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)
}

func (s *ImportPublicationsSuite) TestValidImport() {
	t := s.T()

	json := `{
		"id": "00000000000000000000000011",
		"snapshot_id": "00000000000000000000000011",
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

	stdOut, _, err := importPublication(json)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `stored and indexed publication .* at line .*`, string(stdOut))

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000011")

	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, json, stdOut, "the imported and retrieved records aren't equal.")
}

func (s *ImportPublicationsSuite) TestImportNonJSONLInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "import"})

	in := strings.NewReader("invalid\n")
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not read json input: invalid character 'i' looking for beginning of value\n", string(stdOut))
}

func (s *ImportPublicationsSuite) TestImportEmptyJSONLInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "import"})

	in := strings.NewReader("{}\n")
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `failed to validate publication at line .: publication.id.required\[\/id\], publication.type.required\[\/type\], publication.classification.required\[\/classification\], publication.status.required\[\/status\]`, string(stdOut))
}

func TestImportPublicationsSuite(t *testing.T) {
	suite.Run(t, new(ImportPublicationsSuite))
}

func importPublication(jsonl string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "import"})

	in := strings.NewReader(jsonl)
	rootCmd.SetIn(in)

	err := rootCmd.Execute()
	if err != nil {
		return "", "", err
	}

	impCmdOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	impCmdErr, err := ioutil.ReadAll(stdErr)
	if err != nil {
		return "", "", err
	}

	return string(impCmdOut), string(impCmdErr), nil
}
