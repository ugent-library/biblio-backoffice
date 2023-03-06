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

	// Book chapter

	file, errFile := os.ReadFile("../../etc/fixtures/complete.book_chapter.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication as user A
	addCmdOutFile, errJSONL := toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	_, _, errAdd := addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}
}

// Test empty input
func (s *UpdatePublicationSuite) TestEmptyInput() {
	t := s.T()

	jsonl := ""

	_, _, err := updatePublication(jsonl)

	assert.Equal(t, "could not read from stdin: EOF", err.Error())
}

// Test invalid JSON input
func (s *UpdatePublicationSuite) TestInvalidJSON() {
	// t := s.T()

	//TODO
}

// Test if a publication can be succesfully updated
func (s *UpdatePublicationSuite) TestUpdate() {
	t := s.T()

	// Retrieve the publication as user A
	pubA, _, errGetCmdOut := getPublication("00000000000000000000000001")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
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

	pubB, _, err := getPublication("00000000000000000000000001")
	if err != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	updPubA, err = removeKey(updPubA, "snapshot_id", "date_from", "date_updated", "user", "last_user")
	if err != nil {
		t.Fatal(err)
	}

	pubB, err = removeKey(pubB, "snapshot_id", "date_from", "date_updated", "user", "last_user")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, updPubA, pubB, "the added and retrieved records aren't equal.")
}

// Test an update conflict
func (s *UpdatePublicationSuite) TestUpdateConflict() {
	t := s.T()

	// Retrieve the publication as user A
	getCmdStdOut, _, errGetCmdOut := getPublication("00000000000000000000000001")

	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
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
	var m map[string]interface{}
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
