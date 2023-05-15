package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AddPublicationSuite struct {
	suite.Suite
}

func (s *AddPublicationSuite) SetupTest() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)
}

// Test empty input
func (s *AddPublicationSuite) TestAddEmptyInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "add"})

	in := strings.NewReader("")
	rootCmd.SetIn(in)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", string(stdOut))
}

// Test non-JSONL input
func (s *AddPublicationSuite) TestAddNonJSONLInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "add"})

	in := strings.NewReader("invalid\n")
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not read json input: invalid character 'i' looking for beginning of value\n", string(stdOut))
}

// Test empty JSONL object
func (s *AddPublicationSuite) TestAddEmptyJSONLInput() {
	t := s.T()

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "add"})

	in := strings.NewReader("{}\n")
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `failed to validate publication .* at line .: publication.type.required\[\/type\]`, string(stdOut))
}

// Create new minimal valid record
func (s *AddPublicationSuite) TestAddMinimalValidRecord() {
	t := s.T()

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

	jsonl, errJSONL := toJSONL([]byte(json))
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "add"})

	in := strings.NewReader(jsonl)
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `stored and indexed publication .* at line .*`, string(stdOut))
}

// Test adding multiple records at once
func (s *AddPublicationSuite) TestAddingMultiplePublications() {
	t := s.T()

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

	input := ""
	for i := 1; i <= 5; i += 1 {
		jsonl, errJSONL := toJSONL([]byte(json))
		if errJSONL != nil {
			t.Fatal(errJSONL)
		}
		input = input + jsonl
	}

	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "add"})

	in := strings.NewReader(input)
	rootCmd.SetIn(in)

	rootCmd.Execute()

	stdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(stdOut), "\n")

	seq := 0
	for _, line := range lines {
		if line != "" {
			seq++
			assert.Regexp(t, fmt.Sprintf("stored and indexed publication .* at line %d", seq), line)
		}
	}
}

// Test if all fields return properly
func (s *AddPublicationSuite) TestAddAndGetCompletePublications() {
	t := s.T()

	// Book chapter

	file, errFile := os.ReadFile("../../etc/fixtures/complete.book_chapter.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ := addKey(string(file), "id", "00000000000000000000000001")

	jsonl, err := toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err := getPublication("00000000000000000000000001")

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

	// Book editor

	file, errFile = os.ReadFile("../../etc/fixtures/complete.book_editor.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ = addKey(string(file), "id", "00000000000000000000000002")

	jsonl, err = toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000002")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(jsonl, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(stdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Book

	file, errFile = os.ReadFile("../../etc/fixtures/complete.book.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ = addKey(string(file), "id", "00000000000000000000000003")

	jsonl, err = toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000003")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(jsonl, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(stdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Conference

	file, errFile = os.ReadFile("../../etc/fixtures/complete.conference.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ = addKey(string(file), "id", "00000000000000000000000004")

	jsonl, err = toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000004")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(jsonl, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(stdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Dissertation

	file, errFile = os.ReadFile("../../etc/fixtures/complete.dissertation.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ = addKey(string(file), "id", "00000000000000000000000005")

	jsonl, err = toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000005")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(jsonl, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(stdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Issue editor

	file, errFile = os.ReadFile("../../etc/fixtures/complete.issue_editor.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ = addKey(string(file), "id", "00000000000000000000000006")

	jsonl, err = toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000006")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(jsonl, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(stdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Journal article

	file, errFile = os.ReadFile("../../etc/fixtures/complete.journal_article.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ = addKey(string(file), "id", "00000000000000000000000007")

	jsonl, err = toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000007")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(jsonl, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(stdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Miscellaneous

	file, errFile = os.ReadFile("../../etc/fixtures/complete.miscellaneous.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	rec, _ = addKey(string(file), "id", "00000000000000000000000008")

	jsonl, err = toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000008")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(jsonl, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(stdOut, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")
}

// Test validation

// Validation: type field is required
// Validation: type is invalid
// Validation: classification is required
// Validation: classification is invalid
// Validation: status is required
// Validation: status is invalid
// Validation: title field is required
// Validation: year field is required
// Validation: year field is invalid
// Validation: language code is invalid
// Validation: abstract => id is required
// Validation: abstract => lang is required
// Validation: abstract => lang is invalid
// Validation: abstraxt => text is required
// Validation: laysummary => id is required
// Validation: laysummary => lang is required
// Validation: laysummary => lang is invalid
// Validation: laysummary => text is required
// Validation: if status = public => at least 1 author
// Validation: if status = public, !legacy, usesAuthor(), !extern => at least 1 ugent author
// Validation: if status = public, usesEditor, !usesAuthor => at least 1 editor
// Validation: if status = public, usesEditor, !usesAuthor, !extern => at least 1 ugent editor
// Validation: if status = public, !legacy, usesSupervisor => at least 1 supervisor
// Validation: author => id is required
// Validation: author => first name is required
// Validation: author => last name is required
// Validation: author => credit role is required
// Validation: author => credit role is invalid
// Validation: editor => id is required
// Validation: editor => first name is required
// Validation: editor => last name is required
// Validation: editor => credit role is required
// Validation: editor => credit role is invalid
// Validation: supervisor => id is required
// Validation: supervisor => first name is required
// Validation: supervisor => last name is required
// Validation: supervisor => credit role is required
// Validation: supervisor => credit role is invalid
// Validation: project => id is required
// Validation: department => id is required
// Validation: relatedDataset => id is required
// Validation: file ........
// Validation: link => invalid relation

// Specifics

// Dissertation
// JournalArticle
// Miscellaneous
// book
// book_chapter
// conference
// book_editor
// issue_editor

func TestAddPublicationSuite(t *testing.T) {
	suite.Run(t, new(AddPublicationSuite))
}

func (s *AddPublicationSuite) TearDownSuite() {
	purgePublication("00000000000000000000000001")
	purgePublication("00000000000000000000000002")
	purgePublication("00000000000000000000000003")
	purgePublication("00000000000000000000000004")
	purgePublication("00000000000000000000000005")
	purgePublication("00000000000000000000000006")
	purgePublication("00000000000000000000000007")
	purgePublication("00000000000000000000000008")
}

func addPublication(jsonl string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"publication", "add"})

	in := strings.NewReader(jsonl)
	rootCmd.SetIn(in)

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

func toJSONL(file []byte) (string, error) {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, file); err != nil {
		return "", err
	}

	output, _ := ioutil.ReadAll(buffer)
	json := string(output)
	json = json + "\n"

	return json, nil
}

func removeKey(input string, keys ...string) (string, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		return "", err
	}

	output := input

	for _, key := range keys {
		if _, exists := m[key]; exists {
			delete(m, key)
			outputData, err := json.Marshal(m)
			if err != nil {
				return "", err
			}
			output = string(outputData)
		}
	}

	return output, nil
}

func addKey(input, key, value string) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		return "", err
	}

	m[key] = value

	output, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
