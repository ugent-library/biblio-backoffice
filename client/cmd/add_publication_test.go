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

	errAdd := rootCmd.Execute()
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	addCmdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", string(addCmdOut))
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

	addCmdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Error: could not read json: invalid character 'i' looking for beginning of value\n", string(addCmdOut))
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

	addCmdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `Error: validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(addCmdOut))
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

	addCmdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(addCmdOut))

	// @todo Add test
}

// Update existing valid record

// Test if all fields return properly
func (s *AddPublicationSuite) estAddAndGetCompletePublications() {
	t := s.T()

	// Book chapter

	file, errFile := os.ReadFile("../../etc/fixtures/complete.book_chapter.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL := toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd := addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut := getPublication("00000000000000000000000001")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err := removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err := removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Book editor

	file, errFile = os.ReadFile("../../etc/fixtures/complete.book_editor.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000002")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Book

	file, errFile = os.ReadFile("../../etc/fixtures/complete.book.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000003")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Conference

	file, errFile = os.ReadFile("../../etc/fixtures/complete.conference.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000004")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Dissertation

	file, errFile = os.ReadFile("../../etc/fixtures/complete.dissertation.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000005")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Issue editor

	file, errFile = os.ReadFile("../../etc/fixtures/complete.issue_editor.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000006")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Journal article

	file, errFile = os.ReadFile("../../etc/fixtures/complete.journal_article.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000007")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// Miscellaneous

	file, errFile = os.ReadFile("../../etc/fixtures/complete.miscellaneous.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000008")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
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
// Validation: abstract => ID is required
// Validation: abstract => lang is required
// Validation: abstract => lang is invalid
// Validation: abstraxt => text is required
// Validation: laysummary => ID is required
// Validation: laysummary => lang is required
// Validation: laysummary => lang is invalid
// Validation: laysummary => text is required
// Validation: if status = public => at least 1 author
// Validation: if status = public, !legacy, usesAuthor(), !extern => at least 1 ugent author
// Validation: if status = public, usesEditor, !usesAuthor => at least 1 editor
// Validation: if status = public, usesEditor, !usesAuthor, !extern => at least 1 ugent editor
// Validation: if status = public, !legacy, usesSupervisor => at least 1 supervisor
// Validation: author => ID is required
// Validation: author => first name is required
// Validation: author => last name is required
// Validation: author => credit role is required
// Validation: author => credit role is invalid
// Validation: editor => ID is required
// Validation: editor => first name is required
// Validation: editor => last name is required
// Validation: editor => credit role is required
// Validation: editor => credit role is invalid
// Validation: supervisor => ID is required
// Validation: supervisor => first name is required
// Validation: supervisor => last name is required
// Validation: supervisor => credit role is required
// Validation: supervisor => credit role is invalid
// Validation: project => ID is required
// Validation: department => ID is required
// Validation: relatedDataset => ID is required
// Validation: file ........
// Validation: link => invalid relation

// Specifics

// Dissertation xx
// JournalArticle xx
// Miscellaneous xx
// book xx
// book_chapter xxx
// conference xx
// book_editor xxx
// issue_editor x

func TestAddPublicationSuite(t *testing.T) {
	suite.Run(t, new(AddPublicationSuite))
}

func addPublication(jsonl string) error {
	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "add"})

	in := strings.NewReader(jsonl)
	rootCmd.SetIn(in)

	errAdd := rootCmd.Execute()
	if errAdd != nil {
		return errAdd
	}

	// addCmdOut, err := ioutil.ReadAll(actual)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	return nil
}

func getPublication(id string) (string, error) {
	actual := bytes.NewBufferString("")
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	rootCmd.SetArgs([]string{"publication", "get", id})

	errAdd := rootCmd.Execute()
	if errAdd != nil {
		return "", errAdd
	}

	getCmdOut, err := ioutil.ReadAll(actual)
	if err != nil {
		return "", err
	}

	return string(getCmdOut), nil
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
