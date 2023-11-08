package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UpdatePublicationSuite struct {
	suite.Suite
}

func (s *UpdatePublicationSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)

	t := s.T()

	file, err := os.ReadFile("../../etc/fixtures/complete.book_chapter.json")
	if err != nil {
		t.Fatal(err)
	}

	rec, _ := addKey(string(file), "id", "00000000000000000000000081")

	jsonl, err := toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}
}

// Test empty input
func (s *UpdatePublicationSuite) TestUpdateEmptyInput() {
	t := s.T()

	jsonl := ""

	_, _, err := updatePublication(jsonl)

	assert.Equal(t, "could not read from stdin: EOF", err.Error())
}

// Test invalid JSON input
func (s *UpdatePublicationSuite) TestUpdateInvalidJSON() {
	t := s.T()

	jsonl := "invalid\n"

	stdOut, _, err := updatePublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not read json input: invalid character 'i' looking for beginning of value\n", stdOut)
}

func (s *UpdatePublicationSuite) TestUpdateEmptyJSONLInput() {
	t := s.T()
	jsonl := "{}\n"

	stdOut, _, err := updatePublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `failed to validate publication : publication.id.required\[\/id\], publication.type.required\[\/type\], publication.classification.required\[\/classification\], publication.status.required\[\/status\]`, string(stdOut))
}

// Test if a publication can be succesfully updated
func (s *UpdatePublicationSuite) TestUpdateValid() {
	t := s.T()

	// Retrieve the publication as user A
	pubA, _, err := getPublication("00000000000000000000000081")
	if err != nil {
		t.Fatal(err)
	}

	updPubA, err := changeTitle(pubA, "update title")
	if err != nil {
		t.Fatal(err)
	}

	jsonl, err := toJSONL([]byte(updPubA))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = updatePublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	pubB, _, err := getPublication("00000000000000000000000081")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	updPubA, err = removeKey(updPubA, "snapshot_id", "date_from", "date_updated", "user_id", "last_user_id")
	if err != nil {
		t.Fatal(err)
	}

	pubB, err = removeKey(pubB, "snapshot_id", "date_from", "date_updated", "use_idr", "last_user_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, updPubA, pubB, "the added and retrieved records aren't equal.")
}

// Test an update conflict
func (s *UpdatePublicationSuite) TestUpdateConflict() {
	t := s.T()

	// Retrieve the publication as user A
	getCmdStdOut, _, err := getPublication("00000000000000000000000081")

	if err != nil {
		t.Fatal(err)
	}

	// user B pushes a change to the publication

	updated, err := changeTitle(getCmdStdOut, "new title")
	if err != nil {
		t.Fatal(err)
	}

	jsonl, err := toJSONL([]byte(updated))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = updatePublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// // User A tries to change the publication with the old snapshotID

	updated, err = changeTitle(getCmdStdOut, "alternate title")
	if err != nil {
		t.Fatal(err)
	}

	jsonl, err = toJSONL([]byte(updated))
	if err != nil {
		t.Fatal(err)
	}

	updCmdOut, _, err := updatePublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `failed to update publication: conflict detected for publication\[snapshot_id: .*, id: .*\] : version conflict`, updCmdOut)
}

func (s *UpdatePublicationSuite) TearDownSuite() {
	t := s.T()

	_, _, err := purgePublication("00000000000000000000000081")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePublicationSuite(t *testing.T) {
	suite.Run(t, new(UpdatePublicationSuite))
}

func updatePublication(jsonl string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "update"})

	in := strings.NewReader(jsonl)
	rootCmd.SetIn(in)

	err := rootCmd.Execute()
	if err != nil {
		return "", "", err
	}

	updCmdOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	updCmdErr, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	return string(updCmdOut), string(updCmdErr), nil
}

func changeTitle(input string, value string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		return "", err
	}

	m["title"] = value

	output, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
